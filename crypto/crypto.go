package crypto

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
