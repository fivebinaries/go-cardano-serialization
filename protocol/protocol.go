package protocol

import (
	"encoding/json"
	"io/ioutil"
)

type ProtocolVersion struct {
	Major uint8 `json:"major"`
	Minor uint8 `json:"minor"`
}

// Protocol contains protocol parameters.
type Protocol struct {
	// The 'a' parameter to calculate the minimum transaction fee in the
	// linear equation a * byte_size(tx) + b
	TxFeePerByte uint `json:"txFeePerByte"`

	// The 'b' parameter to calculate the minimum transaction fee.
	TxFeeFixed uint `json:"txFeeFixed"`

	// The maximum transaction size (in bytes).
	MaxTxSize uint `json:"maxTxSize"`

	// The Protocol Version
	ProtocolVersion ProtocolVersion `json:"protocolVersion"`

	// Minimum UTXO Value
	MinUTXOValue uint `json:"minUTxOValue"`
}

// LOadProtocol returns a pointer to a unmarshalled Protocol given a file path of a
// protocol parameters file in the cardano-cli generated format.
func LoadProtocol(fp string) (*Protocol, error) {
	pByte, err := ioutil.ReadFile(fp)
	if err != nil {
		return &Protocol{}, err
	}
	return loadProtocolFromBytes(pByte)
}

func loadProtocolFromBytes(pByte []byte) (*Protocol, error) {
	p := &Protocol{}

	err := json.Unmarshal(pByte, p)
	return p, err
}
