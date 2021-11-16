package types

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/utils"
)

type TestTask struct {
	AddressBase58 string
	PublicKey     []byte
}

var tests = []TestTask{
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L414
	{"DdzFFzCqrhsrcTVhLygT24QwTnNqQqQ8mZrq5jykUzMveU26sxaH529kMpo7VhPrt5pwW3dXeB2k3EEvKcNBRmzCfcQ7dTkyGzTs658C",
		[]byte{
			0x6a, 0x50, 0x96, 0x89, 0xc6, 0x53, 0x17, 0x58, 0x65, 0x98, 0x5a, 0xd1, 0xe0, 0xeb,
			0x5f, 0xf9, 0xad, 0xa6, 0x99, 0x7a, 0xa4, 0x03, 0xe6, 0x48, 0x61, 0x4b, 0x3b, 0x78,
			0xfc, 0xba, 0x9c, 0x27, 0x30, 0x82, 0x28, 0xd9, 0x87, 0x2a, 0xf8, 0xb6, 0x5b, 0x98,
			0x7f, 0xf2, 0x3e, 0x1a, 0x20, 0xcd, 0x90, 0xd8, 0x34, 0x6c, 0x31, 0xf0, 0xed, 0xb8,
			0x99, 0x89, 0x52, 0xdc, 0x67, 0x66, 0x55, 0x80,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L427
	{"DdzFFzCqrht4it4GYgBp4J39FNnKBsPFejSppARXHCf2gGiTJcwXzpRvgDmxPvKQ8aZZmVqcLUz5L66a8Ja46pfKVtFRaKyn9eKdvpaC",
		[]byte{
			0xff, 0x7b, 0xf1, 0x29, 0x9d, 0xf3, 0xd7, 0x17, 0x98, 0xae, 0xfd, 0xc4, 0xae, 0xa7,
			0xdb, 0x2f, 0x8d, 0xb7, 0x60, 0x46, 0x56, 0x94, 0x41, 0xea, 0xe5, 0x8b, 0x72, 0x23,
			0xb6, 0x8b, 0x44, 0x04, 0x82, 0x15, 0xcb, 0xac, 0x94, 0xbc, 0xb7, 0xf2, 0xcf, 0x33,
			0x6c, 0x6c, 0x18, 0xbc, 0x3e, 0x71, 0x3f, 0xfd, 0x82, 0x67, 0x59, 0x4f, 0xf6, 0x34,
			0x93, 0x32, 0xce, 0x4f, 0x98, 0x04, 0xa7, 0xff,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L440
	{"DdzFFzCqrhsvNQtyViTvEdGxfdc5T1E5RorzFWjYodqjhFDy8fQxfDPccmTc4ePbvkiwvRkR8dtqQ1SHpH53fDSoxD17fo9f6WkRjjAA",
		[]byte{
			0x5c, 0x36, 0x51, 0xe0, 0xeb, 0x9d, 0x6d, 0xc9, 0x64, 0x07, 0x13, 0x7c, 0xcc, 0x1f,
			0x37, 0x7a, 0x87, 0x94, 0x61, 0x77, 0xa5, 0x2c, 0xa3, 0x77, 0x2c, 0x6b, 0x4b, 0xeb,
			0x72, 0x39, 0x50, 0xdc, 0x50, 0x22, 0x46, 0x68, 0x21, 0x8b, 0x8b, 0x36, 0x62, 0x02,
			0xfe, 0x5b, 0x7d, 0x55, 0x6f, 0x50, 0x1c, 0x5c, 0x4e, 0x2d, 0x58, 0xe0, 0x54, 0x67,
			0xe1, 0xab, 0xc0, 0x44, 0xc6, 0xc1, 0xbf, 0x8e,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L454
	{"DdzFFzCqrhsn7ZAhKy8mxkzW6G3wryM7K6bH38VAjE2FesJMxia3UviivMvGz146TP1FpDharxTE6nUgCCnZx2fmtKpmxAosg9Tf5b8y",
		[]byte{
			0xcd, 0x84, 0x2e, 0x01, 0x0d, 0x81, 0xa6, 0xbe, 0x1e, 0x16, 0x9f, 0xd6, 0x35, 0x21,
			0xdb, 0xb9, 0x5f, 0x42, 0x41, 0xfc, 0x82, 0x3f, 0x45, 0xb1, 0xcf, 0x1a, 0x1c, 0xb4,
			0xc5, 0x89, 0x57, 0x27, 0x1d, 0x4d, 0x14, 0x2a, 0x22, 0x94, 0xea, 0x5f, 0xa3, 0x16,
			0xa4, 0xad, 0xbf, 0xcd, 0x59, 0x7a, 0x7c, 0x89, 0x6a, 0x52, 0xa9, 0xa3, 0xa9, 0xce,
			0x49, 0x64, 0x4a, 0x10, 0x2d, 0x00, 0x71, 0x99,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L466
	{"DdzFFzCqrhssTCJf4sv664bdQURovAwzx1hNKkMkNLwMNyaxZFuPSDdZTTRMcoDyXHuCiZhbD4umvMJcWGkvFMMzBoBUW5UBdBbDqXGX",
		[]byte{
			0x5a, 0xac, 0x2d, 0xd0, 0xa8, 0xdc, 0x5d, 0x61, 0x0a, 0x4b, 0x6f, 0xdf, 0x3f, 0x5e,
			0xf1, 0xb6, 0x4a, 0xcb, 0x76, 0xb1, 0xe8, 0x1f, 0x6a, 0x35, 0x70, 0x31, 0xfa, 0x19,
			0xd5, 0xe6, 0x56, 0x9d, 0xcc, 0x37, 0xb7, 0xae, 0x6f, 0x39, 0x15, 0x82, 0xfb, 0x05,
			0x4b, 0x72, 0xba, 0xda, 0x90, 0xab, 0x14, 0x6c, 0xdd, 0x01, 0x42, 0x0e, 0x4b, 0x40,
			0x18, 0xf1, 0xa0, 0x55, 0x29, 0x82, 0xd2, 0x31,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L479
	{"DdzFFzCqrhsfi5fFjJUHYPSnfTYrnMohzh3PrrtrVQgwua33HWPKUdTJXo3o77pSGCmDNrjYaAiZmJddaPW9iHyUDatvU2WhX7MgnNMy",
		[]byte{
			0x2a, 0x6a, 0xd1, 0x51, 0x09, 0x96, 0xff, 0x2d, 0x10, 0x89, 0xcb, 0x8e, 0xd5, 0xf5,
			0xc0, 0x61, 0xf6, 0xad, 0x0a, 0xfb, 0xb5, 0x3d, 0x95, 0x40, 0xa0, 0xfc, 0x89, 0xef,
			0xc0, 0xa2, 0x63, 0xb9, 0x6d, 0xac, 0x00, 0xbd, 0x0d, 0x7b, 0xda, 0x7d, 0x16, 0x3a,
			0x08, 0xdb, 0x20, 0xba, 0x64, 0xb6, 0x33, 0x4d, 0xca, 0x34, 0xea, 0xc8, 0x2c, 0xf7,
			0xb4, 0x91, 0xc3, 0x5f, 0x5c, 0xae, 0xc7, 0xb0,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L492
	{"DdzFFzCqrhsy2zYMDQRCF4Nw34C3P7aT5B7JwHFQ6gLAeoHgVXurCLPCm3AeV1nTa1Nd46uDoNt16cnsPFkb4fpLi1J17AmvphCtGFz2",
		[]byte{
			0x0c, 0xd2, 0x15, 0x54, 0xa0, 0xf9, 0xb8, 0x25, 0x9c, 0x46, 0x88, 0xdd, 0x00, 0xfc,
			0x01, 0x88, 0x43, 0x50, 0x79, 0x76, 0x4f, 0xa5, 0x50, 0xfb, 0x57, 0x38, 0x2b, 0xff,
			0x43, 0xe2, 0xd8, 0xd8, 0x27, 0x27, 0x4e, 0x2a, 0x12, 0x9f, 0x86, 0xc3, 0x80, 0x88,
			0x34, 0x37, 0x4d, 0xfe, 0x3f, 0xda, 0xa6, 0x28, 0x48, 0x30, 0xb8, 0xf6, 0xe4, 0x0d,
			0x29, 0x93, 0xde, 0xa2, 0xfb, 0x0a, 0xbe, 0x82,
		},
	}, // https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L505
	{"DdzFFzCqrht8ygB5pLM4uVbS2x4ek2NTDx6R3DJqP7fUaWEkx8RA9UFR8CHitp2R74XLDP876Pe3KLUByHnrWrKWnffpqPpm14rPCxeP",
		[]byte{
			0x1f, 0x0a, 0xb8, 0x33, 0xfd, 0xb1, 0xfa, 0x49, 0x58, 0xce, 0x74, 0x04, 0x81, 0x84,
			0x5b, 0x3a, 0x26, 0x6e, 0xfa, 0xab, 0x2d, 0x65, 0xd1, 0x6b, 0xdd, 0x3d, 0xfe, 0x7f,
			0xcb, 0xe4, 0x46, 0x30, 0x25, 0x9e, 0xd1, 0x91, 0x98, 0x93, 0x03, 0x9d, 0xfd, 0x40,
			0x02, 0x4a, 0x72, 0x03, 0x45, 0x5b, 0x03, 0xd6, 0xd0, 0x0d, 0x0a, 0x5c, 0xd6, 0xee,
			0x82, 0xde, 0x2e, 0xce, 0x73, 0x8a, 0xa1, 0xbf,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L518
	{"DdzFFzCqrhssTywqjv3dw3EakpEydWQcc3phQzR3YF9NPgQN9Ftkx68FfLLnpJ4vhWo9mAjx5EcpM1wNvorSySrpARZGfk5QugHkVs58",
		[]byte{
			0x16, 0xf7, 0xd2, 0x55, 0x32, 0x6d, 0x77, 0x6e, 0xc1, 0xb5, 0xed, 0xd2, 0x5f, 0x75,
			0xd3, 0xe3, 0xeb, 0xe0, 0xb9, 0xd4, 0x9c, 0xdd, 0xb2, 0x46, 0xd8, 0x0c, 0xf4, 0x1b,
			0x25, 0x24, 0x64, 0xb6, 0x24, 0x50, 0xa2, 0x4e, 0xf5, 0x98, 0x7b, 0x4b, 0xd6, 0x5e,
			0x0d, 0x25, 0x23, 0x43, 0xab, 0xa8, 0xef, 0x77, 0x93, 0x34, 0x79, 0xde, 0xa8, 0xdd,
			0xe2, 0x9e, 0xec, 0x56, 0xcc, 0x6a, 0xc0, 0x69,
		},
	},
	// https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L531
	{"DdzFFzCqrhsqTG4t3uq5UBqFrxhxGVM6bvF4q1QcZXqUpizFddEEip7dx5rbife2s9o2fRU3hVKhRp4higog7As8z42s4AMw6Pcu8vL4",
		[]byte{
			0x97, 0xb8, 0x6c, 0x69, 0xd1, 0x2a, 0xf1, 0x64, 0xdc, 0x87, 0xf2, 0x71, 0x26, 0x8f,
			0x33, 0xbc, 0x4d, 0xee, 0xb0, 0xdf, 0xd3, 0x73, 0xc3, 0xfd, 0x3b, 0xac, 0xd4, 0x47,
			0x53, 0xa3, 0x1d, 0xe7, 0x8f, 0x10, 0xe5, 0x55, 0x03, 0x7c, 0xd4, 0x00, 0x43, 0x6c,
			0xcf, 0xd5, 0x38, 0x0d, 0xbb, 0xcd, 0x4d, 0x7c, 0x28, 0x0a, 0xef, 0x9e, 0xc7, 0x57,
			0x4a, 0xe0, 0xac, 0xac, 0x0c, 0xf7, 0x9e, 0x89,
		},
	},
}

func TestFromToBytes(t *testing.T) {
	address, err := FromBytes(base58.Decode(tests[0].AddressBase58))
	if err != nil {
		t.Fatal("Error in test: ", err)
	}
	addr, err := address.ToAddr()
	if err != nil {
		t.Fatal("Error in test: ", err)
	}
	strRepr, err := addr.ToString()
	if err != nil {
		t.Fatal("Error in test: ", err)
	}
	if strRepr != tests[0].AddressBase58 {
		t.Fatal("For test for \"from/to bytes\" unexpected answer")
	}
}

func TestFromBytes(t *testing.T) {
	for i, test := range tests {
		address, err := FromBytes(base58.Decode(test.AddressBase58))
		if err != nil {
			t.Fatal(
				"Error in test", test,
				"error", err,
			)
		}
		if !address.IdenticalWithPubKey((*bip32.XPub)(&test.PublicKey)) {
			t.Fatal("For", test, "unexpected answer")
		} else {
			t.Logf("Test #%d is successful", i+1)
		}
	}
}

func TestFromBytes2(t *testing.T) {
	rawBytes := []byte{130, 216, 24, 88, 66, 131, 88, 28, 98, 20, 93, 160, 196, 223, 73, 74, 239, 128, 24, 81, 94, 84,
		14, 150, 209, 121, 236, 157, 75, 138, 206, 238, 123, 185, 188, 9, 161, 1, 88, 30, 88, 28, 54, 3, 60, 125, 235,
		15, 7, 94, 174, 1, 220, 144, 222, 86, 44, 185, 172, 19, 170, 210, 84, 142, 65, 88, 80, 223, 47, 243, 0, 26,
		103, 3, 88, 25}
	address, err := FromBytes(rawBytes)
	if err != nil {
		t.Fatal(
			"Error in test",
			"error", err,
		)
	}
	testPubKey := []byte{
		0x6a, 0x50, 0x96, 0x89, 0xc6, 0x53, 0x17, 0x58, 0x65, 0x98, 0x5a, 0xd1, 0xe0, 0xeb,
		0x5f, 0xf9, 0xad, 0xa6, 0x99, 0x7a, 0xa4, 0x03, 0xe6, 0x48, 0x61, 0x4b, 0x3b, 0x78,
		0xfc, 0xba, 0x9c, 0x27, 0x30, 0x82, 0x28, 0xd9, 0x87, 0x2a, 0xf8, 0xb6, 0x5b, 0x98,
		0x7f, 0xf2, 0x3e, 0x1a, 0x20, 0xcd, 0x90, 0xd8, 0x34, 0x6c, 0x31, 0xf0, 0xed, 0xb8,
		0x99, 0x89, 0x52, 0xdc, 0x67, 0x66, 0x55, 0x80,
	}
	if !address.IdenticalWithPubKey((*bip32.XPub)(&testPubKey)) {
		t.Fatal("unexpected answer")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L675
func TestVariableNetEncoding(t *testing.T) {
	cases := []uint64{0, 127, 128, 255, 256275757658493284}
	for i, cur_case := range cases {
		encoded := VariableNatEncode(cur_case)
		decoded, _, err := VariableNatDecode(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if cur_case != decoded {
			t.Fatal("unexpected answer in case #", i, " ", cur_case)
		}
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L691
func TestBaseSerializeConsistency(t *testing.T) {
	keyBytes := [crypto.Ed25519KeyHashLen]byte{}
	scriptBytes := [crypto.ScriptHashLen]byte{}
	for i := range keyBytes {
		keyBytes[i] = 23
	}
	for i := range scriptBytes {
		scriptBytes[i] = 42
	}
	addr := NewBaseAddress(5, StakeCredentialFromKeyHash(keyBytes[:]), StakeCredentialFromScriptHash(scriptBytes[:]))
	addrBytes := addr.ToBytes()
	addr2, err := AddressFromBytes(addrBytes)
	if err != nil {
		t.Fatal(err)
	}

	baseAddr, ok := addr2.(*BaseAddress)

	if !ok {
		t.Fatal("unexpected address")
	}
	baseAddrBytes := baseAddr.ToBytes()
	if !bytes.Equal(addrBytes, baseAddrBytes) {
		t.Fatal("unexpected bytes")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L702
func TestPtrSerializeConsistency(t *testing.T) {
	keyBytes := [crypto.Ed25519KeyHashLen]byte{}
	for i := range keyBytes {
		keyBytes[i] = 23
	}
	ptr := NewPointerAddress(25, StakeCredentialFromKeyHash(keyBytes[:]),
		NewPointer(2354556573, 127, 0),
	)
	addrBytes := ptr.ToBytes()
	addr2, err := AddressFromBytes(addrBytes)
	if err != nil {
		t.Fatal(err)
	}
	pointerAddr, ok := addr2.(*PointerAddress)
	if !ok {
		t.Fatal("unexpected address")
	}
	pointerAddrBytes := pointerAddr.ToBytes()
	if !bytes.Equal(addrBytes, pointerAddrBytes) {
		t.Fatal("unexpected bytes")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L713
func TestEnterpiseSerializeConsistency(t *testing.T) {
	keyBytes := [crypto.Ed25519KeyHashLen]byte{}
	for i := range keyBytes {
		keyBytes[i] = 23
	}
	ptr := NewEnterpriseAddress(64, StakeCredentialFromKeyHash(keyBytes[:]))
	addrBytes := ptr.ToBytes()
	addr2, err := AddressFromBytes(addrBytes)
	if err != nil {
		t.Fatal(err)
	}
	enterpriseAddr, ok := addr2.(*EnterpriseAddress)
	if !ok {
		t.Fatal("unexpected address")
	}
	enterpriseAddrBytes := enterpriseAddr.ToBytes()
	if !bytes.Equal(addrBytes, enterpriseAddrBytes) {
		t.Fatal("unexpected bytes")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L723
func TestRewardSerializeConsistency(t *testing.T) {
	keyBytes := [crypto.Ed25519KeyHashLen]byte{}
	for i := range keyBytes {
		keyBytes[i] = 127
	}
	ptr := NewRewardAddress(9, StakeCredentialFromScriptHash(keyBytes[:]))
	addrBytes := ptr.ToBytes()
	addr2, err := AddressFromBytes(addrBytes)
	if err != nil {
		t.Fatal(err)
	}
	rewardAddr, ok := addr2.(*RewardAddress)
	if !ok {
		t.Fatal("unexpected address")
	}
	rewardAddrBytes := rewardAddr.ToBytes()
	if !bytes.Equal(addrBytes, rewardAddrBytes) {
		t.Fatal("unexpected bytes")
	}
}

//implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L732
func rootKey12() bip32.XPrv {
	//test walk nut penalty hip pave soap entry language right filter choice
	entropy := []byte{0xdf, 0x9e, 0xd2, 0x5e, 0xd1, 0x46, 0xbf, 0x43, 0x33, 0x6a, 0x5d, 0x7c, 0xf7, 0x39, 0x59, 0x94}
	return bip32.FromBip39Entropy(entropy, []byte{})
}

//implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L738
func rootKey15() bip32.XPrv {
	//test walk nut penalty hip pave soap entry language right filter choice
	entropy := []byte{0x0c, 0xcb, 0x74, 0xf3, 0x6b, 0x7d, 0xa1, 0x64, 0x9a, 0x81, 0x44, 0x67, 0x55, 0x22, 0xd4, 0xd8, 0x09, 0x7c, 0x64, 0x12}
	return bip32.FromBip39Entropy(entropy, []byte{})
}

//implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L744
func rootKey24() bip32.XPrv {
	//test walk nut penalty hip pave soap entry language right filter choice
	entropy := []byte{0x4e, 0x82, 0x8f, 0x9a, 0x67, 0xdd, 0xcf, 0xf0,
		0xe6, 0x39, 0x1a, 0xd4, 0xf2, 0x6d, 0xdb, 0x75,
		0x79, 0xf5, 0x9b, 0xa1, 0x4b, 0x6d, 0xd4, 0xba,
		0xf6, 0x3d, 0xcf, 0xdb, 0x9d, 0x24, 0x20, 0xda}
	return bip32.FromBip39Entropy(entropy, []byte{})
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L754
func TestBech32Parsing(t *testing.T) {
	addr, err := AddressFromBech32("addr1u8pcjgmx7962w6hey5hhsd502araxp26kdtgagakhaqtq8sxy9w7g")
	if err != nil {
		t.Fatal(err)
	}
	prefix := "foobar"
	encodeStr, err := addr.ToBech32(&prefix)
	if err != nil {
		t.Fatal(err)
	}
	if encodeStr != "foobar1u8pcjgmx7962w6hey5hhsd502araxp26kdtgagakhaqtq8s92n4tm" {
		t.Fatal("unexpected answer")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L760
func TestByronMagicParsing(t *testing.T) {
	addr, err := FromBytes(base58.Decode("Ae2tdPwUPEZ4YjgvykNpoFeYUxoyhNj2kg8KfKWN2FizsSpLUPv68MpTVDo"))
	if err != nil {
		t.Fatal(err)
	}
	byronAddr := addr.ToByronAddress()

	if byronAddr.ProtocolMagic() != MainNet().ProtocolMagic {
		t.Fatal("Wrong the protocol magic for first byron address")
	}
	netId, err := byronAddr.NetworkId()
	if err != nil {
		t.Fatal(err)
	}
	if netId != MainNet().NetworkId {
		t.Fatal("Wrong the protocol magic for first byron address")
	}

	addr2, err := FromBytes(base58.Decode("2cWKMJemoBaipzQe9BArYdo2iPUfJQdZAjm4iCzDA1AfNxJSTgm9FZQTmFCYhKkeYrede"))
	if err != nil {
		t.Fatal(err)
	}
	byronAddr2 := addr2.ToByronAddress()
	if byronAddr2.ProtocolMagic() != TestNet().ProtocolMagic {
		t.Fatal("Wrong the protocol magic for second byron address")
	}
	netId2, err := byronAddr2.NetworkId()
	if err != nil {
		t.Fatal(err)
	}
	if netId2 != TestNet().NetworkId {
		t.Fatal("Wrong the protocol magic for second byron address")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L773
func TestBip32_12Base(t *testing.T) {
	spend := rootKey12().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	stake := rootKey12().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(2).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	stakeHash := stake.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	stakeCred := StakeCredentialFromKeyHash(stakeHash[:])
	addrNet0 := NewBaseAddress(TestNet().NetworkId, spendCred, stakeCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1qz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzer3jcu5d8ps7zex2k2xt3uqxgjqnnj83ws8lhrn648jjxtwq2ytjqp" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewBaseAddress(MainNet().NetworkId, spendCred, stakeCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1qx2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzer3jcu5d8ps7zex2k2xt3uqxgjqnnj83ws8lhrn648jjxtwqfjkjv7" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L797
func TestBip32_12Enterprise(t *testing.T) {
	spend := rootKey12().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewEnterpriseAddress(TestNet().NetworkId, spendCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1vz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzerspjrlsz" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewEnterpriseAddress(MainNet().NetworkId, spendCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1vx2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzers66hrl8" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L813
func TestBip32_12Pointer(t *testing.T) {
	spend := rootKey12().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewPointerAddress(TestNet().NetworkId, spendCred, NewPointer(1, 2, 3))
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1gz2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzerspqgpsqe70et" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewPointerAddress(MainNet().NetworkId, spendCred, NewPointer(24157, 177, 42))
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1gx2fxv2umyhttkxyxp8x0dlpdt3k6cwng5pxj3jhsydzer5ph3wczvf2w8lunk" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L829
func TestBip32_15Base(t *testing.T) {
	spend := rootKey15().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	stake := rootKey15().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(2).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	stakeHash := stake.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	stakeCred := StakeCredentialFromKeyHash(stakeHash[:])
	addrNet0 := NewBaseAddress(TestNet().NetworkId, spendCred, stakeCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1qpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5ewvxwdrt70qlcpeeagscasafhffqsxy36t90ldv06wqrk2qum8x5w" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewBaseAddress(MainNet().NetworkId, spendCred, stakeCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1q9u5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5ewvxwdrt70qlcpeeagscasafhffqsxy36t90ldv06wqrk2qld6xc3" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L853
func TestBip32_15Enterprise(t *testing.T) {
	spend := rootKey15().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewEnterpriseAddress(TestNet().NetworkId, spendCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1vpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg57c2qv" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewEnterpriseAddress(MainNet().NetworkId, spendCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1v9u5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5eg0kvk0f" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L869
func TestBip32_15Pointer(t *testing.T) {
	spend := rootKey15().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewPointerAddress(TestNet().NetworkId, spendCred, NewPointer(1, 2, 3))
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1gpu5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5egpqgpsdhdyc0" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewPointerAddress(MainNet().NetworkId, spendCred, NewPointer(24157, 177, 42))
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1g9u5vlrf4xkxv2qpwngf6cjhtw542ayty80v8dyr49rf5evph3wczvf2kd5vam" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L885
func TestParseRedeemAddress(t *testing.T) {
	extAddr, err := FromBytes(base58.Decode("Ae2tdPwUPEZ3MHKkpT5Bpj549vrRH7nBqYjNXnCV8G2Bc2YxNcGHEa8ykDp"))
	if err != nil {
		t.Fatal(err)
	}
	addr, err := extAddr.ToAddr()
	if err != nil {
		t.Fatal(err)
	}
	bytesAddr, err := addr.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	extAddr2, err := FromBytes(bytesAddr)
	if err != nil {
		t.Fatal(err)
	}
	addr2, err := extAddr2.ToAddr()
	if err != nil {
		t.Fatal(err)
	}
	bytesAddr2, err := addr2.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytesAddr2, bytesAddr) {
		t.Fatal("unexpected answer")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L894
func TestBip32_15Byron(t *testing.T) {
	byronKey := rootKey15().Derive(utils.Harden(44)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	byronAddr := IcarusFromKey(byronKey, MainNet().ProtocolMagic)
	addr, err := byronAddr.ToAddr()
	if err != nil {
		t.Fatal(err)
	}
	addrBase58, err := addr.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if addrBase58 != "Ae2tdPwUPEZHtBmjZBF4YpMkK9tMSPTE2ADEZTPN97saNkhG78TvXdp3GDk" {
		t.Fatal("unexpected base58 address")
	}

	netId, err := byronAddr.NetworkId()
	if err != nil {
		t.Fatal(err)
	}
	if netId != 0b0001 {
		t.Fatal("unexpected network id")
	}

	extAddr2, err := FromBytes(byronAddr.ToBytes())
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := extAddr2.ToAddr()
	if err != nil {
		t.Fatal(err)
	}
	addr2Base58, err := addr2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	if addr2Base58 != addrBase58 {
		t.Fatal("unexpected addresses")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L912
func TestBip32_24Base(t *testing.T) {
	spend := rootKey24().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	stake := rootKey24().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(2).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	stakeHash := stake.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	stakeCred := StakeCredentialFromKeyHash(stakeHash[:])
	addrNet0 := NewBaseAddress(TestNet().NetworkId, spendCred, stakeCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1qqy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmn8k8ttq8f3gag0h89aepvx3xf69g0l9pf80tqv7cve0l33sw96paj" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewBaseAddress(MainNet().NetworkId, spendCred, stakeCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1qyy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmn8k8ttq8f3gag0h89aepvx3xf69g0l9pf80tqv7cve0l33sdn8p3d" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L936
func TestBip32_25Enterprise(t *testing.T) {
	spend := rootKey24().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewEnterpriseAddress(TestNet().NetworkId, spendCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1vqy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmnqtjtf68" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewEnterpriseAddress(MainNet().NetworkId, spendCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1vyy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmnqs6l44z" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L952
func TestBip32_24Pointer(t *testing.T) {
	spend := rootKey24().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(0).Derive(0).Public()
	spendHash := spend.PublicKey().Hash()
	spendCred := StakeCredentialFromKeyHash(spendHash[:])
	addrNet0 := NewPointerAddress(TestNet().NetworkId, spendCred, NewPointer(1, 2, 3))
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "addr_test1gqy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmnqpqgps5mee0p" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewPointerAddress(MainNet().NetworkId, spendCred, NewPointer(24157, 177, 42))
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "addr1gyy6nhfyks7wdu3dudslys37v252w2nwhv0fw2nfawemmnyph3wczvf2dqflgt" {
		t.Fatal("undefined address net 3")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/address.rs#L968
func TestBip32_12Reward(t *testing.T) {
	stakingKey := rootKey12().Derive(utils.Harden(1852)).Derive(utils.Harden(1815)).Derive(utils.Harden(0)).Derive(2).Derive(0).Public()
	stakingKeyHash := stakingKey.PublicKey().Hash()
	stakingKeyCred := StakeCredentialFromKeyHash(stakingKeyHash[:])
	addrNet0 := NewRewardAddress(TestNet().NetworkId, stakingKeyCred)
	addrNet0Bech32, err := addrNet0.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet0Bech32 != "stake_test1uqevw2xnsc0pvn9t9r9c7qryfqfeerchgrlm3ea2nefr9hqp8n5xl" {
		t.Fatal("undefined address net 0")
	}
	addrNet3 := NewRewardAddress(MainNet().NetworkId, stakingKeyCred)
	addrNet3Bech32, err := addrNet3.ToBech32(nil)
	if err != nil {
		t.Fatal(err)
	}
	if addrNet3Bech32 != "stake1uyevw2xnsc0pvn9t9r9c7qryfqfeerchgrlm3ea2nefr9hqxdekzz" {
		t.Fatal("undefined address net 3")
	}
}
