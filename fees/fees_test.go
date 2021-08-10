package fees

import (
	"encoding/hex"
	"github.com/fivebinaries/go-cardano-serialization/common"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"testing"
)

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/lib.rs#L2441
func TestTxSimpleUTXO(t *testing.T) {
	var inputs []types.TransactionInput
	transactionIdBytes, err := hex.DecodeString("3b40265111d8bb3c3c608d95b3a0bf83461ace32d79336579a1939b3aad1c0b7")
	if err != nil {
		t.Fatal(err)
	}
	transactionId, err := crypto.TransactionHashFromBytes(transactionIdBytes)
	inputs = append(inputs, types.TransactionInput{
		TransactionId: transactionId[:],
		Index:         0,
	})
	var outputs []types.TransactionOutput
	addrBytes, err := hex.DecodeString("611c616f1acb460668a9b2f123c80372c2adad3583b9c6cd2b1deeed1c")
	if err != nil {
		t.Fatal(err)
	}
	addr, err := types.AddressFromBytes(addrBytes)
	if err != nil {
		t.Fatal(err)
	}
	coin := types.Coin(1)
	outputs = append(outputs, types.TransactionOutput{
		V1: addr,
		Amount: types.Value{
			V1Coin:      &coin,
			V2SomeArray: nil,
		},
	})
	var ttl uint = 10
	body := types.TransactionBody{
		V1SetTransactionInput:   inputs,
		V2TransactionOutputList: outputs,
		V3Coin:                  94002,
		V4Uint:                  &ttl,
	}

	txBodyHash, err := common.HashTransaction(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(txBodyHash) != 1 {
	}

}
