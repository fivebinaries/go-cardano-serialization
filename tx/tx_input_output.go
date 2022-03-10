package tx

import (
	"encoding/hex"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fxamacker/cbor/v2"
)

type TxInput struct {
	cbor.Marshaler

	TxHash []byte
	Index  uint16
	Amount uint
}

// NewTxInput creates and returns a *TxInput from Transaction Hash(Hex Encoded), Transaction Index and Amount.
func NewTxInput(txHash string, txIx uint16, amount uint) *TxInput {
	hash, _ := hex.DecodeString(txHash)

	return &TxInput{
		TxHash: hash,
		Index:  txIx,
		Amount: amount,
	}
}

func (txI *TxInput) MarshalCBOR() ([]byte, error) {
	type arrayInput struct {
		_      struct{} `cbor:",toarray"`
		TxHash []byte
		Index  uint16
	}
	input := arrayInput{
		TxHash: txI.TxHash,
		Index:  txI.Index,
	}
	return cbor.Marshal(input)
}

type TxOutput struct {
	_       struct{} `cbor:",toarray"`
	Address address.Address
	Amount  uint
}

func NewTxOutput(addr address.Address, amount uint) *TxOutput {
	return &TxOutput{
		Address: addr,
		Amount:  amount,
	}
}
