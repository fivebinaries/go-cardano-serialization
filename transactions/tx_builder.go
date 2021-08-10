package transactions

import (
	"errors"
	"fmt"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/common"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/fees"
	"github.com/fivebinaries/go-cardano-serialization/lib"
	"github.com/fivebinaries/go-cardano-serialization/metadata"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fivebinaries/go-cardano-serialization/utils"
)

// TxBuilderInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L109
type TxBuilderInput struct {
	Input  types.TransactionInput
	Amount types.Value
}

// MockWitnessSet implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L102
type MockWitnessSet struct {
	//todo change to set type
	VKeys      []crypto.Ed25519KeyHash
	Scripts    []crypto.ScriptHash
	Bootstraps [][]uint8
}

// TransactionBuilder implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L116
type TransactionBuilder struct {
	MinimumUTXOval        common.BigNum
	PoolDeposit           common.BigNum
	KeyDeposit            common.BigNum
	FeeAlgo               fees.LinearFee
	Inputs                []TxBuilderInput
	Outputs               []types.TransactionOutput
	Fee                   *types.Coin
	TTL                   *lib.Slot
	Certs                 []types.Certificate
	Withdrawals           *types.Withdrawals
	Metadata              *metadata.TransactionMetadata
	ValidityStartInterval *lib.Slot
	InputTypes            MockWitnessSet
	Mint                  *types.Mint
}

// NewTransactionBuilder implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L284
func NewTransactionBuilder(linearFee *fees.LinearFee, minimumUTXOVal types.Coin, poolDeposit common.BigNum, keyDeposit common.BigNum) TransactionBuilder {
	return TransactionBuilder{
		MinimumUTXOval:        common.BigNum(minimumUTXOVal),
		PoolDeposit:           poolDeposit,
		KeyDeposit:            keyDeposit,
		FeeAlgo:               *linearFee,
		Inputs:                nil,
		Outputs:               nil,
		Fee:                   nil,
		TTL:                   nil,
		Certs:                 nil,
		Withdrawals:           nil,
		Metadata:              nil,
		ValidityStartInterval: nil,
		InputTypes:            MockWitnessSet{},
		Mint:                  nil,
	}
}

// AddKeyInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L138
func (t *TransactionBuilder) AddKeyInput(hash crypto.Ed25519KeyHash, input *types.TransactionInput, amount *types.Value) {
	t.Inputs = append(t.Inputs, TxBuilderInput{
		Input:  *input,
		Amount: *amount,
	})
	t.InputTypes.VKeys = append(t.InputTypes.VKeys, hash)
}

// AddScriptInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L145
func (t *TransactionBuilder) AddScriptInput(hash crypto.ScriptHash, input *types.TransactionInput, amount *types.Value) {
	t.Inputs = append(t.Inputs, TxBuilderInput{
		Input:  *input,
		Amount: *amount,
	})
	t.InputTypes.Scripts = append(t.InputTypes.Scripts, hash)
}

// AddBootstrapInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L152
func (t *TransactionBuilder) AddBootstrapInput(hash *types.ByronAddress, input *types.TransactionInput, amount *types.Value) {
	t.Inputs = append(t.Inputs, TxBuilderInput{
		Input:  *input,
		Amount: *amount,
	})
	t.InputTypes.Bootstraps = append(t.InputTypes.Bootstraps, hash.ToBytes())
}

// AddInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L160
func (t *TransactionBuilder) AddInput(inputAddress *types.Address, input *types.TransactionInput, amount *types.Value) {
	switch addr := (*inputAddress).(type) {
	case *types.BaseAddress:
		if hash := addr.Payment.ToKeyHash(); hash != nil {
			t.AddKeyInput(*hash, input, amount)
			return
		}
		if hash := addr.Payment.ToScriptHash(); hash != nil {
			t.AddScriptInput(*hash, input, amount)
		}
	case *types.EnterpriseAddress:
		if hash := addr.Payment.ToKeyHash(); hash != nil {
			t.AddKeyInput(*hash, input, amount)
			return
		}
		if hash := addr.Payment.ToScriptHash(); hash != nil {
			t.AddScriptInput(*hash, input, amount)
		}
	case *types.PointerAddress:
		if hash := addr.Payment.ToKeyHash(); hash != nil {
			t.AddKeyInput(*hash, input, amount)
			return
		}
		if hash := addr.Payment.ToScriptHash(); hash != nil {
			t.AddScriptInput(*hash, input, amount)
		}
	case *types.ByronAddress:
		t.AddBootstrapInput(addr, input, amount)
	}
}

