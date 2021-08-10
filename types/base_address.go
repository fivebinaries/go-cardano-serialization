package types

import "github.com/btcsuite/btcutil/bech32"

// BaseAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L480
type BaseAddress struct {
	Network uint8
	Payment StakeCredential
	Stake   StakeCredential
}

// NetworkId implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L455
func (b *BaseAddress) NetworkId() (byte, error) {
	return b.Network, nil
}

// ToBech32 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L429
func (b *BaseAddress) ToBech32(prefix *string) (string, error) {
	finalPrefix := ""
	if prefix == nil {
		prefixHeader := "addr"
		prefixTail := ""
		netId, err := b.NetworkId()
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
	data, err := bech32.ConvertBits(b.ToBytes(), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(finalPrefix, data)
}

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L274
func (b *BaseAddress) ToBytes() []byte {
	var buf []byte
	header := (b.Payment.Kind() << 4) | (b.Stake.Kind() << 5) | (b.Network & 0xF)
	buf = append(buf, header)
	buf = append(buf, b.Payment.ToRawBytes()...)
	return append(buf, b.Stake.ToRawBytes()...)
}

// NewBaseAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L487
func NewBaseAddress(network uint8, payment *StakeCredential, stake *StakeCredential) BaseAddress {
	return BaseAddress{
		Network: network,
		Payment: *payment,
		Stake:   *stake,
	}
}
