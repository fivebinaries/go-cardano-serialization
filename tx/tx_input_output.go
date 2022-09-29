package tx

import (
	"encoding/hex"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fxamacker/cbor/v2"
)

type TxInput struct {
	TxHash []byte
	Index  uint16
	Amount uint
}

var (
	_ cbor.Marshaler   = &TxInput{}
	_ cbor.Unmarshaler = &TxInput{}
)

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

func (txI *TxInput) UnmarshalCBOR(in []byte) error {
	type arrayInput struct {
		_      struct{} `cbor:",toarray"`
		TxHash []byte
		Index  uint16
	}
	input := arrayInput{}
	if err := cbor.Unmarshal(in, &input); err != nil {
		return err
	}
	txI.TxHash = input.TxHash
	txI.Index = input.Index
	return nil
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

func (txO *TxOutput) UnmarshalCBOR(in []byte) (err error) {
	type arrayOutput struct {
		_       struct{} `cbor:",toarray"`
		Address []byte
		Amount  uint
	}
	output := arrayOutput{}
	if err = cbor.Unmarshal(in, &output); err != nil {
		return err
	}
	txO.Address, err = address.NewAddressFromBytes(output.Address)
	txO.Amount = output.Amount
	return err
}
