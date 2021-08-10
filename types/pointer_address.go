package types

import (
	"github.com/btcsuite/btcutil/bech32"
	"github.com/fivebinaries/go-cardano-serialization/lib"
)

// Pointer implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L609
type Pointer struct {
	Slot      lib.Slot
	TxIndex   TransactionIndex
	CertIndex lib.CertificateIndex
}

// PointerAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L633
type PointerAddress struct {
	Network uint8
	Payment StakeCredential
	Stake   Pointer
}

// NetworkId implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L455
func (p *PointerAddress) NetworkId() (byte, error) {
	return p.Network, nil
}

// ToBech32 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L429
func (p *PointerAddress) ToBech32(prefix *string) (string, error) {
	finalPrefix := ""
	if prefix == nil {
		prefixHeader := "addr"
		prefixTail := ""
		netId, err := p.NetworkId()
		if err != nil {
			return "", err
		}
		if netId == TestNet().NetworkId {
			prefixTail = "_test"
		}
		finalPrefix = prefixHeader + prefixTail
	} else {
		finalPrefix = *prefix
	}
	data, err := bech32.ConvertBits(p.ToBytes(), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(finalPrefix, data)
}

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L282
func (p *PointerAddress) ToBytes() []byte {
	var buf []byte
	header := 0b0100_0000 | (p.Payment.Kind() << 4) | (p.Network & 0xF)
	buf = append(buf, header)
	buf = append(buf, p.Payment.ToRawBytes()...)
	buf = append(buf, VariableNatEncode(uint64(p.Stake.Slot))...)
	buf = append(buf, VariableNatEncode(uint64(p.Stake.TxIndex))...)
	return append(buf, VariableNatEncode(uint64(p.Stake.CertIndex))...)
}

// NewPointerAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L641
func NewPointerAddress(network uint8, payment *StakeCredential, stake *Pointer) PointerAddress {
	return PointerAddress{
		Network: network,
		Payment: *payment,
		Stake:   *stake,
	}
}

// NewEnterpriseAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L525
func NewPointer(slot lib.Slot, txIndex TransactionIndex, certIndex lib.CertificateIndex) *Pointer {
	return &Pointer{
		Slot:      slot,
		TxIndex:   txIndex,
		CertIndex: certIndex,
	}
}
