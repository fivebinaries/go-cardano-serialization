package fees

type LinearFee struct {
	TxFeePerByte uint
	TxFeeFixed   uint
}

func NewLinearFee(feePerByte uint, fixedFee uint) *LinearFee {
	return &LinearFee{
		TxFeePerByte: feePerByte,
		TxFeeFixed:   fixedFee,
	}
}
