package address

import (
	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

// RewardAddress contains content of the reward/staking address.
// Reward account addresses are used to distribute rewards for participating in the proof-of-stake protocol (either directly or via delegation).
type RewardAddress struct {
	Network network.NetworkInfo
	Stake   StakeCredential
}

// Bytes returns a  length 29 byte slice represantation of the Address.
func (r *RewardAddress) Bytes() []byte {
	data := make([]byte, 29)
	data[0] = 0b1110_0000 | (byte(r.Stake.Kind) << 4) | (byte(r.Network.NetworkId) & 0xf)
	copy(data[1:], r.Stake.Payload[:])
	return data
}

// String returns a bech32 encoded string of the reward address.
func (r *RewardAddress) String() string {
	str, _ := bech32.Encode(r.Prefix(), r.Bytes())

	return str
}

// MarshalCBOR returns a cbor encoded byte slice of the base address.
func (r *RewardAddress) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(r.Bytes)
}

// NetworkInfo returns pointer to NetworkInfo{ProtocolMagigic and NetworkId}.
func (r *RewardAddress) NetworkInfo() *network.NetworkInfo {
	return &(r.Network)
}

// Prefix returns the string prefix for the base address. Prefix `stake` for mainnet addresses and `stake_test` for testnet.
func (r *RewardAddress) Prefix() string {
	if r.Network.NetworkId == network.TestNet().NetworkId {
		return "stake_test"
	}
	return "stake"
}

// NewRewardAddress returns a pointer to a new RewardAddress given the network and stake credentials.
func NewRewardAddress(net *network.NetworkInfo, stake *StakeCredential) *RewardAddress {
	return &RewardAddress{
		Network: *net,
		Stake:   *stake,
	}
}
