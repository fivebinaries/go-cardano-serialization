package bip32

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/crypto/edwards25519"
	"golang.org/x/crypto/pbkdf2"
)

const (
	XPrvSize = 96
)

type XPrv []byte
type XPub []byte
type PublicKey []byte
type PrivateKey []byte

func NewXPrv(seed []byte) (XPrv, error) {
	// Let k~ be 256-bit master secret.
	if len(seed) != 32 {
		return XPrv{}, errors.New("Seed needs to be 256 bits long")
	}

	// Then derive k = H_512(k~) and denote its left 32-byte by k_L and right one by k_R.
	extendedPrivateKey := sha512.Sum512(seed)

	// Otherwise additionally set the bits in k_L as follows:
	// The lowest 3 bits of the first byte of k_L are cleared.
	extendedPrivateKey[0] &= 0b_1111_1000
	// The highest bit of the last byte is cleared.
	// If the third highest bit of the last byte of k_L is not zero, discard k~.
	extendedPrivateKey[31] &= 0b_0101_1111
	// The second highest bit of the last byte is set
	extendedPrivateKey[31] |= 0b_0100_0000

	// The resulting pair (k_L, k_R) is the extended root private key.

	// And A <- [k_L]B is the root public key after encoding.

	// Derive c <- H_256(0x01 || k~) and call it the root chain code.
	rootChainCode := sha256.Sum256(append([]byte{0x01}, seed...))

	return append(extendedPrivateKey[:], rootChainCode[:]...), nil
}

func isHardened(index uint32) bool {
	return index >= 0x80000000
}

func (key XPrv) publicKey() [crypto.PublicKeyLen]byte {
	return MakePublicKey(key.extendedPrivateKey())
}

func (key XPrv) extendedPrivateKey() []byte {
	return key[:64]
}

func (key XPrv) chainCode() []byte {
	return key[64:]
}

func add28mul8(x, y []byte) []byte {
	var carry uint16

	out := make([]byte, 32, 32)
	for i := 0; i < 28; i++ {
		r := uint16(x[i]) + (uint16(y[i]) << 3) + carry
		out[i] = byte(r & 0xff)
		carry = r >> 8
	}
	for i := 28; i < 32; i++ {
		r := uint16(x[i]) + carry
		out[i] = byte(r & 0xff)
		carry = r >> 8
	}
	return out
}

func add256bits(x, y []byte) []byte {
	var carry uint16

	out := make([]byte, 32, 32)
	for i := 0; i < 32; i++ {
		r := uint16(x[i]) + uint16(y[i]) + carry
		out[i] = byte(r & 0xff)
		carry = r >> 8
	}
	return out
}

func (key XPrv) Derive(index uint32) XPrv {
	zmac := hmac.New(sha512.New, key.chainCode())
	imac := hmac.New(sha512.New, key.chainCode())

	serializedIndex := make([]byte, 4)
	binary.LittleEndian.PutUint32(serializedIndex, index)
	//hash.Write(serializedIndex)
	if isHardened(index) {
		// pk := []byte(ed25519.PrivateKey(key.extendedPrivateKey()).Public())
		zmac.Write([]byte{0x00})
		zmac.Write(key.extendedPrivateKey())
		zmac.Write(serializedIndex)
		imac.Write([]byte{0x01})
		imac.Write(key.extendedPrivateKey())
		imac.Write(serializedIndex)
	} else {
		pk := MakePublicKey(key.extendedPrivateKey())
		zmac.Write([]byte{0x02})
		zmac.Write(pk[:])
		zmac.Write(serializedIndex)
		imac.Write([]byte{0x03})
		imac.Write(pk[:])
		imac.Write(serializedIndex)
	}

	zout := zmac.Sum(nil)
	iout := imac.Sum(nil)

	left := add28mul8(key[:32], zout[:32])
	right := add256bits(key[32:64], zout[32:64])

	out := make([]byte, 0, 92)
	out = append(out, left...)
	out = append(out, right...)
	out = append(out, iout[32:]...)

	imac.Reset()
	zmac.Reset()

	return out
}

func (key XPrv) Public() XPub {
	out := make([]byte, 0, 64)
	pk := key.publicKey()
	out = append(out, pk[:]...)
	out = append(out, key.chainCode()...)
	return out
}

func (pub XPub) PublicKey() PublicKey {
	return PublicKey(pub[:32])
}

func (pub XPub) chainCode() PrivateKey {
	return PrivateKey(pub[32:])
}

//FromBip39Entropy implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/chain_crypto/derive.rs#L30
// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/crypto.rs#L123
func FromBip39Entropy(entropy []byte, password []byte) XPrv {
	const Iter = 4096
	pbkdf2_result := pbkdf2.Key(password, entropy, Iter, XPrvSize, sha512.New)
	return NormalizeBytesForce3rd(pbkdf2_result)
}

func NormalizeBytesForce3rd(bytes []byte) XPrv {
	bytes[0] &= 0b1111_1000
	bytes[31] &= 0b0001_1111
	bytes[31] |= 0b0100_0000
	return bytes
}

func (pub PublicKey) Hash() crypto.Ed25519KeyHash {
	return crypto.Blake2b224(pub)
}

func MakePublicKey(extendedSecret []byte) [crypto.PublicKeyLen]byte {
	var result [crypto.PublicKeyLen]byte
	var key [crypto.PublicKeyLen]byte
	copy(key[:], extendedSecret[:32])

	h := edwards25519.ExtendedGroupElement{}
	edwards25519.GeScalarMultBase(&h, &key)
	h.ToBytes(&result)
	return result
}
