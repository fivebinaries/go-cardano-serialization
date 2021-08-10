package metadata

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fivebinaries/go-cardano-serialization/hash_map"
	"github.com/fxamacker/cbor/v2"
	"strconv"
	"strings"
)

// MetadataJsonSchema
// Different schema methods for mapping between JSON and the metadata CBOR.
// This conversion should match TxMetadataJsonSchema in cardano-node defined (at time of writing) here:
// https://github.com/input-output-hk/cardano-node/blob/master/cardano-api/src/Cardano/Api/MetaData.hs
// but has 2 additional schemas for more or less conversionse
// Note: Byte/Strings (including keys) in any schema must be at most 64 bytes in length
type MetadataJsonSchema int

type TransactionMetadatum interface {
	AsMap() (MetadataMap, error)
	AsList() (MetadataList, error)
	AsBytes() (MetadataBytes, error)
	AsInt() (MetadataInt, error)
	AsText() (MetadataText, error)
	cbor.Unmarshaler
}

const (
	MdMaxLen = 64
	// NoConversions
	// Does zero implicit conversions.
	// Round-trip conversions are 100% consistent
	// Treats maps DIRECTLY as maps in JSON in a natural way e.g. {"key1": 47, "key2": [0, 1]]}
	// From JSON:
	// * null/true/false NOT supported.
	// * keys treated as strings only
	// To JSON
	// * Bytes, non-string keys NOT supported.
	// Stricter than any TxMetadataJsonSchema in cardano-node but more natural for JSON -> Metadata
	NoConversions MetadataJsonSchema = 0
	// BasicConversions
	// Does some implicit conversions.
	// Round-trip conversions MD -> JSON -> MD is NOT consistent, but JSON -> MD -> JSON is.
	// Without using bytes
	// Maps are treated as an array of k-v pairs as such: [{"key1": 47}, {"key2": [0, 1]}, {"key3": "0xFFFF"}]
	// From JSON:
	// * null/true/false NOT supported.
	// * Strings parseable as bytes (0x starting hex) or integers are converted.
	// To JSON:
	// * Non-string keys partially supported (bytes as 0x starting hex string, integer converted to string).
	// * Bytes are converted to hex strings starting with 0x for both values and keys.
	// Corresponds to TxMetadataJsonSchema's TxMetadataJsonNoSchema in cardano-node
	BasicConversions MetadataJsonSchema = 1
	// DetailedSchema
	// Supports the annotated schema presented in cardano-node with tagged values e.g. {"int": 7}, {"list": [0, 1]}
	// Round-trip conversions are 100% consistent
	// Maps are treated as an array of k-v pairs as such: [{"key1": {"int": 47}}, {"key2": {"list": [0, 1]}}, {"key3": {"bytes": "0xFFFF"}}]
	// From JSON:
	// * null/true/false NOT supported.
	// * Strings parseable as bytes (hex WITHOUT 0x prefix) or integers converted.
	// To JSON:
	// * Non-string keys are supported. Any key parseable as JSON is encoded as metadata instead of a string
	// Corresponds to TxMetadataJsonSchema's TxMetadataJsonDetailedSchema in cardano-node
	DetailedSchema MetadataJsonSchema = 2
)

func NewTransactionMetadatum(raw interface{}) (TransactionMetadatum, error) {
	switch v := raw.(type) {
	case string:
		return NewMetadataText(v)
	case []interface{}:
		var res MetadataList
		for _, rv := range v {
			val, err := NewTransactionMetadatum(rv)
			if err != nil {
				return nil, err
			}
			res = append(res, val)
		}
		return &res, nil
	case []byte:
		return NewMetadataBytes(v)
	default:
		panic("1!!!!")
	}
}

// NewMetadataBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L152
func NewMetadataBytes(bytes []byte) (TransactionMetadatum, error) {
	if len(bytes) > MdMaxLen {
		return nil, fmt.Errorf("max metadata bytes too long: %d, max = %d", len(bytes), MdMaxLen)
	}
	res := MetadataBytes(bytes)
	return &res, nil
}

// NewMetadataText implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L160
func NewMetadataText(text string) (TransactionMetadatum, error) {
	if len(text) > MdMaxLen {
		return nil, fmt.Errorf("max metadata string too long: %d, max = %d", len(text), MdMaxLen)
	}
	res := MetadataText(text)
	return &res, nil
}

// EncodeArbitraryBytesAsMetadatum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L320
// encodes arbitrary bytes into chunks of 64 bytes (the limit for bytes) as a list to be valid Metadata
func EncodeArbitraryBytesAsMetadatum(bytes []byte) TransactionMetadatum {
	var list MetadataList
	for i := 0; i < len(bytes); i += MdMaxLen {
		end := i + MdMaxLen
		if end > len(bytes) {
			end = len(bytes)
		}
		// this should never fail as we are already chunking it
		newMetadata, _ := NewMetadataBytes(bytes[i:end])
		list = append(list, newMetadata)
	}
	return &list
}

