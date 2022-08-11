package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
	"github.com/fivebinaries/go-cardano-serialization/tx"
)

type cardanoCli struct {
	cliPath string
	network *network.NetworkInfo
}

func (cli *cardanoCli) execCommand(args ...string) (data []byte, err error) {
	buf := &bytes.Buffer{}
	var cmdSuffix string

	if cli.network.NetworkId == 0 {
		cmdSuffix = "--mainnet"
	} else {
		cmdSuffix = fmt.Sprintf("--testnet-magic %d", cli.network.ProtocolMagic)
	}

	args = append(args, cmdSuffix)
	cmd := exec.Command(cli.cliPath, args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

func (cli *cardanoCli) ProtocolParameters() (p protocol.Protocol, err error) {
	data, err := cli.execCommand("query", "protocol-parameters")

	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	return
}

func (cli *cardanoCli) UTXOs(addr address.Address) (txIs []tx.TxInput, err error) {
	data, err := cli.execCommand("query", "utxos", fmt.Sprintf("--address %s", addr))
	if err != nil {
		return
	}

	ldata := strings.Split(string(data), "\n")
	lenData := len(data)

	if lenData < 3 {
		return
	}

	for _, it := range ldata[2 : lenData-1] {
		sec := strings.Fields(it)
		txIx, err := strconv.Atoi(sec[1])
		if err != nil {
			return txIs, err
		}

		amSec := strings.Fields(sec[2])

		amount, err := strconv.ParseUint(amSec[0], 10, 16)
		if err != nil {
			return txIs, err
		}

		utxo := *tx.NewTxInput(
			sec[0],
			uint16(txIx),
			uint(amount),
		)

		txIs = append(txIs, utxo)
	}
	return
}

func (cli *cardanoCli) QueryTip() (tip NetworkTip, err error) {
	data, err := cli.execCommand("query", "tip")

	if err := json.Unmarshal(data, &tip); err != nil {
		return tip, err
	}
	return
}

func (cli *cardanoCli) SubmitTx(txFinal tx.Tx) (txHash string, err error) {
	type sb struct {
		TxType      string `json:"type"`
		Description string `json:"description"`
		CborHex     string `json:"cborHex"`
	}

	txHex, err := txFinal.Hex()
	if err != nil {
		return
	}

	outTx := sb{
		CborHex: txHex,
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "gada-tx-")
	if err != nil {
		return
	}

	defer os.Remove(tmpFile.Name())

	data, err := json.MarshalIndent(outTx, "", "	")
	if err != nil {
		return
	}

	if _, err := tmpFile.Write((data)); err != nil {
		return txHash, err
	}

	cliData, err := cli.execCommand("submit", fmt.Sprintf("--tx-file %s", tmpFile.Name()))
	if err != nil {
		return
	}

	txHash = string(cliData)

	if err := tmpFile.Close(); err != nil {
		return txHash, err
	}
	return
}

// NewCardanoCliNode returns a wrapper for the cardano-cli with the Node interface
func NewCardanoCliNode(network *network.NetworkInfo, cliPaths ...string) Node {
	var cliPath string

	if len(cliPaths) > 0 {
		if cliPaths[0] != "" {
			cliPath = cliPaths[0]
		}
	} else {
		cliPath = "cardano-cli"
	}
	return &cardanoCli{
		cliPath: cliPath,
		network: network,
	}
}
