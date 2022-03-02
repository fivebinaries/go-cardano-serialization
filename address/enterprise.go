package address

import (
	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

type EnterpriseAddress struct {
	Network network.NetworkInfo
	Payment StakeCredential
}

func (e *EnterpriseAddress) Bytes() []byte {
	bytes := make([]byte, 29)
	bytes[0] = 0b01100000 | (byte(e.Payment.Kind) << 4) | (byte(e.Network.NetworkId) & 0xf)
	copy(bytes[1:], e.Payment.Payload[:])
	return bytes
}

func (e *EnterpriseAddress) String() string {
	str, _ := bech32.Encode(e.Prefix(), e.Bytes())
	return str
}

func (e *EnterpriseAddress) NetworkInfo() *network.NetworkInfo {
	return &(e.Network)
}

func (e *EnterpriseAddress) Prefix() string {
	if e.Network.NetworkId == network.TestNet().NetworkId {
		return "addr_test"
	}
	return "addr"
}

func (e *EnterpriseAddress) MarshalCBOR() (bytes []byte, err error) {
	return cbor.Marshal(e.Bytes())
}

func NewEnterpriseAddress(network *network.NetworkInfo, payment *StakeCredential) *EnterpriseAddress {
	return &EnterpriseAddress{
		Network: *network,
		Payment: *payment,
	}
}
