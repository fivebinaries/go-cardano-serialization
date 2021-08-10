package transactions

import (
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/common"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fivebinaries/go-cardano-serialization/lib"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"testing"
)

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L482
func genesisId() crypto.TransactionHash {
	var bytes []byte
	for len(bytes) < crypto.TransactionHashLen {
		bytes = append(bytes, 0)
	}
	res, err := crypto.TransactionHashFromBytes(bytes)
	if err != nil {
		panic("undefined error")
	}
	return res
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L482
func rootKey15() bip32.XPrv {
	// art forum devote street sure rather head chuckle guard poverty release quote oak craft enemy
	entropy := []byte{0x0c, 0xcb, 0x74, 0xf3, 0x6b, 0x7d, 0xa1, 0x64, 0x9a, 0x81, 0x44, 0x67, 0x55, 0x22, 0xd4, 0xd8, 0x09, 0x7c, 0x64, 0x12}
	return bip32.FromBip39Entropy(entropy, []byte{})
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L497
func TestBuildTxWithChange(t *testing.T) {
	linearFee := fees.LinearFee{
		Constant:    2,
		Coefficient: 500,
	}
	txBuilder := NewTransactionBuilder(&linearFee, 1, 1, 1)
	spend := rootKey15().
		Derive(common.Harden(1852)).
		Derive(common.Harden(1815)).
		Derive(common.Harden(0)).
		Derive(0).
		Derive(0).
		Public()
	//changeKey := rootKey15().
	//	Derive(common.Harden(1852)).
	//	Derive(common.Harden(1815)).
	//	Derive(common.Harden(0)).
	//	Derive(1).
	//	Derive(0).
	//	Public()
	stake := rootKey15().
		Derive(common.Harden(1852)).
		Derive(common.Harden(1815)).
		Derive(common.Harden(0)).
		Derive(2).
		Derive(0).
		Public()
	spendHash := spend.PublicKey().Hash()
	stakeHash := stake.PublicKey().Hash()
	spendCred := types.StakeCredentialFromKeyHash(spendHash[:])
	stakeCred := types.StakeCredentialFromKeyHash(stakeHash[:])
	addrNet0 := types.NewBaseAddress(types.TestNet().NetworkId, spendCred, stakeCred)
	coin := types.Coin(1000000)
	coinOut := types.Coin(10)
	genId := genesisId()
	txBuilder.AddKeyInput(spendHash, &types.TransactionInput{
		TransactionId: genId[:],
		Index:         0,
	}, &types.Value{
		V1Coin:      &coin,
		V2SomeArray: nil,
	})
	err := txBuilder.AddOutput(&types.TransactionOutput{
		V1: &addrNet0,
		Amount: types.Value{
			V1Coin:      &coinOut,
			V2SomeArray: nil,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ttl := lib.Slot(1000)
	txBuilder.TTL = &ttl

	//changeCred := types.StakeCredentialFromKeyHash(changeKey)
	//changeAddr := types.NewBaseAddress(types.TestNet().NetworkId, changeCred, stakeCred)
	//addedChage := txBuilder.add

}
