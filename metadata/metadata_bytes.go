package metadata

import (
	"errors"
)

type MetadataBytes []byte

func (m *MetadataBytes) UnmarshalCBOR(bytes []byte) error {
	m = (*MetadataBytes)(&bytes)
	return nil
}

func (m *MetadataBytes) AsMap() (MetadataMap, error) {
	return MetadataMap{}, errors.New("not a map")
}
func (m *MetadataBytes) AsList() (MetadataList, error) {
	return MetadataList{}, errors.New("not a list")
}
func (m *MetadataBytes) AsBytes() (MetadataBytes, error) {
	return *m, nil
}
func (m *MetadataBytes) AsInt() (MetadataInt, error) {
	return MetadataInt{}, errors.New("not an int")
}
func (m *MetadataBytes) AsText() (MetadataText, error) {
	return "", errors.New("not text")
}
