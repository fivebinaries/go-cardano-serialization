package address

import (
	"errors"
	"hash/crc32"

	"github.com/btcsuite/btcutil/base58"
	"github.com/fivebinaries/go-cardano-serialization/network"
	"github.com/fxamacker/cbor/v2"
)

var (
	ErrInvalidByronAddress  = errors.New("invalid byron address")
	ErrInvalidByronChecksum = errors.New("invalid byron checksum")
)

type ByronAddressAttributes struct {
	Payload []byte `cbor:"1,keyasint,omitempty"`
	Network *uint8 `cbor:"2,keyasint,omitempty"`
}

type ByronAddress struct {
	Hash       []byte
	Attributes ByronAddressAttributes
	Tag        uint
}

// Bytes returns byte slice represantation of the Address.
func (b *ByronAddress) Bytes() (bytes []byte) {
	bytes, _ = b.MarshalCBOR()
	return bytes
}

// String returns base58 encoded byron address.
func (b *ByronAddress) String() (str string) {
	return base58.Encode(b.Bytes())
}

// NetworkInfo returns NetworkInfo{ProtocolMagigic and NetworkId}.
func (b *ByronAddress) NetworkInfo() (ni *network.NetworkInfo) {
	if b.Attributes.Network == nil {
		return network.MainNet()
	}
	return network.TestNet()
}

// MarshalCBOR returns a cbor encoded byte slice of the base address.
func (b *ByronAddress) MarshalCBOR() (bytes []byte, err error) {
	raw, err := cbor.Marshal([]interface{}{b.Hash, b.Attributes, b.Tag})
	if err != nil {
		return nil, err
	}
	return cbor.Marshal([]interface{}{
		cbor.Tag{Number: 24, Content: raw},
		uint64(crc32.ChecksumIEEE(raw)),
	})
}

// UnmarshalCBOR deserializes raw byron address, encoded in cbor, into a Byron Address.
func (b *ByronAddress) UnmarshalCBOR(data []byte) error {
	type RawAddr struct {
		_        struct{} `cbor:",toarray"`
		Tag      cbor.Tag
		Checksum uint32
	}

	var rawAddr RawAddr

	if err := cbor.Unmarshal(data, &rawAddr); err != nil {
		return err
	}

	rawTag, ok := rawAddr.Tag.Content.([]byte)
	if !ok || rawAddr.Tag.Number != 24 {
		return ErrInvalidByronAddress
	}

	cheksum := crc32.ChecksumIEEE(rawTag)
	if rawAddr.Checksum != cheksum {
		return ErrInvalidByronChecksum
	}

	var byron struct {
		_      struct{} `cbor:",toarray"`
		Hashed []byte
		Attrs  ByronAddressAttributes
		Tag    uint
	}

	if err := cbor.Unmarshal(rawTag, &byron); err != nil {
		return err
	}

	if len(byron.Hashed) != 28 || byron.Tag != 0 {
		return errors.New("")
	}

	*b = ByronAddress{
		Hash:       byron.Hashed,
		Attributes: byron.Attrs,
		Tag:        byron.Tag,
	}

	return nil
}

//Pref returns the string prefix for the base address. "" for byron address since it has no prefix.
func (b *ByronAddress) Prefix() string {
	return ""
}
