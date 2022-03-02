package tx

import (
	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fivebinaries/go-cardano-serialization/protocol"
)

type TxBuilder struct {
	tx       *Tx
	xprvs    []bip32.XPrv
	protocol protocol.Protocol
}

// Sign adds a private key to create signature for witness
func (tb *TxBuilder) Sign(xprv bip32.XPrv) {
	tb.xprvs = append(tb.xprvs, xprv)
}

// Build creates hash of transaction, signs the hash using supplied witnesses and adds them to the transaction.
func (tb *TxBuilder) Build() (tx Tx, err error) {
	txKeys := []*VKeyWitness{}
	for _, prv := range tb.xprvs {
		hash, err := tb.tx.Hash()
		if err != nil {
			return tx, err
		}
		publicKey := prv.Public().PublicKey()
		signature := prv.Sign(hash[:])

		txKeys = append(txKeys, NewVKeyWitness(publicKey, signature[:]))
	}

	tb.tx.Witness = NewTXWitness(
		txKeys...,
	)

	return *tb.tx, nil
}

// Tx returns a pointer to the transaction
func (tb *TxBuilder) Tx() (tx *Tx) {
	return tb.tx
}

func (tb *TxBuilder) AddChangeIfNeeded(addr address.Address) {
	// change is amount in utxo minus outputs minus fee
	tb.tx.SetFee(1000000)
	fee := tb.MinFee()
	totalI := uint(0)
	totalO := uint(0)
	for _, inp := range tb.tx.Body.Inputs {
		totalI += inp.Amount
	}
	for _, out := range tb.tx.Body.Outputs {
		totalO += uint(out.Amount)
	}

	change := totalI - totalO - fee
	tb.tx.AddOutputs(
		NewTxOutput(
			addr,
			change,
		),
	)

	tb.tx.SetFee(tb.MinFee())
}

// SetTTL sets the
func (tb *TxBuilder) SetTTL(ttl uint32) {
	tb.tx.Body.TTL = ttl
}

// MinFee calculates the minimum fee for the provided transaction.
func (tb TxBuilder) MinFee() (fee uint) {
	vWitness := NewVKeyWitness(
		make([]byte, 64),
		make([]byte, 64),
	)
	innerTx := *tb.tx
	if len(innerTx.Witness.Keys) == 0 {
		innerTx.Witness.Keys = append(innerTx.Witness.Keys, vWitness)
	}
	lfee := fees.NewLinearFee(tb.protocol.TxFeePerByte, tb.protocol.TxFeeFixed)
	fee, _ = innerTx.Fee(lfee)

	return
}

// AddInputs adds inputs to the transaction body
func (tb *TxBuilder) AddInputs(inputs ...*TxInput) {
	tb.tx.AddInputs(inputs...)
}

// AddOutputs add outputs to the transaction body
func (tb *TxBuilder) AddOutputs(outputs ...*TxOutput) {
	tb.tx.AddOutputs(outputs...)
}

// NewTxBuilder returns pointer to a new TxBuilder.
func NewTxBuilder(pr protocol.Protocol, xprvs []bip32.XPrv) *TxBuilder {
	return &TxBuilder{
		tx:       NewTx(),
		xprvs:    xprvs,
		protocol: pr,
	}
}
