package address_test

import (
	"errors"
	"testing"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/stretchr/testify/assert"
)

func TestAddressInit(t *testing.T) {
	type AddressTestCase struct {
		description string
		address     string
		err         error
		network     network.NetworkInfo
	}
	addressesScenarios := []AddressTestCase{
		{
			description: "valid base address decode and encode",
			address:     "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9v6lkgaj5cunufgvz5f7snql2fawd",
			network:     *network.TestNet(),
			err:         nil,
		},
		{
			description: "invalid base address(invalid/missing checksum)",
			address:     "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9",
			network:     *network.TestNet(),
			err:         errors.New("invalid checksum"),
		},
	}

	for _, scenario := range addressesScenarios {
		t.Run(scenario.description, func(t *testing.T) {
			addr, err := address.NewAddress(scenario.address)

			// TODO: Provide Errors in implementation
			if err != nil {
				assert.Equal(t, err.Error(), scenario.err.Error())
			}

			if addr != nil {
				assert.Equal(t, scenario.address, addr.String())
				assert.Equal(t, scenario.network, *addr.NetworkInfo())
			}
		})
	}
}
