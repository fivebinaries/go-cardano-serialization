package address

import "github.com/btcsuite/btcutil/bech32"

// EnterpriseAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L519
type EnterpriseAddress struct {
	Network uint8
	Payment StakeCredential
}

// NetworkId implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L455
func (e EnterpriseAddress) NetworkId() (byte, error) {
	return e.Network, nil
}

// ToBech32 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L429
func (e EnterpriseAddress) ToBech32(prefix *string) (string, error) {
	finalPrefix := ""
	if prefix == nil {
		prefixHeader := "addr"
		prefixTail := ""
		netId, err := e.NetworkId()
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
	data, err := bech32.ConvertBits(e.ToBytes(), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(finalPrefix, data)
}

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L292
func (e EnterpriseAddress) ToBytes() []byte {
	var buf []byte
	header := 0b0110_0000 | (e.Payment.Kind() << 4) | (e.Network & 0xF)
	buf = append(buf, header)
	return append(buf, e.Payment.ToRawBytes()...)
}

// NewEnterpriseAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L525
func NewEnterpriseAddress(network uint8, payment *StakeCredential) EnterpriseAddress {
	return EnterpriseAddress{
		Network: network,
		Payment: *payment,
	}
}
