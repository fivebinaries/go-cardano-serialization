package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fivebinaries/go-cardano-serialization/bip32"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fxamacker/cbor/v2"
	"golang.org/x/crypto/sha3"
	"hash/crc32"
	"log"
)

// AddrAttributes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L80
type AddrAttributes struct {
	DerivationPath []byte
	ProtocolMagic  *uint32
}

// SpendingData implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L383
type SpendingData struct {
	_    interface{} `cbor:",toarray"`
	Type uint64
	Data []byte
}

// hashSpendingDataCbor Special type for https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L178
type hashSpendingDataCbor struct {
	_            interface{} `cbor:",toarray"`
	AddrType     AddrType
	SpendingData SpendingData
	Attributes   AddrAttributesRaw
}

// Special type for marshalling ExtendedAddr
type extendedAddrCBOR struct {
	_          interface{} `cbor:",toarray"`
	Addr       []byte
	Attributes AddrAttributesRaw
	AddrType   AddrType
}

// ExtendedAddr implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L281
type ExtendedAddr struct {
	_          interface{} `cbor:",toarray"`
	Addr       []byte
	Attributes AddrAttributes
	AddrType   AddrType
}

// HashSpendingData implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L178
func HashSpendingData(addrType AddrType, pub bip32.XPub, attributes AddrAttributes) ([]byte, error) {
	cborRaw, err := cbor.Marshal(hashSpendingDataCbor{
		SpendingData: SpendingData{
			Type: SPENDING_DATA_TAG_PUBKEY,
			Data: pub,
		},
		AddrType:   addrType,
		Attributes: attributes.ProcessAttributes(),
	})
	if err != nil {
		log.Fatalf("Error in HashSpendingData: %s", err)
	}
	return SHA3ThenBlake2b224(cborRaw), nil
}

// Convert ExtendedAddr to ByronAddress
func (ea *ExtendedAddr) ToByronAddress() ByronAddress {
	return ByronAddress{*ea}
}

// NewExtendedAddr implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L287
func NewExtendedAddr(pub bip32.XPub, attributes AddrAttributes) *ExtendedAddr {
	hashSpendingData, err := HashSpendingData(ATPubKey, pub, attributes)
	if err != nil {
		log.Fatalf("Error in create new Extended Addr: %s", err)
	}
	return &ExtendedAddr{
		Addr:       hashSpendingData,
		Attributes: attributes,
		AddrType:   ATPubKey,
	}
}

// NewSimpleExtendedAddr implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L295
// bootstrap era + no hdpayload address
func NewSimpleExtendedAddr(pub bip32.XPub, protocolMagic *uint32) *ExtendedAddr {
	return NewExtendedAddr(pub, NewBootstrapEra(nil, protocolMagic))
}

// ToExtendedAddr convert extendedAddrCBOR to ExtendedAddr
func (addr *extendedAddrCBOR) ToExtendedAddr() (ExtendedAddr, error) {
	attr, err := addr.Attributes.ProcessAttributes()
	if err != nil {
		return ExtendedAddr{}, nil
	}
	return ExtendedAddr{
		Attributes: attr,
		AddrType:   addr.AddrType,
		Addr:       addr.Addr,
	}, nil
}

// toExtendedAddrCbor convert ExtendedAddr to extendedAddrCBOR
func (ea *ExtendedAddr) toExtendedAddrCbor() extendedAddrCBOR {
	attr := ea.Attributes.ProcessAttributes()
	return extendedAddrCBOR{
		Attributes: attr,
		AddrType:   ea.AddrType,
		Addr:       ea.Addr,
	}
}

// SHA3ThenBlake2b224 implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L165
// calculate the hash of the data using SHA3 digest then using Blake2b224
func SHA3ThenBlake2b224(data []byte) []byte {
	sh3 := sha3.New256()
	sh3.Write(data)
	sh3Result := sh3.Sum(nil)
	b2result := crypto.Blake2b224(sh3Result)
	return b2result[:]
}

