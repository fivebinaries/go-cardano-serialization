package tx

type Witness struct {
	Keys []*VKeyWitness `cbor:"0,keyasint,omitempty"`
}

// NewTXWitness returns a pointer to a Witness created from VKeyWitnesses.
func NewTXWitness(keys ...*VKeyWitness) *Witness {
	if len(keys) == 0 {
		return &Witness{
			Keys: make([]*VKeyWitness, 0),
		}
	}

	return &Witness{
		Keys: keys,
	}
}

// VKeyWitness - Witness for use with Shelley based transactions
type VKeyWitness struct {
	_         struct{} `cbor:",toarray"`
	VKey      []byte
	Signature []byte
}

// NewVKeyWitness creates a Witness for Shelley Based transactions from a verification key and transaction signature.
func NewVKeyWitness(vkey, signature []byte) *VKeyWitness {
	return &VKeyWitness{
		VKey: vkey, Signature: signature,
	}
}

// BootstrapWitness for use with Byron/Legacy based transactions
type BootstrapWitness struct {
	_          struct{} `cbor:",toarray"`
	VKey       []byte
	Signature  []byte
	ChainCode  []byte
	Attributes []byte
}
