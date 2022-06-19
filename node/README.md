# Node
[![GoDoc](https://godoc.org/github.com/fivebinaries/go-cardano-serialization/node?status.svg)](https://godoc.org/github.com/fivebinaries/go-cardano-serialization/node)


Package node provides an interface for cardano backends/nodes. Node implements two backends; blockfrost API and cardano-cli.

## Installation

```bash
go get github.com/fivebinaries/go-cardano-serialization/node
```

## Usage

Basic usage on how to setup a Blockfrost Node  on testnet and query for utxos for an address.

```golang
package main

import (
    "fmt"
    "log"
    "github.com/fivebinaries/go-cardano-serialization/node"
)

func prettyPrint(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "	")
	return string(bytes)
}

func main() {
    cli := node.NewBlockfrostClient(
		os.Getenv("BLOCKFROST_PROJECT_ID"),
		network.TestNet(),
	)

    addr, err := address.NewAddress(
        "addr_test1qqe6zztejhz5hq0xghlf72resflc4t2gmu9xjlf73x8dpf88d78zlt4rng3ccw8g5vvnkyrvt96mug06l5eskxh8rcjq2wyd63"
    )
    if err != nil {
        log.Fatal(err)
    }

    utxos, err := cli.UTXOs(addr)
    if err != nil {
        log.Fatal(err)
    }

    prettyPrint(utxos)
}

```

More examples on node usage can be found in the [`examples`](../examples/node/)

## License

Licensed under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0), see [`LICENSE`](https://github.com/fivebinaries/go-cardano-serialization/blob/master/LICENSE)