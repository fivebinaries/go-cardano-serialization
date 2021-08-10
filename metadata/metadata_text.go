package metadata

import (
	"errors"
	"github.com/fxamacker/cbor/v2"
)

type MetadataText string

func (m *MetadataText) AsMap() (MetadataMap, error) {
	return MetadataMap{}, errors.New("not a map")
}
func (m *MetadataText) AsList() (MetadataList, error) {
	return MetadataList{}, errors.New("not a list")
}
func (m *MetadataText) AsBytes() (MetadataBytes, error) {
	return MetadataBytes{}, errors.New("not bytes")
}
func (m *MetadataText) AsInt() (MetadataInt, error) {
	return MetadataInt{}, errors.New("not an int")
}
func (m *MetadataText) AsText() (MetadataText, error) {
	return *m, nil
}

func (m *MetadataText) UnmarshalCBOR(bytes []byte) error {
	var tmp string
	err := cbor.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	*m = MetadataText(tmp)
	return nil
}
