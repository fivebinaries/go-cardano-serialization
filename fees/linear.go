package fees

// LinearFee contains parameters for the linear fee equation `TxFeeFixed + len_bytes(tx) * TxFeePerByte`.
// These are provided in the protocol parameters
type LinearFee struct {
	TxFeePerByte uint
	TxFeeFixed   uint
}

// NewLinearFee returns a pointer to a new LinearFee from the provided params
func NewLinearFee(feePerByte uint, fixedFee uint) *LinearFee {
	return &LinearFee{
		TxFeePerByte: feePerByte,
		TxFeeFixed:   fixedFee,
	}
}