// ProcessAttributes method for converting AddrAttributesRaw to AddrAttributes
func (rawAttributes *AddrAttributesRaw) ProcessAttributes() (AddrAttributes, error) {
	attributes := AddrAttributes{
		DerivationPath: nil,
		ProtocolMagic:  nil,
	}
	for key, value := range *rawAttributes {
		switch key {
		case ATTRIBUTE_NAME_TAG_DERIVATION:
			attributes.DerivationPath = value
		case ATTRIBUTE_NAME_TAG_PROTOCOL_MAGIC:
			var protocolMagic uint32
			switch value[0] {
			case CBORInt8:
				protocolMagic = uint32(value[1])
			case CBORInt16:
				protocolMagic = uint32(binary.BigEndian.Uint16(value[1:]))
			case CBORInt32:
				protocolMagic = binary.BigEndian.Uint32(value[1:])
			case CBORInt64:
				protocolMagic = uint32(binary.BigEndian.Uint64(value[1:]))
			default:
				return AddrAttributes{}, errors.New("unexpected type of integer")
			}
			attributes.ProtocolMagic = &protocolMagic

		default:
			return attributes, errors.New("unknown attributes")
		}
	}
	return attributes, nil
}

// FromBytes Deserialize address
func FromBytes(data []byte) (ExtendedAddr, error) {
	var outer Addr

	if err := cbor.Unmarshal(data, &outer); err != nil {
		return ExtendedAddr{}, fmt.Errorf("cbor unmarshalling error: %w", err)
	}

	if outer.Crc != crc32.ChecksumIEEE(outer.Payload) {
		return ExtendedAddr{}, errors.New("Failed checksum")
	}

	var address extendedAddrCBOR
	if err := cbor.Unmarshal(outer.Payload, &address); err != nil {
		return ExtendedAddr{}, fmt.Errorf("cbor unmarshalling inner payload error: %w", err)
	}

	return address.ToExtendedAddr()
}

// ProcessAttributes method for converting AddrAttributes to AddrAttributesRaw
func (attributes *AddrAttributes) ProcessAttributes() AddrAttributesRaw {
	rawAttributes := make(AddrAttributesRaw)
	if len(attributes.DerivationPath) != 0 {
		rawAttributes[ATTRIBUTE_NAME_TAG_DERIVATION] = attributes.DerivationPath
	}
	if attributes.ProtocolMagic != nil {
		magicBytes := make([]byte, 5)
		magicBytes[0] = CBORInt32
		binary.BigEndian.PutUint32(magicBytes[1:], *attributes.ProtocolMagic)
		rawAttributes[ATTRIBUTE_NAME_TAG_PROTOCOL_MAGIC] = magicBytes
	}
	return rawAttributes
}

// ToBytes convert ExtendedAddr to bytes
// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L346
func (ea *ExtendedAddr) ToBytes() ([]byte, error) {
	return cbor.Marshal(ea.toExtendedAddrCbor())
}

// ToAddr convert ExtendedAddr to Addr
func (ea *ExtendedAddr) ToAddr() (Addr, error) {
	addrBytes, err := ea.ToBytes()
	if err != nil {
		return Addr{}, err
	}
	res := Addr{
		Payload: addrBytes,
	}
	res.CalcCRC()
	return res, nil
}

// MarshalCBOR method for cbor marshalling PayloadBytes
func (pb *PayloadBytes) MarshalCBOR() ([]byte, error) {
	tag := cbor.RawTag{Number: 24}

	tagBytes, err := cbor.Marshal(tag)
	if err != nil {
		return nil, err
	}

	encodedPayload, err := cbor.Marshal([]byte(*pb))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, len(tagBytes)+len(encodedPayload))
	copy(buf, tagBytes)
	copy(buf[len(tagBytes):], encodedPayload)

	return buf, nil
}

// IdenticalWithPubKey implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L201
func (ea *ExtendedAddr) IdenticalWithPubKey(xpub *bip32.XPub) bool {
	newea := NewExtendedAddr(*xpub, ea.Attributes)
	neweaBytes, err := newea.ToBytes()
	if err != nil {
		log.Fatalf("Error in ToBytes: %s", err)
	}
	addrBytes, err := ea.ToBytes()
	if err != nil {
		log.Fatalf("Error in ToBytes: %s", err)
	}
	return bytes.Equal(addrBytes, neweaBytes)
}

// IdenticalWithPubKey implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/legacy_address/address.rs#L85
func NewBootstrapEra(hdap []byte, protocolMagic *uint32) AddrAttributes {
	return AddrAttributes{DerivationPath: hdap, ProtocolMagic: protocolMagic}
}
