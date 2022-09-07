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
		description  string
		address      string
		addressBytes []byte
		addressHex   string
		err          error
		network      network.NetworkInfo
		addrType     interface{}
	}
	addressesScenarios := []AddressTestCase{
		{
			description:  "valid base address decode and encode",
			address:      "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9v6lkgaj5cunufgvz5f7snql2fawd",
			addressBytes: []byte{0, 50, 85, 4, 149, 218, 98, 148, 141, 123, 151, 29, 32, 139, 154, 147, 190, 231, 139, 99, 163, 238, 124, 147, 94, 6, 51, 127, 237, 100, 165, 176, 54, 33, 77, 129, 83, 96, 162, 42, 93, 35, 79, 149, 154, 253, 145, 217, 83, 28, 159, 18, 134, 10, 137, 244, 38},
			addressHex:   "0032550495da62948d7b971d208b9a93bee78b63a3ee7c935e06337fed64a5b036214d815360a22a5d234f959afd91d9531c9f12860a89f426",
			network:      *network.TestNet(),
			err:          nil,
			addrType:     &address.BaseAddress{},
		},
		{
			description: "invalid base address(invalid/missing checksum)",
			address:     "addr_test1qqe92py4mf3ffrtmjuwjpzu6jwlw0zmr50h8ey67qcehlmty5kcrvg2ds9fkpg32t535l9",
			network:     *network.TestNet(),
			err:         errors.New("invalid checksum"),
			addrType:    nil,
		},
		{
			description:  "valid stake/reward address",
			address:      "stake1u9w862n8jtje5fuc32l20mqqvwaslpveja6paugnnezz99shxsy55",
			addressBytes: []byte{225, 92, 125, 42, 103, 146, 229, 154, 39, 152, 138, 190, 167, 236, 0, 99, 187, 15, 133, 153, 151, 116, 30, 241, 19, 158, 68, 34, 150},
			addressHex:   "e15c7d2a6792e59a27988abea7ec0063bb0f859997741ef1139e442296",
			network:      *network.MainNet(),
			err:          nil,
			addrType:     &address.RewardAddress{},
		},
		{
			description: "invalid stake/reward address",
			address:     "stake1u9w862n8jtje5fuc32l20mqqvwaslpveja6paugnnezz99shxsy",
			network:     *network.MainNet(),
			err:         errors.New("invalid checksum"),
			addrType:    nil,
		},
		{
			description:  "valid enterprise address",
			address:      "addr1vy2qrg3afcprp3lklswy7lux7srcdcd7vghu3md4f0qtd9cszg2k2",
			addressBytes: []byte{97, 20, 1, 162, 61, 78, 2, 48, 199, 246, 252, 28, 79, 127, 134, 244, 7, 134, 225, 190, 98, 47, 200, 237, 181, 75, 192, 182, 151},
			addressHex:   "611401a23d4e0230c7f6fc1c4f7f86f40786e1be622fc8edb54bc0b697",
			network:      *network.MainNet(),
			err:          nil,
			addrType:     &address.EnterpriseAddress{},
		},
		{
			description:  "valid yoroi legacy address",
			address:      "Ae2tdPwUPEZFRbyhz3cpfC2CumGzNkFBN2L42rcUc2yjQpEkxDbkPodpMAi",
			addressHex:   "82d818582183581cba970ad36654d8dd8f74274b733452ddeab9a62a397746be3c42ccdda0001a9026da5b",
			addressBytes: []byte{130, 216, 24, 88, 33, 131, 88, 28, 186, 151, 10, 211, 102, 84, 216, 221, 143, 116, 39, 75, 115, 52, 82, 221, 234, 185, 166, 42, 57, 119, 70, 190, 60, 66, 204, 221, 160, 0, 26, 144, 38, 218, 91},
			network:      *network.MainNet(),
			err:          nil,
			addrType:     &address.ByronAddress{},
		},
		{
			description:  "valid deadulus style legacy address",
			address:      "DdzFFzCqrhsf6zq32tPdqzCqL4JxNSw5aDkiKQp9x8PWUHBXNhR6UNtEeBthFGuf7oSGT2uLKYjoDTyJochABBPCjs6VN4V8eVk7acbe",
			addressHex:   "82d818584283581c0a1e1b7f0e38e24fbe2e30af04c7c7aab10d838cf2bd24b89e81eb12a101581e581c9b1771bd305e4a6a5a37e8a962ed18fb0a339639e039355feecaa7ff001ab79e49ad",
			addressBytes: []byte{130, 216, 24, 88, 66, 131, 88, 28, 10, 30, 27, 127, 14, 56, 226, 79, 190, 46, 48, 175, 4, 199, 199, 170, 177, 13, 131, 140, 242, 189, 36, 184, 158, 129, 235, 18, 161, 1, 88, 30, 88, 28, 155, 23, 113, 189, 48, 94, 74, 106, 90, 55, 232, 169, 98, 237, 24, 251, 10, 51, 150, 57, 224, 57, 53, 95, 238, 202, 167, 255, 0, 26, 183, 158, 73, 173},
			network:      *network.MainNet(),
			err:          nil,
			addrType:     &address.ByronAddress{},
		},
	}

	checkAddrScenario := func(t *testing.T, scenario AddressTestCase, addr address.Address, err error) {
		// TODO: Provide Errors in implementation
		if err != nil {
			assert.Equal(t, err.Error(), scenario.err.Error())
		}

		if addr != nil {
			assert.Equal(t, scenario.address, addr.String())
			assert.Equal(t, scenario.network, *addr.NetworkInfo())
			assert.Equal(t, reflect.TypeOf(addr).String(), reflect.TypeOf(scenario.addrType).String())
		}
	}

	for _, scenario := range addressesScenarios {
		t.Run(scenario.description, func(t *testing.T) {

			// Default
			addr, err := address.NewAddress(scenario.address)
			checkAddrScenario(t, scenario, addr, err)

			// Test from bytes
			if len(scenario.addressBytes) > 0 {
				addr, err := address.NewAddressFromBytes(scenario.addressBytes)
				checkAddrScenario(t, scenario, addr, err)
			}

			// Test from hex
			if scenario.addressHex != "" {
				addr, err := address.NewAddressFromHex(scenario.addressHex)
				checkAddrScenario(t, scenario, addr, err)
			}
		})
	}
}
