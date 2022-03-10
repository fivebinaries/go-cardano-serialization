package tx_test

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
	"github.com/fivebinaries/go-cardano-serialization/tx"
	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _  = runtime.Caller(0)
	packagepath = filepath.Dir(b)
)

var (
	// update   = flag.Bool("update", false, "update .golden files")
	generate = flag.Bool("gen", false, "generate .golden files")
)

type utxoIn struct {
	TxHash         string `json:"txHash"`
	TxIndex        uint   `json:"txIndex"`
	AmountLovelace uint   `json:"amountLovelace"`
}

type txDetails struct {
	ReceiverAddress string `json:"receiverAddress"`
	AmountLovelace  uint   `json:"amountLovelace"`
	ChangeAddress   string `json:"changeAddress"`
	SlotNo          uint   `json:"slotNo"`
	UtxoIn          utxoIn `json:"utxoIn"`
}

type txScenario struct {
	Description string
	GoldenFile  string
	AddrProc    func(*address.BaseAddress) address.Address
}

func createRootKey() bip32.XPrv {
	rootKey := bip32.FromBip39Entropy(
		[]byte{214, 64, 138, 69, 145, 210, 32, 51, 202, 45, 90, 151, 33, 194, 153, 176, 188, 94, 94, 186, 67, 118, 194, 227, 207, 157, 54, 49, 34, 12, 83, 93},
		[]byte{},
	)
	return rootKey
}

func harden(num uint) uint32 {
	return uint32(0x80000000 + num)
}

func generateBaseAddress(net *network.NetworkInfo) (addr *address.BaseAddress, utxoPrvKey bip32.XPrv, err error) {
	rootKey := createRootKey()
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

func getTxDetails(fp string) (txD txDetails) {
	data, err := readJson(fp)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &txD); err != nil {
		log.Fatal("Err", err)
	}

	return

}

func readJson(fp string) (data []byte, err error) {
	file, err := os.Open(fp)
	if err != nil {
		return
	}

	defer file.Close()

	data, err = ioutil.ReadAll(file)
	if err != nil {
		return
	}

	return
}

func TestTxBuilderRaw(t *testing.T) {
	createRootKey()
	basepath := packagepath[:strings.LastIndex(packagepath, "/")]
	pr, err := protocol.LoadProtocol(filepath.Join(basepath, "testdata", "protocol", "protocol.json"))
	if err != nil {
		log.Fatal(err)
	}

	txD := getTxDetails(filepath.Join(basepath, "testdata", "transaction", "tx_builder", "json", "raw_tx.json"))

	txScenarios := []txScenario{
		{
			Description: "Transaction with base address marshalling",
			GoldenFile:  "raw_tx_base.golden",
			AddrProc:    func(addr *address.BaseAddress) address.Address { return addr },
		},
		{
			Description: "Transaction with enterprise address marshalling",
			GoldenFile:  "raw_tx_ent.golden",
			AddrProc:    func(addr *address.BaseAddress) address.Address { return addr.ToEnterprise() },
		},
	}

	for _, sc := range txScenarios {
		t.Run(sc.Description, func(t *testing.T) {
			builder := tx.NewTxBuilder(
				*pr,
				[]bip32.XPrv{},
			)
			addr, utxoPrv, err := generateBaseAddress(network.MainNet())
			if err != nil {
				log.Fatal(err)
			}

			builder.AddInputs(
				tx.NewTxInput(
					txD.UtxoIn.TxHash,
					uint16(txD.UtxoIn.TxIndex),
					txD.UtxoIn.AmountLovelace,
				),
			)

			builder.AddOutputs(
				tx.NewTxOutput(
					sc.AddrProc(addr),
					txD.AmountLovelace,
				),
			)

			changeAddr, err := address.NewAddress(txD.ChangeAddress)
			if err != nil {
				log.Fatal(err)
			}
			builder.SetTTL(uint32(txD.SlotNo))
			builder.AddChangeIfNeeded(changeAddr)

			builder.Sign(
				utxoPrv,
			)
			txFinal, err := builder.Build()
			if err != nil {
				t.Fatal(err)
			}

			txHex, err := txFinal.Hex()
			if err != nil {
				log.Fatal(err)
			}
			golden := filepath.Join(basepath, "testdata", "transaction", "tx_builder", "golden", sc.GoldenFile)
			assert.Equal(t, txHex, ReadOrGenerateGoldenFile(t, golden, txFinal))
		})
	}
}

func WriteGoldenFile(t *testing.T, path string, data []byte) {
	t.Helper()
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		t.Fatal(err)
	}
}

func ReadOrGenerateGoldenFile(t *testing.T, path string, txF tx.Tx) string {
	t.Helper()
	b, err := ioutil.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if *generate {
			txHex, err := txF.Hex()
			if err != nil {
				t.Fatal("golden-gen: Failed to hex encode transaction")
			}
			if err != nil {
				t.Fatal("golden-gen: Failed to hex encode transaction")
			}
			WriteGoldenFile(t, path, []byte(txHex))
			return txHex
		}
		t.Fatalf("golden-read: Missing golden file. Run `go test -args -gen` to generate it.")
	case err != nil:
		t.Fatal(err)
	}
	return string(b)
}
