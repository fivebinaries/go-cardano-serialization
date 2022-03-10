package address

type StakeCredentialType byte

const (
	KeyStakeCredentialType StakeCredentialType = iota
	ScriptStakeCredentialType
)

type StakeCredential struct {
	Kind    StakeCredentialType `cbor:"0,keyasint,omitempty"`
	Payload []byte              `cbor:"1,keyasint,omitempty"`
}

func NewKeyStakeCredential(hash []byte) *StakeCredential {
	return &StakeCredential{
		Kind:    KeyStakeCredentialType,
		Payload: hash,
	}
}

func NewScriptStakeCredential(hash []byte) *StakeCredential {
	return &StakeCredential{
		Kind:    ScriptStakeCredentialType,
		Payload: hash,
	}
}
