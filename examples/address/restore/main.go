// An example showing how to generate base address and enterprise addresses from a
// 24 word mnemonic phrase.
//
// For help of on usage run `go run main.go --help`
//
// Example: To generate an enterprise address run
// 		`go run main.go -mnemonic "YOUR_24_WORD_PHRASE_HERE" -type 0`
// Alternatively replace type with 1 when using this cli to generate base addresses

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/tyler-smith/go-bip39"
)

const (
	enterpriseType uint = iota
	baseType
)

var mnemFlag string
var addrTypeFlag uint
var networkFlag string

func init() {

	flag.StringVar(&mnemFlag, "mnemonic", "", "Mnemonic to restore wallet")
	flag.UintVar(&addrTypeFlag, "type", enterpriseType, "Enum of address type(0: Enterprise Address, 1: Base Address)")
	flag.StringVar(&networkFlag, "network", "mainnet", "The network ie mainnet or testnet")
	flag.Parse()

	if mnemFlag == "" {
		log.Fatal("mnemonic cannot be empty")
	}

	if len(strings.Split(mnemFlag, " ")) != 24 {
		log.Fatal("mnemonic should be 24 words long")
	}

}

func harden(num uint) uint32 {
	return uint32(0x80000000 + num)
}

// checkHandleErr - utility
func checkHandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateEnterprise(net *network.NetworkInfo, entropy []byte) (addr *address.EnterpriseAddress, err error) {
	rootKey := bip32.FromBip39Entropy(
		entropy,
		[]byte{},
	)

	accountKey := rootKey.Derive(harden(1852)).Derive(harden(1815)).Derive(harden(0))

	utxoPubKey := accountKey.Derive(0).Derive(0).Public()
	utxoPubKeyHash := utxoPubKey.PublicKey().Hash()

	addr = address.NewEnterpriseAddress(
		net,
		address.NewKeyStakeCredential(utxoPubKeyHash[:]),
	)
	return
}

func generateBaseAddress(net *network.NetworkInfo, entropy []byte) (addr address.BaseAddress, err error) {
	rootKey := bip32.FromBip39Entropy(
		entropy,
		[]byte{},
	)

	accountKey := rootKey.Derive(harden(1852)).Derive(harden(1815)).Derive(harden(0))

	utxoPubKey := accountKey.Derive(0).Derive(0).Public()
	utxoPubKeyHash := utxoPubKey.PublicKey().Hash()

	stakeKey := accountKey.Derive(2).Derive(0).Public()
	stakeKeyHash := stakeKey.PublicKey().Hash()

	addr = *address.NewBaseAddress(
		net,
		&address.StakeCredential{
			Kind:    address.KeyStakeCredentialType,
			Payload: utxoPubKeyHash[:],
		},
		&address.StakeCredential{
			Kind:    address.KeyStakeCredentialType,
			Payload: stakeKeyHash[:],
		})
	return
}

func main() {

	networks := make(map[string]*network.NetworkInfo)

	networks["mainnet"] = network.MainNet()
	networks["testnet"] = network.TestNet()

	if _, ok := networks[networkFlag]; !ok {
		log.Fatalf("Unsupported network type (%s)", networkFlag)
	}
	net := networks[networkFlag]

	if addrTypeFlag == enterpriseType {
		entropy, err := bip39.EntropyFromMnemonic(mnemFlag)
		checkHandleErr(err)

		addr, err := generateEnterprise(net, entropy)
		checkHandleErr(err)

		fmt.Printf("Enterprise Address: %s", addr.String())

	} else if addrTypeFlag == baseType {
		entropy, err := bip39.EntropyFromMnemonic(mnemFlag)
		checkHandleErr(err)

		addr, err := generateBaseAddress(net, entropy)
		checkHandleErr(err)

		fmt.Printf("Base Address: %s", addr.String())
	}
}
