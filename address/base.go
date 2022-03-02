package address

import (
	"log"

	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
	"github.com/fivebinaries/go-cardano-serialization/network"

	"github.com/fxamacker/cbor/v2"
)

type BaseAddress struct {
	Network network.NetworkInfo
	Payment StakeCredential
	Stake   StakeCredential
}

func (b *BaseAddress) Bytes() []byte {
	bytes := make([]byte, 57)
	bytes[0] = (byte(b.Payment.Kind) << 4) | (byte(b.Stake.Kind) << 5) | (byte(b.Network.NetworkId) & 0xf)
	copy(bytes[1:29], b.Payment.Payload[:])
	copy(bytes[29:], b.Stake.Payload[:])
	return bytes
}

func (b *BaseAddress) String() string {
	str, err := bech32.Encode(b.Prefix(), b.Bytes())

	if err != nil {
		log.Println(err)
	}
	return str
}

func (b *BaseAddress) NetworkInfo() *network.NetworkInfo {
	return &(b.Network)
}

func (b *BaseAddress) Prefix() string {
	if b.Network.NetworkId == network.TestNet().NetworkId {
		return "addr_test"
	}
	return "addr"
}

func (b *BaseAddress) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(b.Bytes())
}

func (b *BaseAddress) ToEnterprise() (addr *EnterpriseAddress, err error) {
	addr = NewEnterpriseAddress(
		&b.Network,
		&b.Payment,
	)
	return
}

func NewBaseAddress(network *network.NetworkInfo, payment *StakeCredential, stake *StakeCredential) *BaseAddress {
	return &BaseAddress{
		Network: *network,
		Payment: *payment,
		Stake:   *stake,
	}
}
