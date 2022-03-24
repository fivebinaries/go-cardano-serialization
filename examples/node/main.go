// An example on how to use the `node` wrapper provided to query address UTXOs,
// network tip, protocol parameters and submit transactions to the network.
//
// The node interface currently has two backends;
//  - Blockfrost API
//  - cardano-cli
//
// Using Blockfrost API requires a `project_id` from blockfrost.io
// cardano-cli backend is not recommended when accepting raw user input

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fivebinaries/go-cardano-serialization/node"
)

func prettyPrint(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "	")
	return string(bytes)
}

func main() {
	// Init a new BlockfrostClient with provided project_id on the testnet network.
	cli := node.NewBlockfrostClient(
		os.Getenv("BLOCKFROST_PROJECT_ID"),
		network.TestNet(),
	)

	p, err := cli.ProtocolParameters()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nProtocol Parameters: %+v\n\n", prettyPrint(p))

	tip, err := cli.QueryTip()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tip: %+v\n\n", prettyPrint(tip))

	addr, err := address.NewAddress(
		"addr_test1qqe6zztejhz5hq0xghlf72resflc4t2gmu9xjlf73x8dpf88d78zlt4rng3ccw8g5vvnkyrvt96mug06l5eskxh8rcjq2wyd63",
	)
	if err != nil {
		log.Fatal(err)
	}

	utxos, err := cli.UTXOs(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UTXOS: %+v\n\n", prettyPrint(utxos))
}
