package address

import (
	"errors"
	"reflect"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

type StakePointer struct {
	Slot      uint64
	TxIndex   uint64
	CertIndex uint64
}

// A pointer address indirectly specifies the staking key that should control the stake for the address.
type PointerAddress struct {
	Network network.NetworkInfo
	Payment StakeCredential
	Stake   StakePointer
}

func VariableNatEncode(num uint64) []byte {
	var output []byte
	output = append(output, byte(num)&0x7F)
	num /= 128
	for num > 0 {
		output = append(output, byte(num)&0x7F|0x80)
		num /= 128
	}
	reverseAny(output)
	return output
}

func VariableNatDecode(raw []byte) (uint64, int, error) {
	var output uint64
	output = 0
	bytes_read := 0
	for _, rbyte := range raw {
		output = (output << 7) | uint64(rbyte&0x7F)
		bytes_read += 1
		if (rbyte & 0x80) == 0 {
			return output, bytes_read, nil
		}
	}
	return 0, 0, errors.New("unexpected bytes")
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// MarshalCBOR returns a cbor encoded byte slice of the enterprise address.
func (p *PointerAddress) MarshalCBOR() (bytes []byte, err error) {
	return cbor.Marshal(p.Bytes())
}

// Bytes retuns a byte slice representation of the pointer address.
func (p *PointerAddress) Bytes() (bytes []byte) {
	var buf []byte

	header := 0b0100_0000 | (byte(p.Payment.Kind) << 4) | (byte(p.Network.NetworkId) & 0xF)
	buf = append(buf, header)
	buf = append(buf, p.Payment.Payload...)
	buf = append(buf, VariableNatEncode(p.Stake.Slot)...)
	buf = append(buf, VariableNatEncode(p.Stake.TxIndex)...)

	return append(buf, VariableNatEncode(p.Stake.CertIndex)...)
}

// String returns a bech32 encoded string of the Enterprise Address.
func (p *PointerAddress) String() string {
	str, _ := bech32.Encode(p.Prefix(), p.Bytes())
	return str
}

// Prefix returns the string prefix for the base address. Prefix `addr` for mainnet addresses and `addr_test` for testnet.
func (p *PointerAddress) Prefix() string {
	if p.Network == *network.TestNet() {
		return "addr_test"
	}
	return "addr"
}

// NetworkInfo returns NetworkInfo{ProtocolMagigic and NetworkId}.
func (p *PointerAddress) NetworkInfo() *network.NetworkInfo {
	return &(p.Network)
}

// NewPointer returns a pointer to a new StakePointer given slot, transaction index and certificate index.
func NewPointer(slot, txIndex, certIndex uint64) *StakePointer {
	return &StakePointer{
		Slot:      slot,
		TxIndex:   txIndex,
		CertIndex: certIndex,
	}
}

// NewPointerAddress returns a pointer to a new Pointer address given the network, payment and stake credentials.
func NewPointerAddress(net network.NetworkInfo, payment StakeCredential, stake StakePointer) *PointerAddress {
	return &PointerAddress{
		Network: net,
		Payment: payment,
		Stake:   stake,
	}
}
