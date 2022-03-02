package tx

import (
	"encoding/hex"

	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fxamacker/cbor/v2"
	"golang.org/x/crypto/blake2b"
)

type Tx struct {
	_        struct{} `cbor:",toarray"`
	Body     *TxBody
	Witness  *Witness
	Metadata interface{}
}

func NewTx() *Tx {
	return &Tx{
		Body:    NewTxBody(),
		Witness: NewTXWitness(),
	}
}

func (t *Tx) Bytes() ([]byte, error) {
	bytes, err := cbor.Marshal(t)
	return bytes, err
}

func (t *Tx) Hex() (string, error) {
	bytes, err := t.Bytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (t *Tx) Hash() ([32]byte, error) {
	txBody, err := cbor.Marshal(t.Body)
	if err != nil {
		var bt [32]byte
		return bt, err
	}

	txHash := blake2b.Sum256(txBody)
	return txHash, nil
}

func (t *Tx) Fee(lfee *fees.LinearFee) (uint, error) {
	txCbor, err := cbor.Marshal(t)
	if err != nil {
		return 0, err
	}
	txBodyLen := len(txCbor)
	fee := lfee.TxFeeFixed + lfee.TxFeePerByte*uint(txBodyLen)

	return fee, nil
}

func (t *Tx) SetFee(fee uint) {
	t.Body.Fee = uint64(fee)
}

func (t *Tx) AddInputs(inputs ...*TxInput) error {
	t.Body.Inputs = append(t.Body.Inputs, inputs...)

	return nil
}

func (t *Tx) AddOutputs(outputs ...*TxOutput) error {
	t.Body.Outputs = append(t.Body.Outputs, outputs...)

	return nil
}
