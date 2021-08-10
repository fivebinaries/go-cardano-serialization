package hash_map

import (
	"encoding/json"
	"github.com/fxamacker/cbor/v2"
	"sort"
)

type Key interface{}
type KeyValue struct {
	Key   Key         `json:"k"`
	Value interface{} `json:"v"`
}
type KeyInfo struct {
	Hash  string
	Value Key
}

type HashMap map[string]KeyValue

func NewHashMap() HashMap {
	return make(HashMap)
}

func (hs *HashMap) Get(key Key) (interface{}, bool) {
	hashKey := hashFunc(key)
	return hs.GetByHash(hashKey)
}

func (hs *HashMap) GetByHash(hashKey string) (interface{}, bool) {
	if res, ok := (*hs)[hashKey]; ok {
		return res.Value, ok
	} else {
		return nil, false
	}
}

func hashFunc(key Key) string {
	r, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}

	return string(r)
}

func (hs *HashMap) Set(key Key, value interface{}) {
	hashKey := hashFunc(key)

	(*hs)[hashKey] = KeyValue{Value: value, Key: key}
}

func (hs *HashMap) Count() int {
	return len(*hs)
}

func (hs *HashMap) MarshalJSON() ([]byte, error) {
	list := make([]KeyValue, 0)
	var keys []string
	for k := range *hs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		list = append(list, (*hs)[key])
	}
	return json.Marshal(list)
}

func (hs *HashMap) MarshalCBOR() ([]byte, error) {
	list := make([]KeyValue, 0)
	var keys []string
	for k := range *hs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		list = append(list, (*hs)[key])
	}
	return cbor.Marshal(list)
}

func (hs *HashMap) Keys() []KeyInfo {
	keys := make([]KeyInfo, 0)
	for k, kv := range *hs {
		keys = append(keys, KeyInfo{
			Hash:  k,
			Value: kv.Key,
		})
	}
	return keys
}
