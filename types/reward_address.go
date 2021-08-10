package types

import (
	"github.com/btcsuite/btcutil/bech32"
)

// RewardAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L551
type RewardAddress struct {
	Network uint8
	Payment StakeCredential
}

// NetworkId implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L455
func (r *RewardAddress) NetworkId() (byte, error) {
	return r.Network, nil
}

// ToBech32 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L429
func (r *RewardAddress) ToBech32(prefix *string) (string, error) {
	finalPrefix := ""
	if prefix == nil {
		prefixHeader := "stake"
		prefixTail := ""
		netId, err := r.NetworkId()
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
	data, err := bech32.ConvertBits(r.ToBytes(), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(finalPrefix, data)
}

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L299
func (r *RewardAddress) ToBytes() []byte {
	var buf []byte
	header := 0b1110_0000 | (r.Payment.Kind() << 4) | (r.Network & 0xF)
	buf = append(buf, header)
	return append(buf, r.Payment.ToRawBytes()...)
}

// NewRewardAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L558
func NewRewardAddress(network uint8, payment *StakeCredential) RewardAddress {
	return RewardAddress{
		Network: network,
		Payment: *payment,
	}
}
