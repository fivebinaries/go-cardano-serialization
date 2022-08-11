package tx

import (
	"encoding/hex"

	"github.com/fxamacker/cbor/v2"
)

// TxBody contains the inputs, outputs, fee and titme to live for the transaction.
type TxBody struct {
	Inputs            []*TxInput  `cbor:"0,keyasint"`
	Outputs           []*TxOutput `cbor:"1,keyasint"`
	Fee               uint64      `cbor:"2,keyasint"`
	TTL               uint32      `cbor:"3,keyasint,omitempty"`
	AuxiliaryDataHash []byte      `cbor:"7,keyasint,omitempty"`
}

// NewTxBody returns a pointer to a new transaction body.
func NewTxBody() *TxBody {
	return &TxBody{
		Inputs:  make([]*TxInput, 0),
		Outputs: make([]*TxOutput, 0),
	}
}

// Bytes returns a slice of cbor Marshalled bytes.
func (b *TxBody) Bytes() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	return bytes, err
}

// Hex returns hex encoded string of the transaction bytes.
func (b *TxBody) Hex() (string, error) {
	by, err := b.Bytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(by), nil
}
