package utils

import (
	"errors"
	"math/bits"
)

// BigNum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#72
type BigNum uint64

func Quot(a int64, b int64) int64 {
	return (a - (a % b)) / b
}

func GetFilledArray(length int, val byte) []byte {
	var res []byte
	for len(res) < length {
		res = append(res, val)
	}
	return res
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

// Harden implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L749
// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/tx_builder.rs#L492
func Harden(index uint32) uint32 {
	return index | 0x80000000
}
