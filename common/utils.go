package common

import (
	"errors"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fivebinaries/go-cardano-serialization/utils"
	"github.com/fxamacker/cbor/v2"
	"math/bits"
)

// BigNum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#72
type BigNum uint64

//todo remove
// Value implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#223
//type Value struct {
//	Coin Coin
//	MultiAsSet *hash_map.HashMap
//}

// HashTransaction implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#614
func HashTransaction(txBody types.TransactionBody) (crypto.TransactionHash, error) {
	txBodyBytes, err := cbor.Marshal(txBody)
	if err != nil {
		return [32]byte{}, err
	}
	b2bBytes := crypto.Blake2b256(txBodyBytes)
	return crypto.TransactionHashFromBytes(b2bBytes[:])
}

// Harden implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L749
// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L492
func Harden(index uint32) uint32 {
	return index | 0x80000000
}

// MinAdaRequired implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#750
func MinAdaRequired(assets *types.Value, minimumUTXOVal BigNum) BigNum {
	if assets.V2SomeArray == nil {
		return minimumUTXOVal
	}
	// NOTE: should be 2, but a bug in Haskell set this to 0
	coinSize := int64(0)
	txOutLenNoVal := int64(14)
	txInLen := int64(7)
	utxoEntrySizeWithoutVal := 6 + txOutLenNoVal + txInLen

	// NOTE: should be 29 but a bug in Haskell set this to 27
	adaOnlyUTXOSize := utxoEntrySizeWithoutVal + coinSize

	size := BundleSize(assets, &OutputSizeConstants{
		K0: 6,
		K1: 12,
		K2: 1,
	})

	v2 := BigNum(utils.Quot(int64(minimumUTXOVal), adaOnlyUTXOSize) * (utxoEntrySizeWithoutVal + int64(size)))
	if minimumUTXOVal < v2 {
		return v2
	} else {
		return minimumUTXOVal
	}
}

// OutputSizeConstants implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#699
type OutputSizeConstants struct {
	K0 uint
	K1 uint
	K2 uint
}

// BundleSize implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#710
func BundleSize(assets *types.Value, constants *OutputSizeConstants) uint {
	if assets.V2SomeArray == nil {
		return 1
	}

	//todo
	panic("implement me")
}

// MakeIcarusBootstrapWitness implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#581
func MakeIcarusBootstrapWitness(txBodyHash *crypto.TransactionHash, addr *types.ByronAddress, key *bip32.XPrv) (types.BootstrapWitness, error) {
	chainCode := key.ChainCode()
	vkey := key.Public()
	signature := key.Sign(txBodyHash[:])
	addr.Attributes.ProcessAttributes()
	attrRaw, err := cbor.Marshal(addr.Attributes.ProcessAttributes())
	if err != nil {
		return types.BootstrapWitness{}, err
	}

	return types.BootstrapWitness{
		PublicKey:  types.Vkey(vkey),
		Signature:  signature[:],
		ChainCode:  chainCode,
		Attributes: attrRaw,
	}, nil
}

// CheckedMul implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#159
func (b *BigNum) CheckedMul(other BigNum) (BigNum, error) {
	carryOut, res := bits.Mul64(uint64(*b), uint64(other))
	if carryOut != 0 {
		return 0, errors.New("overflow")
	}
	return BigNum(res), nil
}

// CheckedAdd implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#173
func (b *BigNum) CheckedAdd(other BigNum) (BigNum, error) {
	res, carryOut := bits.Add64(uint64(*b), uint64(other), 0)
	if carryOut != 0 {
		return 0, errors.New("overflow")
	}
	return BigNum(res), nil
}

// CheckedSub implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#159
func (b *BigNum) CheckedSub(other BigNum) (BigNum, error) {
	res, carryOut := bits.Sub64(uint64(*b), uint64(other), 0)
	if carryOut != 0 {
		return 0, errors.New("underflow")
	}
	return BigNum(res), nil
}
