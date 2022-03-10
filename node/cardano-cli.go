package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

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

func (cli *cardanoCli) ProtocolParameters() (p *protocol.Protocol, err error) {
	data, err := cli.execCommand("query", "tip")

	if err := json.Unmarshal(data, p); err != nil {
		return p, err
	}
	return
}

func (cli *cardanoCli) UTXOs(address.Address) (txIs []tx.TxInput, err error) {
	return
}

func (cli *cardanoCli) QueryTip() (tip NetworkTip, err error) {
	return
}

func (cli *cardanoCli) SubmitTx(tx.Tx) (err error) {
	return
}

// NewCardanoCliNode returns a wrapper for the cardano-cli with the Node interface
func NewCardanoCliNode(network *network.NetworkInfo, cliPaths ...string) *cardanoCli {
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