// DecodeArbitraryBytesFromMetadatum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L331
func DecodeArbitraryBytesFromMetadatum(metadata TransactionMetadatum) ([]byte, error) {
	var bytes []byte
	metadataList, err := metadata.AsList()
	if err != nil {
		return nil, err
	}
	for _, elem := range metadataList {
		metadataBytes, err := elem.AsBytes()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, metadataBytes...)
	}
	return bytes, nil
}

func SupportsTaggedValues(schema MetadataJsonSchema) bool {
	switch schema {
	case NoConversions, BasicConversions:
		return false
	case DetailedSchema:
		return true
	}
	return false
}

// EncodeJsonStrToMetadatum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L403
func EncodeJsonStrToMetadatum(jsonString string, schema MetadataJsonSchema) (TransactionMetadatum, error) {
	var value interface{}

	dec := json.NewDecoder(strings.NewReader(jsonString))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return nil, err
	}

	return EncodeJsonValueToMetadatum(value, schema)
}

// hexStringToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L389
func hexStringToBytes(hexstr string) []byte {
	if strings.HasPrefix(hexstr, "0x") {
		dst, err := hex.DecodeString(hexstr[2:])
		if err == nil {
			return dst
		}
	}
	return nil
}

// EncodeJsonValueToMetadatum implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L408
func EncodeJsonValueToMetadatum(value interface{}, schema MetadataJsonSchema) (TransactionMetadatum, error) {
	encodeNumber := func(x json.Number) (TransactionMetadatum, error) {
		resUint64, err := strconv.ParseUint(string(x), 10, 64)
		if err != nil {
			//todo check negative numbers
			i64, err := x.Int64()
			if err != nil {
				return nil, errors.New("floats not allowed in metadata")
			}
			return NewMetadataInt(i64), nil
		} else {
			return NewMetadataUInt(resUint64), nil
		}
	}

	encodeString := func(s string, schema MetadataJsonSchema) (TransactionMetadatum, error) {
		if schema == BasicConversions {
			hexBytes := hexStringToBytes(s)
			if hexBytes != nil {
				return NewMetadataBytes(hexBytes)
			} else {
				res := MetadataText(s)
				return &res, nil
			}
		} else {
			res := MetadataText(s)
			return &res, nil
		}
	}

	encodeArray := func(jsonArray []interface{}, schema MetadataJsonSchema) (TransactionMetadatum, error) {
		var arr MetadataList
		for _, value := range jsonArray {
			res, err := EncodeJsonValueToMetadatum(value, schema)
			if err != nil {
				return nil, err
			}
			arr = append(arr, res)
		}
		return &arr, nil
	}
	switch schema {
	case NoConversions, BasicConversions:
		switch value.(type) {
		case nil:
			return nil, errors.New("null not allowed in metadata")
		case bool:
			return nil, errors.New("bools not allowed in metadata")
		case json.Number:
			number, _ := value.(json.Number)
			return encodeNumber(number)
		case string:
			str, _ := value.(string)
			return encodeString(str, schema)
		case []interface{}:
			arr, _ := value.([]interface{})
			return encodeArray(arr, schema)
		case map[string]interface{}:
			resMap := hash_map.NewHashMap()
			sourceMap, _ := value.(map[string]interface{})
			for rawKey, rawValue := range sourceMap {
				var key TransactionMetadatum
				if schema == BasicConversions {
					dec := json.NewDecoder(strings.NewReader(rawKey))
					dec.UseNumber()
					var tmp interface{}
					if err := dec.Decode(&tmp); err != nil || dec.More() {
						tmp = rawKey
					}
					switch tmp.(type) {
					case json.Number:
						number, _ := tmp.(json.Number)
						tmpKey, err := number.Int64()
						if err != nil {
							return nil, err
						}
						key = NewMetadataInt(tmpKey)
					default:
						var err error
						key, err = encodeString(rawKey, schema)
						if err != nil {
							return nil, err
						}
					}
				} else {
					var err error
					key, err = NewMetadataText(rawKey)
					if err != nil {
						return nil, err
					}
				}

				newValue, err := EncodeJsonValueToMetadatum(rawValue, schema)
				if err != nil {
					return nil, err
				}
				resMap.Set(key, newValue)
			}
			res := MetadataMap(resMap)
			return &res, nil
		}
	case DetailedSchema:
		sourceMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.New("DetailedSchema requires types to be tagged objects")
		}

		if len(sourceMap) == 1 {
			tagMismatch := errors.New("key does not match type")
			for k, v := range sourceMap {
				switch k {
				case "int":
					switch v.(type) {
					case json.Number:
						number, _ := v.(json.Number)
						return encodeNumber(number)
					default:
						return nil, tagMismatch
					}
				case "string":
					str, ok := v.(string)
					if !ok {
						return nil, tagMismatch
					}
					return encodeString(str, schema)
				case "bytes":
					str, ok := v.(string)
					if !ok {
						return nil, tagMismatch
					}
					bytes, err := hex.DecodeString(str)
					if err != nil {
						return nil, tagMismatch
					}
					return NewMetadataBytes(bytes)
				case "list":
					arr, ok := v.([]interface{})
					if !ok {
						return nil, tagMismatch
					}
					return encodeArray(arr, schema)
				case "map":
					resMap := hash_map.NewHashMap()
					mapEntryErr := errors.New("entry format in detailed schema map object not correct. Needs to be of form {\"k\": \"key\", \"v\": value}")
					arr, ok := v.([]interface{})
					if !ok {
						return nil, tagMismatch
					}
					for _, entry := range arr {
						entryObj, ok := entry.(map[string]interface{})
						if !ok {
							return nil, mapEntryErr
						}
						if rawKey, ok := entryObj["k"]; ok {
							if rawValue, ok := entryObj["v"]; ok {
								key, err := EncodeJsonValueToMetadatum(rawKey, schema)
								if err != nil {
									return nil, err
								}
								resValue, err := EncodeJsonValueToMetadatum(rawValue, schema)
								if err != nil {
									return nil, err
								}
								resMap.Set(key, resValue)
							} else {
								return nil, mapEntryErr
							}
						} else {
							return nil, mapEntryErr
						}

					}
					resMetadataMap := MetadataMap(resMap)
					return &resMetadataMap, nil
				default:
					return nil, fmt.Errorf("key '%s' in tagged object not valid", v)
				}
			}
		}
	}
	return nil, errors.New("unknown error")
}

