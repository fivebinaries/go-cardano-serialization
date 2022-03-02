package node

import (
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
	"github.com/fivebinaries/go-cardano-serialization/tx"
)

type Node interface {
	// UTXOs returns list of unspent transaction outputs
	UTXOs() ([]tx.TxInput, error)

	// SubmitTx submits a cbor marshalled transaction to the cardano blockchain
	// using blockfrost or cardano-cli
	SubmitTx([]byte) error

	// ProtocolParameters returns Protocol Parameters from the network
	ProtocolParameters(*network.NetworkInfo) (protocol.Protocol, error)

	// QueryTip returns the tip of the network for use in tx building
	//
	// Using `query tip` on cardano-cli requires a synced local node
	QueryTip() (NetworkTip, error)
}
