package address

import (
	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

// EnterpriseAddress contains content for enterprise addresses.
// Enterprise addresses carry no stake rights, so using these addresses means that you are opting out of participation in the proof-of-stake protocol.
type EnterpriseAddress struct {
	Network network.NetworkInfo
	Payment StakeCredential
}

// Bytes returns a  length 29 byte slice represantation of the Address.
func (e *EnterpriseAddress) Bytes() []byte {
	bytes := make([]byte, 29)
	bytes[0] = 0b01100000 | (byte(e.Payment.Kind) << 4) | (byte(e.Network.NetworkId) & 0xf)
	copy(bytes[1:], e.Payment.Payload[:])
	return bytes
}

// String returns a bech32 encoded string of the Enterprise Address.
func (e *EnterpriseAddress) String() string {
	str, _ := bech32.Encode(e.Prefix(), e.Bytes())
	return str
}

// NetworkInfo returns NetworkInfo{ProtocolMagigic and NetworkId}.
func (e *EnterpriseAddress) NetworkInfo() *network.NetworkInfo {
	return &(e.Network)
}

// Prefix returns the string prefix for the Enterprise Address. Prefix `addr` for mainnet addresses and `addr_test` for testnet.
func (e *EnterpriseAddress) Prefix() string {
	if e.Network.NetworkId == network.TestNet().NetworkId {
		return "addr_test"
	}
	return "addr"
}

// MarshalCBOR returns a cbor encoded byte slice of the enterprise address.
func (e *EnterpriseAddress) MarshalCBOR() (bytes []byte, err error) {
	return cbor.Marshal(e.Bytes())
}

// NewEnterpriseAddress returns a pointer to a new Enterprise Address given the network and payment.
func NewEnterpriseAddress(network *network.NetworkInfo, payment *StakeCredential) *EnterpriseAddress {
	return &EnterpriseAddress{
		Network: *network,
		Payment: *payment,
	}
}
