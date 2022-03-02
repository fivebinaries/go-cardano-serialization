package tx

import (
	"encoding/hex"

	"github.com/fxamacker/cbor/v2"
)

type TxBody struct {
	Inputs  []*TxInput  `cbor:"0,keyasint"`
	Outputs []*TxOutput `cbor:"1,keyasint"`
	Fee     uint64      `cbor:"2,keyasint"`
	TTL     uint32      `cbor:"3,keyasint,omitempty"`
}

func NewTxBody() *TxBody {
	return &TxBody{
		Inputs:  make([]*TxInput, 0),
		Outputs: make([]*TxOutput, 0),
	}
}

func (b *TxBody) Bytes() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	return bytes, err
}

func (b *TxBody) Hex() (string, error) {
	by, err := b.Bytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(by), nil
}
