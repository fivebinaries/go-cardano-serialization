package types

import "github.com/fivebinaries/go-cardano-serialization/utils"

// CheckedAdd implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#L256
func (v *Value) CheckedAdd(rhs *Value) (Value, error) {
	coin := v.GetCoin()
	coinR := rhs.GetCoin()
	multiasset := v.GetMultiasset()
	multiassetR := rhs.GetMultiasset()

	num := utils.BigNum(coin)
	resCoin, err := num.CheckedAdd(utils.BigNum(coinR))
	if err != nil {
		return Value{}, err
	}

	switch {
	case multiasset == nil && multiassetR == nil:
		return Value{V2SomeArray: &ValueAdditionalType0{
			V1: Coin(resCoin),
			V2: nil,
		}}, nil
	case multiasset == nil && multiassetR != nil:
		return Value{V2SomeArray: &ValueAdditionalType0{
			V1: Coin(resCoin),
			V2: multiasset,
		}}, nil
	case multiasset != nil && multiassetR == nil:
		return Value{V2SomeArray: &ValueAdditionalType0{
			V1: Coin(resCoin),
			V2: multiassetR,
		}}, nil
	default:
		panic("implement me")
	}
}

func (v *Value) GetCoin() Coin {
	if v.V1Coin != nil {
		return *v.V1Coin
	} else {
		return v.V2SomeArray.V1
	}
}

func (v *Value) GetMultiasset() MultiassetUint {
	if v.V1Coin != nil {
		return nil
	} else {
		return v.V2SomeArray.V2
	}
}

// PartialCmp implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#374
func (v *Value) PartialCmp(other *Value) int {
	compareAssets := func(lhs MultiassetUint, rhs MultiassetUint) int {
		switch {
		case lhs == nil && rhs == nil:
			return 0
		case lhs != nil && rhs == nil:
			panic("implement me")
		case lhs == nil && rhs != nil:
			panic("implement me")
		default:
			panic("implement me")
		}
	}

	order := compareAssets(v.V2SomeArray.V2, other.V2SomeArray.V2)
	coin := v.GetCoin()
	coinR := other.GetCoin()
	var coinOrder int
	switch {
	case coin == coinR:
		coinOrder = 0
	case coin > coinR:
		coinOrder = 1
	default:
		coinOrder = -1
	}

	switch {
	case order == 0:
		return coinOrder
	case coinOrder == 0 && order == -1:
		return -1
	case coinOrder == -1 && order == -1:
		return -1
	case coinOrder == 0 && order == 1:
		return 1
	case coinOrder == 1 && order == 1:
		return 1
	default:
		panic("undefined behavior")
	}
}