// AddOutput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L209
func (t *TransactionBuilder) AddOutput(output *types.TransactionOutput) error {
	minAda := common.MinAdaRequired(&output.Amount, t.MinimumUTXOval)
	var coin types.Coin
	if output.Amount.V1Coin != nil {
		coin = *output.Amount.V1Coin
	} else {
		coin = output.Amount.V2SomeArray.V1
	}
	if common.BigNum(coin) < minAda {
		return fmt.Errorf("value %v less than the minimum UTXO value %v", coin, minAda)
	} else {
		t.Outputs = append(t.Outputs, *output)
	}
	return nil
}

// GetExplicitInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L314
func (t *TransactionBuilder) GetExplicitInput() (types.Value, error) {
	res := types.Value{}
	res.V2SomeArray = &types.ValueAdditionalType0{
		V1: 0,
		V2: nil,
	}
	var err error
	for _, input := range t.Inputs {
		res, err = res.CheckedAdd(&input.Amount)
		if err != nil {
			return types.Value{}, err
		}
	}
	return res, nil
}

// GetImplicitInput implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L322
func (t *TransactionBuilder) GetImplicitInput() (types.Value, error) {
	var withdrawalsSum common.BigNum
	var certificateRefund common.BigNum
	var err error
	for _, w := range *t.Withdrawals {
		withdrawalsSum, err = withdrawalsSum.CheckedAdd(common.BigNum(w.Value.(types.Coin)))
		if err != nil {
			return types.Value{}, err
		}
	}
	for _, c := range t.Certs {
		switch {
		case c.V5 != nil: //pool retirement
			certificateRefund, err = certificateRefund.CheckedAdd(t.PoolDeposit)
		case c.V2 != nil: //stake deregistration
			certificateRefund, err = certificateRefund.CheckedAdd(t.KeyDeposit)
		}
		if err != nil {
			return types.Value{}, err
		}
	}

	resB, err := withdrawalsSum.CheckedAdd(certificateRefund)
	res := types.Value{}
	res.V2SomeArray = &types.ValueAdditionalType0{
		V1: types.Coin(resB),
		V2: nil,
	}
	return res, err
}

// GetDeposit implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#L650
func (t *TransactionBuilder) GetDeposit() (types.Value, error) {
	panic("implement me")
}

// AddChangeIfNeeded implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L354
/// Warning: this function will mutate the /fee/ field
func (t *TransactionBuilder) AddChangeIfNeeded(addressInput types.Address) (bool, error) {
	panic("implement me")
	//var fee types.Coin
	//var err error
	//if t.Fee == nil {
	//	fee, err = t.MinFee()
	//	if err != nil {
	//		return false, err
	//	}
	//} else {
	//	return false, errors.New("cannot calculate change if fee was explicitly specified")
	//}
	//
	//explicitInput, err := t.GetExplicitInput()
	//if err != nil {
	//	return false, err
	//}
	//implicitInput, err := t.GetImplicitInput()
	//if err != nil {
	//	return false, err
	//}
	//inputTotal, err := explicitInput.CheckedAdd(&implicitInput)
	//if err != nil {
	//	return false, err
	//}
	//
	//outputTotal, err := explicitInput.CheckedAdd(types.Value{
	//	V1Coin:      t.de,
	//	V2SomeArray: nil,
	//})
}

// MinFee implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L470
/// warning: sum of all parts of a transaction must equal 0. You cannot just set the fee to the min value and forget about it
/// warning: min_fee may be slightly larger than the actual minimum fee (ex: a few lovelaces)
/// this is done to simplify the library code, but can be fixed later
func (t *TransactionBuilder) MinFee() (types.Coin, error) {
	selfCopy := *t
	coin := types.Coin(0x100000000)
	selfCopy.Fee = &coin
	return MinFee(&selfCopy)
}

