# Restore

This examples demonstrates restoring from 24 word backup keyphrase and creating enterprise and base payment addresses. Addresses can also be created for testnet and mainnet networks.


## Usage

```bash
$ go run main.go --help

Usage of main.exe
    -mnemonic string
        Mnemonic to restore wallet
    -network string
        The network ie mainnet or testnet (default "mainnet")
    -type uint
        Enum of address type(0: Enterprise Address, 1: Base Address)
```

## Example

To restore and create base address(payment key + stake key) for testnet from a recovery phrase run 
```bash
go run main.go -mnemonic "24_WORD_BACKUP_PHRASE" -type 1 -network testnet
// addr_test1qqe6zztejhz5hq0xghlf72resflc4t2gmu9xjlf73x8dpf88d78zlt4rng3ccw8g5vvnkyrvt96mug06l5eskxh8rcjq2wyd63
```