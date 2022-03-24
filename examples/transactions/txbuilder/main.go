package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/node"
	"github.com/fivebinaries/go-cardano-serialization/tx"
	"github.com/tyler-smith/go-bip39"
)

const (
	enterpriseType uint = iota
	baseType
)

var mnemFlag string
var addrTypeFlag uint
var networkFlag string
var receiverAddr string

func init() {

	flag.StringVar(&mnemFlag, "mnemonic", "", "Mnemonic to restore wallet")
	flag.UintVar(&addrTypeFlag, "type", enterpriseType, "Enum of address type(0: Enterprise Address, 1: Base Address)")
	flag.StringVar(&networkFlag, "network", "mainnet", "The network ie mainnet or testnet")
	flag.StringVar(&receiverAddr, "sendTo", "", "The address to send ada to")
	flag.Parse()

	if mnemFlag == "" {
		log.Fatal("mnemonic cannot be empty")
	}

	if mlen := len(strings.Split(mnemFlag, " ")); mlen != 24 {
		log.Fatalf("mnemonic should be 24 words long, found {%d}", mlen)
	}

}

// createRootKey returns a bip32 private key generated from 24 word (bip39) secret.
func createRootKey(mnemonic string) bip32.XPrv {
	entropy, _ := bip39.EntropyFromMnemonic(mnemonic)
	rootKey := bip32.FromBip39Entropy(
		entropy,
		[]byte{},
	)
	return rootKey
}

func harden(num uint) uint32 {
	return uint32(0x80000000 + num)
}

func generateBaseAddress(net *network.NetworkInfo, rootKey bip32.XPrv) (addr *address.BaseAddress, utxoPrvKey bip32.XPrv, err error) {
	accountKey := rootKey.Derive(harden(1852)).Derive(harden(1815)).Derive(harden(0))

	utxoPrvKey = accountKey.Derive(0).Derive(0)
	utxoPubKey := utxoPrvKey.Public()
	utxoPubKeyHash := utxoPubKey.PublicKey().Hash()

	stakeKey := accountKey.Derive(2).Derive(0).Public()
	stakeKeyHash := stakeKey.PublicKey().Hash()

	addr = address.NewBaseAddress(
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
	cli := node.NewBlockfrostClient(
		os.Getenv("BLOCKFROST_PROJECT_ID"),
		network.TestNet(),
	)

	// get protocol parameters for linearfee formula
	// fee = TxFeeFixed + TxFeePerByte * byteCount
	pr, err := cli.ProtocolParameters()
	if err != nil {
		log.Fatal(err)
	}

	rootKey := createRootKey(mnemFlag)

	// Generate a base address on testnet from the rootkey
	// the utxoPrvKey is used to sign the transaction
	sourceAddr, utxoPrvKey, err := generateBaseAddress(network.TestNet(), rootKey)
	if err != nil {
		log.Fatal(err)
	}

	// If no sendTo address is provided use the source address
	if receiverAddr == "" {
		receiverAddr = sourceAddr.String()
	}

	receiver, err := address.NewAddress(receiverAddr)
	if err != nil {
		log.Fatal("Address create:", err)
	}

	utxos, err := cli.UTXOs(sourceAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Pick an Unspent Transaction Output as input for the transaction
	// For this example we'll pick the first utxo that has enough ADA for the
	// outputs and fee.

	builder := tx.NewTxBuilder(
		pr,
		[]bip32.XPrv{utxoPrvKey},
	)

	// Send 5000000 lovelace or 5 ADA
	sendAmount := 5000000
	var firstMatchInput tx.TxInput

	// Loop through utxos to find first input with enough ADA
	for _, utxo := range utxos {
		minRequired := sendAmount + 1000000 + 200000
		if utxo.Amount >= uint(minRequired) {
			firstMatchInput = utxo
		}
	}

	// Add the transaction Input / UTXO
	builder.AddInputs(&firstMatchInput)

	// Add a transaction output with the receiver's address and amount of 5 ADA
	builder.AddOutputs(tx.NewTxOutput(
		receiver,
		uint(sendAmount),
	))

	// Query tip from a node on the network. This is to get the current slot
	// and compute TTL of transaction.
	tip, err := cli.QueryTip()
	if err != nil {
		log.Fatal(err)
	}

	// Set TTL for 5 min into the future
	builder.SetTTL(uint32(tip.Slot) + uint32(300))

	// Route back the change to the source address
	// This is equivalent to adding an output with the source address and change amount
	builder.AddChangeIfNeeded(sourceAddr)

	// Build loops through the witness private keys and signs the transaction body hash
	txFinal, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	txHash, err := cli.SubmitTx(txFinal)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(txHash)
}
