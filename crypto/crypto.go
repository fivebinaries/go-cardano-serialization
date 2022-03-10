package crypto

import (
	"errors"
	"log"

	"golang.org/x/crypto/blake2b"
)

const (
	Ed25519KeyHashLen      = 28
	ScriptHashLen          = 28
	TransactionHashLen     = 32
	GenesisDelegateHashLen = 28
	GenesisHashLen         = 28
	MetadataHashLen        = 32
	VRFKeyHashLen          = 32
	BlockHashLen           = 32
	VRFVKeyLen             = 32
	KESVKeyLen             = 32
	PublicKeyLen           = 32
	Blake2b224Len          = 28
	Blake2b256Len          = 32
)

type Ed25519KeyHash [Ed25519KeyHashLen]byte
type ScriptHash [ScriptHashLen]byte
type TransactionHash [TransactionHashLen]byte
type GenesisDelegateHash [GenesisDelegateHashLen]byte
type GenesisHash [GenesisHashLen]byte
type MetadataHash [MetadataHashLen]byte
type VRFKeyHash [VRFKeyHashLen]byte
type BlockHash [BlockHashLen]byte
type VRFVKey [VRFVKeyLen]byte
type KESVKey [KESVKeyLen]byte

// Blake2b224 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L15
func Blake2b224(data []byte) [Blake2b224Len]byte {
	b2b, err := blake2b.New(Blake2b224Len, nil)
	if err != nil {
		log.Fatalf("error blake2b224 transform: %s", err)
	}
	b2b.Write(data)
	var result [Blake2b224Len]byte
	copy(result[:], b2b.Sum(nil)[:Blake2b224Len])
	return result
}

// Blake2b256 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L15
func Blake2b256(data []byte) [Blake2b256Len]byte {
	b2b, err := blake2b.New(Blake2b256Len, nil)
	if err != nil {
		log.Fatalf("error Blake2b256 transform: %s", err)
	}
	b2b.Write(data)
	var result [Blake2b256Len]byte
	copy(result[:], b2b.Sum(nil)[:Blake2b256Len])
	return result
}

// Ed25519KeyHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func Ed25519KeyHashFromBytes(bytes []byte) (Ed25519KeyHash, error) {
	var res Ed25519KeyHash
	if len(bytes) != Ed25519KeyHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:Ed25519KeyHashLen])
	return res, nil
}

// ScriptHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func ScriptHashFromBytes(bytes []byte) (ScriptHash, error) {
	var res ScriptHash
	if len(bytes) != ScriptHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:ScriptHashLen])
	return res, nil
}

// TransactionHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func TransactionHashFromBytes(bytes []byte) (TransactionHash, error) {
	var res TransactionHash
	if len(bytes) != TransactionHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:TransactionHashLen])
	return res, nil
}

// GenesisDelegateHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func GenesisDelegateHashFromBytes(bytes []byte) (GenesisDelegateHash, error) {
	var res GenesisDelegateHash
	if len(bytes) != GenesisDelegateHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:GenesisDelegateHashLen])
	return res, nil
}

// GenesisHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func GenesisHashFromBytes(bytes []byte) (GenesisHash, error) {
	var res GenesisHash
	if len(bytes) != GenesisHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:GenesisHashLen])
	return res, nil
}

// MetadataHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func MetadataHashFromBytes(bytes []byte) (MetadataHash, error) {
	var res MetadataHash
	if len(bytes) != MetadataHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:MetadataHashLen])
	return res, nil
}

// VRFKeyHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func VRFKeyHashFromBytes(bytes []byte) (VRFKeyHash, error) {
	var res VRFKeyHash
	if len(bytes) != VRFKeyHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:VRFKeyHashLen])
	return res, nil
}

// BlockHashFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func BlockHashFromBytes(bytes []byte) (BlockHash, error) {
	var res BlockHash
	if len(bytes) != BlockHashLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:BlockHashLen])
	return res, nil
}

// VRFVKeyFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func VRFVKeyFromBytes(bytes []byte) (VRFVKey, error) {
	var res VRFVKey
	if len(bytes) != VRFVKeyLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:VRFVKeyLen])
	return res, nil
}

// KESVKeyFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L700
func KESVKeyFromBytes(bytes []byte) (KESVKey, error) {
	var res KESVKey
	if len(bytes) != KESVKeyLen {
		return res, errors.New("unexpected bytes")
	}
	copy(res[:], bytes[:KESVKeyLen])
	return res, nil
}
