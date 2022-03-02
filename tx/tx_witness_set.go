package tx

type Witness struct {
	Keys []*VKeyWitness `cbor:"0,keyasint,omitempty"`
}

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

type VKeyWitness struct {
	_         struct{} `cbor:",toarray"`
	VKey      []byte
	Signature []byte
}

func NewVKeyWitness(vkey, signature []byte) *VKeyWitness {
	return &VKeyWitness{
		VKey: vkey, Signature: signature,
	}
}

type BootstrapWitness struct {
	_          struct{} `cbor:",toarray"`
	VKey       []byte
	Signature  []byte
	ChainCode  []byte
	Attributes []byte
}
