# Address

Package address implements structs for Cardano address types. 

## Overview 

Address handles serialization and deserilization of addresses used on the Cardano network. Currently supports:

- Byron/Legacy Address
- Enterprise Address
- Base Address
- Pointer Address

Address package also provides an `Address` interface and utility to load address from bech32/base58 encoded strings automatically into one of the supported address types.

## License
Package address is licensed under the Apache License