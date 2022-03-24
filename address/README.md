# Address

[![GoDoc](https://godoc.org/github.com/fivebinaries/go-cardano-serialization/address?status.svg)](https://godoc.org/github.com/fivebinaries/go-cardano-serialization/address)

Package address implements structs for Cardano address types. 

## Installation

```bash
go get -u github.com/fivebinaries/go-cardano-serialization/network
```

## Overview 

Address handles serialization and deserilization of addresses used on the Cardano network. Currently supports:

- Byron/Legacy Address
- Enterprise Address
- Base Address
- Pointer Address

Address package also provides an `Address` interface and utility to load address from bech32/base58 encoded strings automatically into one of the supported address types.

## Usage 

Examples on how to use the address package can be found in [`examples`](../examples/address/)

## License
Package address is licensed under the Apache License