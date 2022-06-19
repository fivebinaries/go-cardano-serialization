package address_test

import (
	"errors"
	"flag"
	"reflect"
	"testing"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/stretchr/testify/assert"
)

var (
	generate = flag.Bool("gen", false, "generate .golden files")
)

func TestAddressInit(t *testing.T) {
	type AddressTestCase struct {
		description string
		address     string
		err         error
		network     network.NetworkInfo
		addrType    interface{}
	}
	addressesScenarios := []AddressTestCase{
		{
			description: "valid base address decode and encode",
			address:     "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9v6lkgaj5cunufgvz5f7snql2fawd",
			network:     *network.TestNet(),
			err:         nil,
			addrType:    &address.BaseAddress{},
		},
		{
			description: "invalid base address(invalid/missing checksum)",
			address:     "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9",
			network:     *network.TestNet(),
			err:         errors.New("invalid checksum"),
			addrType:    nil,
		},
		{
			description: "valid stake/reward address",
			address:     "stake1u9w862n8jtje5fuc32l20mqqvwaslpveja6paugnnezz99shxsy55",
			network:     *network.MainNet(),
			err:         nil,
			addrType:    &address.RewardAddress{},
		},
		{
			description: "invalid stake/reward address",
			address:     "stake1u9w862n8jtje5fuc32l20mqqvwaslpveja6paugnnezz99shxsy",
			network:     *network.MainNet(),
			err:         errors.New("invalid checksum"),
			addrType:    nil,
		},
		{
			description: "valid enterprise address",
			address:     "addr1vy2qrg3afcprp3lklswy7lux7srcdcd7vghu3md4f0qtd9cszg2k2",
			network:     *network.MainNet(),
			err:         nil,
			addrType:    &address.EnterpriseAddress{},
		},
		{
			description: "valid yoroi legacy address",
			address:     "Ae2tdPwUPEZFRbyhz3cpfC2CumGzNkFBN2L42rcUc2yjQpEkxDbkPodpMAi",
			network:     *network.MainNet(),
			err:         nil,
			addrType:    &address.ByronAddress{},
		},
		{
			description: "valid deadulus style legacy address",
			address:     "DdzFFzCqrhsf6zq32tPdqzCqL4JxNSw5aDkiKQp9x8PWUHBXNhR6UNtEeBthFGuf7oSGT2uLKYjoDTyJochABBPCjs6VN4V8eVk7acbe",
			network:     *network.MainNet(),
			err:         nil,
			addrType:    &address.ByronAddress{},
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
				assert.Equal(t, reflect.TypeOf(addr).String(), reflect.TypeOf(scenario.addrType).String())
			}
		})
	}
}
