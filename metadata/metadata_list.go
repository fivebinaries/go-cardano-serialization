package metadata

import (
	"errors"

	"github.com/fxamacker/cbor/v2"
)

type MetadataList []TransactionMetadatum

func (m *MetadataList) UnmarshalCBOR(bytes []byte) error {
	err := cbor.Unmarshal(bytes, m)
	return err
}

func (m *MetadataList) AsMap() (MetadataMap, error) {
	return MetadataMap{}, errors.New("not a map")
}
func (m *MetadataList) AsList() (MetadataList, error) {
	return *m, nil
}
func (m *MetadataList) AsBytes() (MetadataBytes, error) {
	return MetadataBytes{}, errors.New("not bytes")
}
func (m *MetadataList) AsInt() (MetadataInt, error) {
	return MetadataInt{}, errors.New("not an int")
}
func (m *MetadataList) AsText() (MetadataText, error) {
	return "", errors.New("not text")
}
