package address

import (
	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"

	"github.com/fxamacker/cbor/v2"
)

// BaseAddress contains information of the base address.
// A base address directly specifies the staking key that should control the stake for that address but can be used for transactions without registering the staking key in advance.
type BaseAddress struct {
	Network network.NetworkInfo
	Payment StakeCredential
	Stake   StakeCredential
}

// Bytes returns a  length 57 byte slice represantation of the Address.
func (b *BaseAddress) Bytes() []byte {
	bytes := make([]byte, 57)
	bytes[0] = (byte(b.Payment.Kind) << 4) | (byte(b.Stake.Kind) << 5) | (byte(b.Network.NetworkId) & 0xf)
	copy(bytes[1:29], b.Payment.Payload[:])
	copy(bytes[29:], b.Stake.Payload[:])
	return bytes
}

// String returns a bech32 encoded string of the Enterprise Address.
func (b *BaseAddress) String() string {
	str, _ := bech32.Encode(b.Prefix(), b.Bytes())

	return str
}

// DecodeAddress decodes a  bech32-encoded string to []byte
func DecodeBech32Address(addr string) (string, []byte, error) {
	return bech32.Decode(addr)
}

// NetworkInfo returns NetworkInfo{ProtocolMagigic and NetworkId}.
func (b *BaseAddress) NetworkInfo() *network.NetworkInfo {
	return &(b.Network)
}

// Prefix returns the string prefix for the base address. Prefix `addr` for mainnet addresses and `addr_test` for testnet.
func (b *BaseAddress) Prefix() string {
	if b.Network.NetworkId == network.TestNet().NetworkId {
		return "addr_test"
	}
	return "addr"
}

// MarshalCBOR returns a cbor encoded byte slice of the base address.
func (b *BaseAddress) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(b.Bytes())
}

// ToEnterprise returns the Enterprise Address from the base address. This is equivalent to removing the stake part of the base address.
func (b *BaseAddress) ToEnterprise() (addr *EnterpriseAddress) {
	addr = NewEnterpriseAddress(
		&b.Network,
		&b.Payment,
	)
	return
}

// ToReward returns the RewardAddress from the base address. This is equivalent to removing the payment part of the base address.
func (b *BaseAddress) ToReward() (addr *RewardAddress) {
	addr = NewRewardAddress(
		&b.Network,
		&b.Stake,
	)
	return
}

// NewBaseAddress returns a pointer to a new BaseAddress given the network, payment and stake credentials.
func NewBaseAddress(network *network.NetworkInfo, payment *StakeCredential, stake *StakeCredential) *BaseAddress {
	return &BaseAddress{
		Network: *network,
		Payment: *payment,
		Stake:   *stake,
	}
}
