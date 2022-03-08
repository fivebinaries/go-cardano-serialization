package address

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

type ByronAddressAttributes struct {
	Payload []byte `cbor:"1,keyasint,omitempty"`
	Network *uint8 `cbor:"2,keyasint,omitempty"`
}

type ByronAddress struct {
	Hash       []byte
	Attributes ByronAddressAttributes
}

func (b *ByronAddress) Bytes() (bytes []byte) {
	bytes, _ = b.MarshalCBOR()
	return bytes
}

func (b *ByronAddress) String() (str string) {
	return base58.Encode(b.Bytes())
}

func (b *ByronAddress) NetworkInfo() (ni *network.NetworkInfo) {
	if b.Attributes.Network == nil {
		return network.MainNet()
	}
	return network.TestNet()
}

func (b *ByronAddress) MarshalCBOR() (bytes []byte, err error) {
	return cbor.Marshal(b.Bytes())
}

func (b *ByronAddress) Prefix() string {
	return ""
}
