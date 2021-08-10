package fees

import (
	"github.com/fivebinaries/go-cardano-serialization/common"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fxamacker/cbor/v2"
)

// LinearFee implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/fees.rs#L6
type LinearFee struct {
	Constant    types.Coin
	Coefficient types.Coin
}

// MinFee implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/fees.rs#L30
func MinFee(tx *types.Transaction, fee *LinearFee) (types.Coin, error) {
	txBody, err := cbor.Marshal(tx)
	if err != nil {
		return 0, err
	}
	txBodyLen := common.BigNum(len(txBody))
	if res, err := txBodyLen.CheckedMul(common.BigNum(uint64(fee.Coefficient))); err == nil {
		if res, err := res.CheckedAdd(common.BigNum(uint64(fee.Constant))); err == nil {
			return types.Coin(res), nil
		}
	}
	return 0, nil
}
