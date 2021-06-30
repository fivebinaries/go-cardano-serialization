package address

import (
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fxamacker/cbor/v2"
	"hash/crc32"
	"reflect"
)

// AddrType implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L43
type AddrType uint64
type AddrAttributesRaw map[uint64][]byte
type PayloadBytes []byte

const (
	ATPubKey AddrType = 0
	ATScript AddrType = 1
	ATRedeem AddrType = 2

	ATTRIBUTE_NAME_TAG_DERIVATION     uint64 = 1
	ATTRIBUTE_NAME_TAG_PROTOCOL_MAGIC uint64 = 2

	SPENDING_DATA_TAG_PUBKEY uint64 = 0

	CBORInt8  = 24
	CBORInt16 = 25
	CBORInt32 = 26
	CBORInt64 = 27
)

// Addr implement https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L194
type Addr struct {
	_       interface{} `cbor:",toarray"`
	Payload PayloadBytes
	Crc     uint32
}

// NetworkInfo implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L34
type NetworkInfo struct {
	NetworkId     uint8
	ProtocolMagic uint32
}

type Address interface {
	ToBytes() []byte
	NetworkId() (byte, error)
	ToBech32(prefix *string) (string, error)
}

// TestNet implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L53
func TestNet() NetworkInfo {
	return NetworkInfo{
		NetworkId:     0b0000,
		ProtocolMagic: 1097911063,
	}
}

// MainNet implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L59
func MainNet() NetworkInfo {
	return NetworkInfo{
		NetworkId:     0b0001,
		ProtocolMagic: 764824073,
	}
}

func (addr *Addr) ToBytes() ([]byte, error) {
	return cbor.Marshal(addr)
}

func (addr *Addr) CalcCRC() uint32 {
	addr.Crc = crc32.ChecksumIEEE(addr.Payload)
	return addr.Crc
}

func (addr *Addr) ToString() (string, error) {
	addrBytes, err := addr.ToBytes()
	if err != nil {
		return "", err
	}
	return base58.Encode(addrBytes), nil
}

// Reverse any slice
func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// VariableNatEncode implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L21
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

// VariableNatEncode implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L8
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

// AddressFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L313
func AddressFromBytes(data []byte) (Address, error) {
	// header has 4 bits addr type discrim then 4 bits network discrim.
	// Copied from shelley.cddl:
	//
	// shelley payment addresses:
	// bit 7: 0
	// bit 6: base/other
	// bit 5: pointer/enterprise [for base: stake cred is keyhash/scripthash]
	// bit 4: payment cred is keyhash/scripthash
	// bits 3-0: network id
	//
	// reward addresses:
	// bits 7-5: 111
	// bit 4: credential is keyhash/scripthash
	// bits 3-0: network id
	//
	// byron addresses:
	// bits 7-4: 1000

	header := data[0]
	network := header & 0x0F
	const hashLen = crypto.Ed25519KeyHashLen
	switch (header & 0xF0) >> 4 {
	// base
	case 0b0000, 0b0001, 0b0010, 0b0011:
		baseAddrSize := 1 + hashLen*2
		if len(data) < baseAddrSize {
			return nil, errors.New("cbor not enough error")
		}
		if len(data) > baseAddrSize {
			return nil, errors.New("cbor trailing data error")
		}
		return NewBaseAddress(network, readAddrCred(data, header, 4, 1), readAddrCred(data, header, 5, 1+hashLen)), nil
	// pointer
	case 0b0100, 0b0101:
		// header + keyhash + 3 natural numbers (min 1 byte each)
		const ptrAddrMinSize = 1 + hashLen + 1 + 1 + 1
		if len(data) < ptrAddrMinSize {
			return nil, errors.New("cbor not enough error")
		}
		byteIndex := 1
		paymentCred := readAddrCred(data, header, 4, 1)
		byteIndex += hashLen
		slot, slot_bytes, err := VariableNatDecode(data[byteIndex:])
		if err != nil {
			return nil, err
		}
		byteIndex += slot_bytes
		txIndex, txBytes, err := VariableNatDecode(data[byteIndex:])
		if err != nil {
			return nil, err
		}
		byteIndex += txBytes
		certIndex, certBytes, err := VariableNatDecode(data[byteIndex:])
		if err != nil {
			return nil, err
		}
		byteIndex += certBytes

		if byteIndex < len(data) {
			return nil, errors.New("cbor trailing data error")
		}

		return NewPointerAddress(network, paymentCred,
			NewPointer(Slot(slot), TransactionIndex(txIndex), CertificateIndex(certIndex)),
		), nil

	// enterprise
	case 0b0110, 0b0111:
		const enterpriseAddrSize = 1 + hashLen
		if len(data) < enterpriseAddrSize {
			return nil, errors.New("cbor not enough error")
		}
		if len(data) > enterpriseAddrSize {
			return nil, errors.New("cbor trailing data error")
		}
		return NewEnterpriseAddress(network, readAddrCred(data, header, 4, 1)), nil

	// reward
	case 0b1110, 0b1111:
		const rewardAddrSize = 1 + hashLen
		if len(data) < rewardAddrSize {
			return nil, errors.New("cbor not enough error")
		}
		if len(data) > rewardAddrSize {
			return nil, errors.New("cbor trailing data error")
		}
		return NewRewardAddress(network, readAddrCred(data, header, 4, 1)), nil

	// byron
	case 0b1000:
		// note: 0b1000 was chosen because all existing Byron addresses actually start with 0b1000
		// Therefore you can re-use Byron addresses as-is

		addr, err := FromBytes(data)
		if err != nil {
			return nil, err
		}
		return addr.ToByronAddress(), nil
	}
	return nil, errors.New("bad address type")
}

// AddressFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L449
func AddressFromBech32(bechStr string) (Address, error) {
	_, u5data, err := bech32.Decode(bechStr)
	if err != nil {
		return nil, err
	}
	data, err := bech32.ConvertBits(u5data, 5, 8, false)
	if err != nil {
		return nil, err
	}
	return AddressFromBytes(data)
}
