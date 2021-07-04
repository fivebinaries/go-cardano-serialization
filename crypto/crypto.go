package crypto

import (
	"golang.org/x/crypto/blake2b"
	"log"
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
