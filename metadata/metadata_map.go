package metadata

import (
	"errors"
	"github.com/fivebinaries/go-cardano-serialization/hash_map"
)

type MetadataMap hash_map.HashMap

func (m *MetadataMap) UnmarshalCBOR(bytes []byte) error {
	panic("implement me")
}

func (m *MetadataMap) AsMap() (MetadataMap, error) {
	return *m, nil
}
func (m *MetadataMap) AsList() (MetadataList, error) {
	return MetadataList{}, errors.New("not a list")
}
func (m *MetadataMap) AsBytes() (MetadataBytes, error) {
	return MetadataBytes{}, errors.New("not bytes")
}
func (m *MetadataMap) AsInt() (MetadataInt, error) {
	return MetadataInt{}, errors.New("not an int")
}
func (m *MetadataMap) AsText() (MetadataText, error) {
	return "", errors.New("not text")
}

func (m *MetadataMap) GetI32(key int32) (TransactionMetadatum, error) {
	isUnsigned := key > 0
	realKey := MetadataInt{Value: uint64(key), IsUnsigned: isUnsigned}
	hashMap := hash_map.HashMap(*m)
	if val, ok := hashMap.Get(realKey); ok {
		return val.(TransactionMetadatum), nil
	}
	return nil, errors.New("undefined key")
}

func (m *MetadataMap) GetStr(key string) (TransactionMetadatum, error) {
	realKey := MetadataText(key)
	hashMap := hash_map.HashMap(*m)
	if val, ok := hashMap.Get(realKey); ok {
		return val.(TransactionMetadatum), nil
	}
	return nil, errors.New("undefined key")
}
