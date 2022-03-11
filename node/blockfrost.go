package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/blockfrost/blockfrost-go"
	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
	"github.com/fivebinaries/go-cardano-serialization/tx"
)

type blockfrostNode struct {
	network   *network.NetworkInfo
	client    blockfrost.APIClient
	projectId string
}

func (b *blockfrostNode) getNetwork() (serverUrl string) {
	if b.network.NetworkId == 0 {
		serverUrl = blockfrost.CardanoTestNet
	} else {
		serverUrl = blockfrost.CardanoMainNet
	}

	return
}

// UTXOs queries the network for Unspent Transaction Outputs belonging to an address.
func (b *blockfrostNode) UTXOs(addr address.Address) (txIs []tx.TxInput, err error) {
	utxos, err := b.client.AddressUTXOs(
		context.TODO(),
		addr.String(),
		blockfrost.APIQueryParams{},
	)
	if err != nil {
		return
	}

	for _, utxo := range utxos {
		var amount uint
		for _, am := range utxo.Amount {
			if am.Unit == "lovelace" {

				amountI, err := strconv.Atoi(am.Quantity)
				if err != nil {
					return []tx.TxInput{}, err
				}
				amount = uint(amountI)
			}

		}
		txIs = append(txIs, *tx.NewTxInput(utxo.TxHash, uint16(utxo.OutputIndex), amount))
	}

	return
}

// ProtocolParameters queries the protocol parameters of the network.
func (b *blockfrostNode) ProtocolParameters() (p protocol.Protocol, err error) {
	params, err := b.client.LatestEpochParameters(context.TODO())
	if err != nil {
		return
	}

	return protocol.Protocol{
		TxFeePerByte: uint(params.MinFeeA),
		TxFeeFixed:   uint(params.MinFeeB),
		MaxTxSize:    uint(params.MaxTxSize),
		ProtocolVersion: protocol.ProtocolVersion{
			uint8(params.ProtocolMajorVer),
			uint8(params.ProtocolMinorVer),
		},
	}, nil
}

// SubmitTx submits a signed transaction to the network and returns the transaction hash or error
func (b *blockfrostNode) SubmitTx(txFinal tx.Tx) (txHash string, err error) {
	txB, err := txFinal.Bytes()
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		fmt.Sprintf("%s/%s", b.getNetwork(), "tx/submit"),
		bytes.NewReader(txB),
	)
	if err != nil {
		return
	}

	req.Header.Add("project_id", b.projectId)
	req.Header.Add("Content-Type", "application/cbor")

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return
	}

	resb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		err = errors.New(string(resb))
	}

	return string(resb), err
}

// NewBlockfrostClient returns a wrapper for the blockfrost API/SDK with Node interface
func NewBlockfrostClient(projectId string, network *network.NetworkInfo) Node {
	var serverUrl string
	if network.NetworkId == 0 {
		serverUrl = blockfrost.CardanoTestNet
	} else {
		serverUrl = blockfrost.CardanoMainNet
	}

	client := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{
			ProjectID: projectId,
			Server:    serverUrl,
		},
	)

	return &blockfrostNode{
		network:   network,
		client:    client,
		projectId: projectId,
	}

}

// QueryTip is the equivalent of
// `cardano-cli query tip ${network_parameters}`
//
func (b *blockfrostNode) QueryTip() (nt NetworkTip, err error) {
	block, err := b.client.BlockLatest(context.TODO())

	if err != nil {
		return
	}

	nt = NetworkTip{
		Slot:  uint(block.Slot),
		Epoch: uint(block.Epoch),
		Block: uint(block.Height),
	}
	return
}
