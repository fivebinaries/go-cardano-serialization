package types

import "github.com/fivebinaries/go-cardano-serialization/crypto"

// StakeCredential implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L67
type StakeCredential struct {
	Key    *crypto.Ed25519KeyHash
	Script *crypto.ScriptHash
}

// StakeCredentialFromKeyHash implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L79
func StakeCredentialFromKeyHash(hash []byte) *StakeCredential {
	var key crypto.Ed25519KeyHash
	copy(key[:], hash[:crypto.Ed25519KeyHashLen])
	return &StakeCredential{
		Key: &key,
	}
}

// StakeCredentialFromScriptHash implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L83
func StakeCredentialFromScriptHash(hash []byte) *StakeCredential {
	var script crypto.ScriptHash
	copy(script[:], hash[:crypto.ScriptHashLen])
	return &StakeCredential{
		Script: &script,
	}
}

// readAddrCred implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L339
func readAddrCred(data []byte, header byte, bit byte, pos int) *StakeCredential {
	hashBytes := data[pos : pos+crypto.Ed25519KeyHashLen]
	if header&(1<<bit) == 0 {
		return StakeCredentialFromKeyHash(hashBytes)
	}
	return StakeCredentialFromScriptHash(hashBytes)
}

// Kind implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L101
func (s *StakeCredential) Kind() byte {
	// don't use len(s.Key) != 0
	if s.Key != nil {
		return 0
	}
	return 1
}

// ToRawBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L108
func (s *StakeCredential) ToRawBytes() []byte {
	if s.Key != nil {
		return (*s.Key)[:]
	}
	return (*s.Script)[:]
}

// ToKeyHash implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L87
func (s *StakeCredential) ToKeyHash() *crypto.Ed25519KeyHash {
	if s.Key != nil {
		return s.Key
	}
	return nil
}

// ToScriptHash implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L87
func (s *StakeCredential) ToScriptHash() *crypto.ScriptHash {
	if s.Script != nil {
		return s.Script
	}
	return nil
}
