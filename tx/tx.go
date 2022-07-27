package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fxamacker/cbor/v2"
	"golang.org/x/crypto/blake2b"
)

type Tx struct {
	_        struct{} `cbor:",toarray"`
	Body     *TxBody
	Witness  *Witness
	Valid    bool
	Metadata interface{}
}

// NewTx returns a pointer to a new Transaction
func NewTx() *Tx {
	return &Tx{
		Body:    NewTxBody(),
		Witness: NewTXWitness(),
		Valid:   true,
	}
}

// Bytes returns a slice of cbor marshalled bytes
func (t *Tx) Bytes() ([]byte, error) {
	if err := t.CalculateAuxiliaryDataHash(); err != nil {
		return nil, err
	}
	bytes, err := cbor.Marshal(t)
	return bytes, err
}

// Hex returns hex encoding of the transacion bytes
func (t *Tx) Hex() (string, error) {
	bytes, err := t.Bytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Hash performs a blake2b hash of the transaction body and returns a slice of [32]byte
func (t *Tx) Hash() ([32]byte, error) {
	if err := t.CalculateAuxiliaryDataHash(); err != nil {
		return [32]byte{}, err
	}
	txBody, err := cbor.Marshal(t.Body)
	if err != nil {
		var bt [32]byte
		return bt, err
	}

	txHash := blake2b.Sum256(txBody)
	return txHash, nil
}

// Fee returns the fee(in lovelaces) required by the transaction from the linear formula
// fee = txFeeFixed + txFeePerByte*tx_len_in_bytes
func (t *Tx) Fee(lfee *fees.LinearFee) (uint, error) {
	if err := t.CalculateAuxiliaryDataHash(); err != nil {
		return 0, err
	}
	txCbor, err := cbor.Marshal(t)
	if err != nil {
		return 0, err
	}
	txBodyLen := len(txCbor)
	fee := lfee.TxFeeFixed + lfee.TxFeePerByte*uint(txBodyLen)

	return fee, nil
}

// SetFee sets the fee
func (t *Tx) SetFee(fee uint) {
	t.Body.Fee = uint64(fee)
}

func (t *Tx) CalculateAuxiliaryDataHash() error {
	if t.Metadata != nil {
		mdBytes, err := cbor.Marshal(&t.Metadata)
		if err != nil {
			return fmt.Errorf("cannot serialize metadata: %w", err)
		}
		auxHash := blake2b.Sum256(mdBytes)
		t.Body.AuxiliaryDataHash = auxHash[:]
	}
	return nil
}

// AddInputs adds the inputs to the transaction body
func (t *Tx) AddInputs(inputs ...*TxInput) error {
	t.Body.Inputs = append(t.Body.Inputs, inputs...)

	return nil
}

// AddOutputs adds the outputs to the transaction body
func (t *Tx) AddOutputs(outputs ...*TxOutput) error {
	t.Body.Outputs = append(t.Body.Outputs, outputs...)

	return nil
}