// DecodeMetadatumToJsonStr implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L504
func DecodeMetadatumToJsonStr(metadatum TransactionMetadatum, schema MetadataJsonSchema) (string, error) {
	value, err := DecodeMetadatumToJsonValue(metadatum, schema)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(&value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DecodeMetadatumToJsonValue implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L509
func DecodeMetadatumToJsonValue(metadatum TransactionMetadatum, schema MetadataJsonSchema) (interface{}, error) {
	decodeKey := func(key TransactionMetadatum, schema MetadataJsonSchema) (string, error) {
		if text, err := key.AsText(); err == nil {
			return string(text), nil
		}
		if bytes, err := key.AsBytes(); err == nil && schema != NoConversions {
			return fmt.Sprintf("0x%s", hex.EncodeToString(bytes)), nil
		}
		if bigInt, err := key.AsInt(); err == nil && schema != NoConversions {
			return bigInt.String(), nil
		}
		if list, err := key.AsList(); err == nil && schema == DetailedSchema {
			return DecodeMetadatumToJsonStr(&list, schema)
		}
		if list, err := key.AsMap(); err == nil && schema == DetailedSchema {
			return DecodeMetadatumToJsonStr(&list, schema)
		}
		return "", fmt.Errorf("key type %v not allowed in JSON under specified schema", key)
	}

	var typeKey string
	var value interface{}
	if sourceMap, err := metadatum.AsMap(); err == nil {
		switch schema {
		case NoConversions, BasicConversions:
			// treats maps directly as JSON maps
			jsonMap := make(map[string]interface{})
			for _, kv := range hash_map.HashMap(sourceMap) {
				dKey, err := decodeKey(kv.Key.(TransactionMetadatum), schema)
				if err != nil {
					return nil, err
				}
				dValue, err := DecodeMetadatumToJsonValue(kv.Value.(TransactionMetadatum), schema)
				if err != nil {
					return nil, err
				}
				jsonMap[dKey] = dValue
			}
			typeKey = "map"
			value = jsonMap
		case DetailedSchema:
			resList := &hash_map.HashMap{}
			for _, kv := range hash_map.HashMap(sourceMap) {
				// must encode maps as JSON lists of objects with k/v keys
				// also in these schemas we support more key types than strings
				k, err := DecodeMetadatumToJsonValue(kv.Key.(TransactionMetadatum), schema)
				if err != nil {
					return nil, err
				}
				v, err := DecodeMetadatumToJsonValue(kv.Value.(TransactionMetadatum), schema)
				if err != nil {
					return nil, err
				}
				resList.Set(k, v)
			}
			typeKey = "map"
			value = resList
		}
	}
	if list, err := metadatum.AsList(); err == nil {
		var resList []interface{}
		for _, val := range list {
			dVal, err := DecodeMetadatumToJsonValue(val, schema)
			if err != nil {
				return nil, err
			}
			resList = append(resList, dVal)
		}
		typeKey = "list"
		value = resList
	}
	if bigInt, err := metadatum.AsInt(); err == nil {
		typeKey = "int"
		if bigInt.IsUnsigned {
			value = bigInt.Value
		} else {
			value = int64(bigInt.Value)
		}
	}
	if bytes, err := metadatum.AsBytes(); err == nil {
		typeKey = "bytes"
		switch schema {
		case NoConversions:
			return nil, errors.New("bytes not allowed in JSON in specified schema")

		case BasicConversions: // 0x prefix
			value = fmt.Sprintf("0x%s", hex.EncodeToString(bytes))

		case DetailedSchema: // no prefix
			value = hex.EncodeToString(bytes)
		}
	}
	if text, err := metadatum.AsText(); err == nil {
		typeKey = "string"
		value = text
	}

	// potentially wrap value in a keyed map to represent more types
	if SupportsTaggedValues(schema) {
		resMap := make(map[string]interface{})
		resMap[typeKey] = value
		return resMap, nil
	} else {
		return value, nil
	}
}
