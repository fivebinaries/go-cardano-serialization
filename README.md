<p align="center">
  <a href="https://fivebinaries.com" target="_blank" align="center">
    <img src="https://raw.githubusercontent.com/fivebinaries/go-cardano-serialization/dev/.github/go-cardano-serialization-logo.svg" width="100">
  </a>
  <br />
</p>

# Go Cardano Serialization Library
[![GoDoc](https://godoc.org/github.com/fivebinaries/go-cardano-serialization?status.svg)](https://godoc.org/github.com/fivebinaries/go-cardano-serialization)

Golang library for serialization and deserialiation of Cardano data structures. 

## Installation

```bash
$ go get https://github.com/fivebinaries/go-cardano-serialization
```

## Usage
### Creating a simple transaction

The simplest transaction on Cardano network contains inputs(Unspent Transaction Outputs) and output.

```golang
package main

import (
    "log"

    "github.com/fivebinaries/go-cardano-serialization/address"
    "github.com/fivebinaries/go-cardano-serialization/tx"
)

func main() {
    adaTx := tx.NewTx()
    adaTx.AddInput(
        tx.NewInput(
            "TX_HASH", // Transaction Hash
            0,         // Transaction Index
            10000000   // Lovelace value of UTXO
        )
    )

    receiverAddr, err := address.NewAddress("addr1bech32_receiver_address_here")
    if err != nil {
        log.Fatal(err)
    }

    adaTx.AddOutput(
        tx.NewOutput(
            receiverAddr,
            5000000
        )
    )

    // Set an estimated transaction cost
    adaTx.SetFee(170000)

    // Set the transaction's time to live
    adaTx.SetTTL(505050505)

    // Encode example transaction to cbor hex.
    fmt.Println(adaTx.Hex())
}
```

More examples covering building through signing and submission of transactions can be found in the [`examples`](./examples/) folder.

## License

Licensed under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0), see [`LICENSE`](https://github.com/fivebinaries/go-cardano-serialization/blob/master/LICENSE)