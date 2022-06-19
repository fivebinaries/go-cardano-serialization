# Examples
## Overview

Folder examples provides examples to get you started on creating addresses, building transactions and submitting transactions to Cardano network(testnet or mainnet).

## Installation
```bash
$ go get github.com/fivebinaries/go-cardano-serialization
```

## Examples
### Address
- [Generate Address](./address/generate/)

    Demonstrates generating addresses and exporting mnemonic backup phrases.

- [Restore Address](./address/restore/)


    Demonstrates generating from 24 word backup phrases and deriving enterprise or base addresses for payments. 

- [Simple Decode](./address/simple_decode/)

    Demonstrates usage of the NewAddress function provided in package address for decoding bech32/base58 encoded addresses into one of the supported types(enterprise, base, byron, pointer). The Address interface also provides a `String()` method to encode to bech32 or base58(byron).

### Node

Node package provides an interface for interacting with Blockfrost API or cardano-cli. The examples in `node/` demonstrate how to query UTXOs, Tip information, Protocol Parameters. The node interface can also be used to submit signed transactions using either of the backends.

### Transactions
- [Build Transaction](./transactions/)

    Demostrates building, signing and submitting transactions to the cardano network. Utlizes packages address and node.
