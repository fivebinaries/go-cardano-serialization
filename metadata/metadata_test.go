package metadata

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/fivebinaries/go-cardano-serialization/hash_map"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/google/go-cmp/cmp"
)

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L809
func TestBinaryEncoding(t *testing.T) {
	inputBytes := make([]byte, 1000)
	for i := range inputBytes {
		inputBytes[i] = byte(i)
	}
	metadata := EncodeArbitraryBytesAsMetadatum(inputBytes)
	outputBytes, err := DecodeArbitraryBytesFromMetadatum(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(inputBytes, outputBytes) {
		t.Fatal("unexpected bytes")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L817
func TestJsonEncodingNoConversions(t *testing.T) {
	inputStr := "{\"int64\": 9223372036854775807, \"uint64\": 18446744073709551615, \"receiver_id\": \"SJKdj34k3jjKFDKfjFUDfdjkfd\",\"sender_id\": \"jkfdsufjdk34h3Sdfjdhfduf873\",\"comment\": \"happy birthday\",\"tags\": [0, 264, -1024, 32]}"
	metadata, err := EncodeJsonStrToMetadatum(inputStr, NoConversions)
	if err != nil {
		t.Fatal(err)
	}
	metadataMap, err := metadata.AsMap()
	metadataMapHash := hash_map.HashMap(metadataMap)
	if err != nil {
		t.Fatal("unexpected metadata")
	}
	if receiverRaw, ok := metadataMapHash.Get(MetadataText("receiver_id")); ok {
		receiverId, err := receiverRaw.(TransactionMetadatum).AsText()
		if err != nil {
			t.Fatal(err)
		}
		if receiverId != "SJKdj34k3jjKFDKfjFUDfdjkfd" {
			t.Fatal("unexpected receiver_id")
		}
	}

	if senderIdRaw, ok := metadataMapHash.Get(MetadataText("sender_id")); ok {
		senderId, err := senderIdRaw.(TransactionMetadatum).AsText()
		if err != nil {
			t.Fatal(err)
		}
		if senderId != "jkfdsufjdk34h3Sdfjdhfduf873" {
			t.Fatal("unexpected sender_id")
		}
	}

	if commentRaw, ok := metadataMapHash.Get(MetadataText("comment")); ok {
		comment, err := commentRaw.(TransactionMetadatum).AsText()
		if err != nil {
			t.Fatal(err)
		}
		if comment != "happy birthday" {
			t.Fatal("unexpected comment")
		}
	}

	if tagsRaw, ok := metadataMapHash.Get(MetadataText("tags")); ok {
		tags, err := tagsRaw.(TransactionMetadatum).AsList()
		if err != nil {
			t.Fatal(err)
		}
		realTags := []int64{0, 264, -1024, 32}
		if len(tags) != len(realTags) {
			t.Fatal("unexpected length of tags")
		}
		for i, tagValue := range tags {
			tagValueInt, err := tagValue.AsInt()
			if err != nil {
				t.Fatal(err)
			}
			if int64(tagValueInt.Value) != realTags[i] {
				t.Fatalf("unexpected tag %d. Expected %d", int64(tagValueInt.Value), realTags[i])
			}
		}
	}

	outputStr, err := DecodeMetadatumToJsonStr(&metadataMap, NoConversions)
	if err != nil {
		t.Fatal("decode failed")
	}
	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputStr), &inputJson)
	if err != nil {
		t.Fatal(err)
	}
	var outputJson map[string]interface{}
	err = json.Unmarshal([]byte(outputStr), &outputJson)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(inputJson, outputJson) {
		t.Fatal("unexpected output json")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L845
func TestJsonEncodingBasic(t *testing.T) {
	inputStr := "{\"0x8badf00d\": \"0xdeadbeef\",\"9\": 5,\"obj\": {\"a\":[{\"5\": 2},{}]}}"
	metadata, err := EncodeJsonStrToMetadatum(inputStr, BasicConversions)
	if err != nil {
		t.Fatalf("encode failed: %s", err)
	}
	jsonEncodingCheckExampleMetadatum(metadata, t)

	outputStr, err := DecodeMetadatumToJsonStr(metadata, BasicConversions)
	if err != nil {
		t.Fatal("decode failed")
	}

	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputStr), &inputJson)
	if err != nil {
		t.Fatal(err)
	}
	var outputJson map[string]interface{}
	err = json.Unmarshal([]byte(outputStr), &outputJson)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(inputJson, outputJson) {
		t.Fatal("unexpected output json")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L845
func TestJsonEncodingDetailed(t *testing.T) {
	inputStr := `{"map":[
            {
                "k":{"bytes":"8badf00d"},
                "v":{"bytes":"deadbeef"}
            },
            {
                "k":{"int":9},
                "v":{"int":5}
            },
            {
                "k":{"string":"obj"},
                "v":{"map":[
                    {
                        "k":{"string":"a"},
                        "v":{"list":[
                        {"map":[
                            {
                                "k":{"int":5},
                                "v":{"int":2}
                            }
                            ]},
                            {"map":[
                            ]}
                        ]}
                    }
                ]}
            }
        ]}`

	metadata, err := EncodeJsonStrToMetadatum(inputStr, DetailedSchema)
	if err != nil {
		t.Fatalf("encode failed: %s", err)
	}

	jsonEncodingCheckExampleMetadatum(metadata, t)

	outputStr, err := DecodeMetadatumToJsonStr(metadata, DetailedSchema)
	if err != nil {
		t.Fatal("decode failed")
	}

	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputStr), &inputJson)
	if err != nil {
		t.Fatal(err)
	}
	var outputJson map[string]interface{}
	err = json.Unmarshal([]byte(outputStr), &outputJson)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(inputJson, outputJson) {
		t.Fatal("unexpected output json")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L883
func jsonEncodingCheckExampleMetadatum(metadata TransactionMetadatum, t *testing.T) {
	metadataMap, err := metadata.AsMap()
	if err != nil {
		t.Fatal("unexpected metadata")
	}
	keyBytes, err := hex.DecodeString("8badf00d")
	if err != nil {
		t.Fatal(err)
	}
	key, err := NewMetadataBytes(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	hashMap := hash_map.HashMap(metadataMap)
	if val, ok := hashMap.Get(key); ok {
		valBytes, err := val.(TransactionMetadatum).AsBytes()
		if err != nil {
			t.Fatal(err)
		}
		checkBytes, err := hex.DecodeString("deadbeef")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(valBytes, checkBytes) {
			t.Fatalf("unexpected bytes for 8badf00d")
		}
	} else {
		t.Fatalf("undefined key: %s", key)
	}

	if val, err := metadataMap.GetI32(9); err == nil {
		if valInt, err := val.AsInt(); err == nil {
			if int32(valInt.Value) != int32(5) {
				t.Fatal("undefined value in map")
			}
		} else {
			t.Fatal("undefined value in map")
		}
	} else {
		t.Fatalf("undefined key: 9")
	}

	if innerMetadata, ok := hashMap.Get(MetadataText("obj")); ok {
		var checkMap1 bool
		var checkMap2 bool
		if innerMetadataMapRaw, err := innerMetadata.(TransactionMetadatum).AsMap(); err == nil {
			innerMetadataMap := hash_map.HashMap(innerMetadataMapRaw)
			if a, ok := innerMetadataMap.Get(MetadataText("a")); ok {
				if aList, err := a.(TransactionMetadatum).AsList(); err == nil && len(aList) == 2 {
					if a1InnerMap, err := aList[0].AsMap(); err == nil {
						if val, err := a1InnerMap.GetI32(5); err == nil {
							if valInt, err := val.AsInt(); err == nil {
								if int32(valInt.Value) == 2 {
									checkMap1 = true
								}
							}
						}
					}

					if a2InnerMapRaw, err := aList[1].AsMap(); err == nil {
						a2InnerMap := hash_map.HashMap(a2InnerMapRaw)
						if a2InnerMap.Count() == 0 {
							checkMap2 = true
						}
					}
				}
			}
		}
		if !checkMap1 || !checkMap2 {
			t.Fatal("undefined inner metadata")
		}
	} else {
		t.Fatal("undefined key: \"obj\"")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L896
func TestJsonEncodingDetailedComplexKey(t *testing.T) {
	inputStr := `{"map":[
            {
            "k":{"list":[
                {"map": [
                    {
                        "k": {"int": 5},
                        "v": {"int": -7}
                    },
                    {
                        "k": {"string": "hello"},
                        "v": {"string": "world"}
                    }
                ]},
                {"bytes": "ff00ff00"}
            ]},
            "v":{"int":5}
            }
        ]}`
	metadata, err := EncodeJsonStrToMetadatum(inputStr, DetailedSchema)
	if err != nil {
		t.Fatal("encode failed")
	}

	if metadataMap, err := metadata.AsMap(); err == nil {
		hashMetadataMap := hash_map.HashMap(metadataMap)
		keys := hashMetadataMap.Keys()
		if len(keys) == 1 {
			if val, ok := hashMetadataMap.GetByHash(keys[0].Hash); ok {
				valInt, err := val.(TransactionMetadatum).AsInt()
				if err != nil || valInt.Value != 5 {
					t.Fatal("unexpected first value in the map")
				}
			} else {
				t.Fatal("unexpected map")
			}

			if keyList, err := keys[0].Value.(TransactionMetadatum).AsList(); err == nil {
				if len(keyList) != 2 {
					t.Fatal("unexpected map")
				}
				if keyMap, err := keyList[0].AsMap(); err == nil {
					if val, err := keyMap.GetI32(5); err == nil {
						if valInt, err := val.AsInt(); err == nil {
							if valInt.IsUnsigned || int64(valInt.Value) != int64(-7) {
								t.Fatal("unexpected map")
							}
						} else {
							t.Fatal("unexpected map")
						}
					} else {
						t.Fatal("unexpected map")
					}
					if val, err := keyMap.GetStr("hello"); err == nil {
						if valText, err := val.AsText(); err == nil {
							if valText != "world" {
								t.Fatal("unexpected map")
							}
						} else {
							t.Fatal("unexpected map")
						}
					} else {
						t.Fatal("unexpected map")
					}
				} else {
					t.Fatal("unexpected map")
				}

				if keyBytes, err := keyList[1].AsBytes(); err == nil {
					checkBytes, err := hex.DecodeString("ff00ff00")
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(keyBytes, checkBytes) {
						t.Fatalf("unexpected bytes for ff00ff00")
					}
				} else {
					t.Fatal("unexpected map")
				}
			} else {
				t.Fatal("unexpected map")
			}

		}
	} else {
		t.Fatal("undefined map")
	}

	outputStr, err := DecodeMetadatumToJsonStr(metadata, DetailedSchema)
	if err != nil {
		t.Fatal("decode failed")
	}

	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputStr), &inputJson)
	if err != nil {
		t.Fatal(err)
	}
	var outputJson map[string]interface{}
	err = json.Unmarshal([]byte(outputStr), &outputJson)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(inputJson, outputJson) {
		t.Fatal("unexpected output json")
	}
}

// implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#L936
func TestAllegraMetadata(t *testing.T) {
	gmd := NewGeneralTransactionMetadata()
	mdatum, err := NewMetadataText("string md")
	if err != nil {
		t.Fatal(err)
	}
	gmd[100] = mdatum
	md1 := NewTransactionMetadata(gmd)
	md1Bytes, err := md1.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	md1Deser, err := TransactionMetadataFromBytes(md1Bytes)
	if err != nil {
		t.Fatal(err)
	}
	md1DeserBytes, err := md1Deser.ToBytes()
	if !bytes.Equal(md1DeserBytes, md1Bytes) {
		t.Fatal("unexpected md1_deser")
	}

	md2 := NewTransactionMetadata(gmd)
	var scripts []types.NativeScript
	scripts = append(scripts, types.NativeScript{V5: &types.InvalidBefore{4, 20}})
	md2.Native = scripts

	md2Bytes, err := md2.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	md2Deser, err := TransactionMetadataFromBytes(md2Bytes)
	if err != nil {
		t.Fatal(err)
	}
	md2DeserBytes, err := md2Deser.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(md2DeserBytes, md2Bytes) {
		t.Fatal(cmp.Diff(md2DeserBytes, md2Bytes))
	}

}
