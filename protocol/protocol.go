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

	// The maximum block size (in bytes).
	// MaxBlockBodySize uint `json:"maxBlockBodySize"`

	// The maximum transaction size (in bytes).
	MaxTxSize uint `json:"maxTxSize"`

	// The maximum block header size (in bytes).
	// MaxBlockHeaderSize uint `json:"maxBlockHeaderSize"`

	// The amount (in Lovelaces) required for a deposit to register a StakeAddress.
	// StakeAddressDeposit uint `json:"stakeAddressDeposit"`

	// The amount (in Lovelaces) required for a deposit to register a stake pool.
	// StakePoolDeposit uint `json:"stakePoolDeposit"`

	// The Protocol Version
	ProtocolVersion ProtocolVersion `json:"protocolVersion"`

	// The decentralisation parameter (1 fully centralised, 0 fully decentralised).
	// Decentralisation float64 `json:"decentralisation"`

	// The monetary expansion rate.
	// MonetaryExpansion float64 `json:"monetaryExpansion"`

	// The influence of the pledge on a stake pool's probability on minting a block.
	// PoolPledgeInfluence float64 `json:"poolPledgeInfluence"`
}

// LoadProtocol unmarshalls protocol parameters from json
// generated from `cardano-cli`
//
//
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
