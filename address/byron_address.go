package address

import (
	"errors"
	"github.com/btcsuite/btcutil/bech32"
	"log"
)

// ByronAddress implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L173
type ByronAddress struct {
	ExtendedAddr
}

// ToBech32 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L429
func (b ByronAddress) ToBech32(prefix *string) (string, error) {
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

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L306
func (b ByronAddress) ToBytes() []byte {
	addr, err := b.ToAddr()
	if err != nil {
		log.Panic(err)
	}
	buf, err := addr.ToBytes()
	if err != nil {
		log.Panic(err)
	}
	return buf
}

// ProtocolMagic implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L191
func (b *ByronAddress) ProtocolMagic() uint32 {
	if b.Attributes.ProtocolMagic != nil {
		return *b.Attributes.ProtocolMagic
	}
	return MainNet().ProtocolMagic
}

// NetworkId implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L202
func (b ByronAddress) NetworkId() (uint8, error) {
	// premise: during the Byron-era, we had one mainnet (764824073) and many many testnets
	// with each testnet getting a different protocol magic
	// in Shelley, this changes so that:
	// 1) all testnets use the same u8 protocol magic
	// 2) mainnet is re-mapped to a single u8 protocol magic

	// recall: in Byron mainnet, the network_id is omitted from the address to save a few bytes
	// so here we return the mainnet id if none is found in the address
	protocolMagic := b.ProtocolMagic()
	mainNet := MainNet()
	testNet := TestNet()
	if protocolMagic == mainNet.ProtocolMagic {
		return mainNet.NetworkId, nil
	}
	if protocolMagic == testNet.ProtocolMagic {
		return testNet.NetworkId, nil
	}
	return 0, errors.New("unexpected protocol magic")
}
