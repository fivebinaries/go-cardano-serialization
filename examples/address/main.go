package main

import (
	"encoding/hex"
	"log"

	"github.com/fivebinaries/go-cardano-serialization/bip32"
)

func harden(num uint) uint32 {
	return uint32(0x80000000 + num)
}

func main() {
	rootKey := bip32.FromBip39Entropy(
		//test walk nut penalty hip pave soap entry language right filter choice
		[]byte("ill swing joy endorse peanut erase width axis useless horse wheel super move seminar slice oil alcohol chronic image embody suggest ritual decline song"),
		[]byte{},
	)
	accountKey := rootKey.Derive(harden(1852)).Derive(harden(1815)).Derive(harden(0))

	//utxoPubKey := accountKey.Derive(0).Derive(0).Public()
	stakeKey := accountKey.Derive(2).Derive(0).Public()

	// utf8_reader, err := charmap.ISO8859_1.NewDecoder().Bytes(utxoPubKey)
	// if err != nil {
	// 	panic(err)
	// }
	log.Printf("%+v", hex.EncodeToString(stakeKey))
}
