package serialization

import (
	"errors"
	"github.com/fxamacker/cbor/v2"
	"hash/crc32"
)

type outerAddress struct {
	_ interface {} `cbor:",toarray"`
	Payload []byte
	Crc uint32
}

type Address struct {
	_           interface {} `cbor:",toarray"`
	Root        []byte
	Attributes  map[uint][]byte
	Type        uint
}

// Deserialize address
func FromBytes(data []byte) (Address, error) {
	var outer outerAddress

	if err := cbor.Unmarshal(data, &outer); err != nil {
		// TODO Encapsulate this error and say we failed to unmarshal outer payload
		return Address{}, err
	}

	if outer.Crc != crc32.ChecksumIEEE(outer.Payload) {
		return Address{}, errors.New("Failed checksum")
	}

	var address Address
	if err := cbor.Unmarshal(outer.Payload, &address); err != nil {
		// TODO Encapsulate this error and say we failed to unmarshal inner payload
		return Address{}, err
	}
	return address, nil
}

func (address Address) ToBytes() ([]byte, error) {
	return cbor.Marshal(address)
}