// MinFee implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L40
func (t *TransactionBuilder) Build() (types.TransactionBody, error) {
	if t.Fee == nil {
		return types.TransactionBody{}, errors.New("fee not specified")
	}

	var inputs []types.TransactionInput
	for _, inp := range t.Inputs {
		inputs = append(inputs, inp.Input)
	}

	ttl := uint(*t.TTL)
	validityStartInterval := uint(*t.ValidityStartInterval)
	var metadat *types.MetadataHash
	if t.Metadata != nil {
		hash, err := metadata.HashMetadata(t.Metadata)
		if err != nil {
			return types.TransactionBody{}, err
		}
		metadat = &hash
	}

	return types.TransactionBody{
		V1SetTransactionInput:   inputs,
		V2TransactionOutputList: t.Outputs,
		V3Coin:                  *t.Fee,
		V4Uint:                  &ttl,
		V5CertificateList:       &t.Certs,
		V6Withdrawals:           t.Withdrawals,
		V7Update:                nil,
		V8MetadataHash:          metadat,
		V9Uint:                  &validityStartInterval,
		V10Mint:                 t.Mint,
	}, nil
}

// MinFee implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L40
func MinFee(txBuilder *TransactionBuilder) (types.Coin, error) {
	body, err := txBuilder.Build()
	if err != nil {
		return 0, err
	}

	fakeKeyRoot := bip32.FromBip39Entropy(
		// art forum devote street sure rather head chuckle guard poverty release quote oak craft enemy
		[]byte{0x0c, 0xcb, 0x74, 0xf3, 0x6b, 0x7d, 0xa1, 0x64, 0x9a, 0x81, 0x44, 0x67, 0x55, 0x22, 0xd4, 0xd8, 0x09, 0x7c, 0x64, 0x12},
		[]byte{})

	var vkeys []types.Vkeywitness

	if len(txBuilder.InputTypes.VKeys) != 0 {
		for _ = range txBuilder.InputTypes.VKeys {
			sign := fakeKeyRoot.Sign(utils.GetFilledArray(100, 1))
			vkeys = append(vkeys, types.Vkeywitness{
				V1: types.Vkey(fakeKeyRoot.Public()),
				V2: sign[:],
			})
		}
	}

	if len(txBuilder.InputTypes.Scripts) != 0 {
		return 0, errors.New("scripts inputs not supported yet")
	}

	var bootstrapKeys []types.BootstrapWitness
	if len(txBuilder.InputTypes.Bootstraps) != 0 {
		for _, addr := range txBuilder.InputTypes.Bootstraps {
			// picking icarus over daedalus for fake witness generation shouldn't matter
			hash, err := common.HashTransaction(body)
			if err != nil {
				return 0, err
			}
			addrTmp, err := types.FromBytes(addr)
			if err != nil {
				return 0, err
			}
			byronAddr := addrTmp.ToByronAddress()
			elem, err := common.MakeIcarusBootstrapWitness(&hash, &byronAddr, &fakeKeyRoot)
			if err != nil {
				return 0, err
			}
			bootstrapKeys = append(bootstrapKeys, elem)
		}
	}

	witnessSet := types.TransactionWitnessSet{
		V1VkeywitnessList:      &vkeys,
		V2NativeScriptList:     nil,
		V3BootstrapWitnessList: &bootstrapKeys,
	}
	auxData := types.AuxiliaryData{
		V1TransactionMetadatumLabelMap: nil,
		V2SomeArray: &types.AuxiliaryDataAdditionalType1{
			TransactionMetadata: *txBuilder.Metadata.General.ToHashMap(),
			AuxiliaryScripts:    txBuilder.Metadata.Native,
		},
	}
	fullTx := types.Transaction{
		V1: body,
		V2: witnessSet,
		V3: &auxData,
	}
	return fees.MinFee(&fullTx, &txBuilder.FeeAlgo)
}
