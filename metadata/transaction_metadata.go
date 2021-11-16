package metadata

import (
	"errors"

	"github.com/fivebinaries/go-cardano-serialization/common"
	"github.com/fivebinaries/go-cardano-serialization/crypto"
	"github.com/fivebinaries/go-cardano-serialization/hash_map"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fivebinaries/go-cardano-serialization/utils"
	"github.com/fxamacker/cbor/v2"
)

// TransactionMetadatumLabel implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#222
type TransactionMetadatumLabel utils.BigNum

// GeneralTransactionMetadata implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#251
type GeneralTransactionMetadata map[TransactionMetadatumLabel]TransactionMetadatum

// TransactionMetadata implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#289
type TransactionMetadata struct {
	General GeneralTransactionMetadata
	Native  []types.NativeScript
}

// NewGeneralTransactionMetadata implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#257
func NewGeneralTransactionMetadata() GeneralTransactionMetadata {
	return make(GeneralTransactionMetadata)
}

// NewTransactionMetadata implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#310
func NewTransactionMetadata(general GeneralTransactionMetadata) *TransactionMetadata {
	return &TransactionMetadata{General: general, Native: nil}
}

func (t *TransactionMetadata) MarshalCBOR() ([]byte, error) {
	if len(t.Native) == 0 {
		return cbor.Marshal(t.General)
	} else {
		return cbor.Marshal([]interface{}{t.General, t.Native})
	}
}

// ToBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#755
func (t *TransactionMetadata) ToBytes() ([]byte, error) {
	return cbor.Marshal(t)
}

func (g *GeneralTransactionMetadata) ToHashMap() *hash_map.HashMap {
	res := &hash_map.HashMap{}
	for key, value := range *g {
		res.Set(key, value)
	}
	return res
}

func (g *GeneralTransactionMetadata) UnmarshalCBOR(bytes []byte) error {
	var tmp map[TransactionMetadatumLabel]interface{}
	err := cbor.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	*g = make(GeneralTransactionMetadata)
	for k, rawV := range tmp {
		val, err := NewTransactionMetadatum(rawV)
		if err != nil {
			return err
		}
		(*g)[k] = val
	}

	//switch tmp.(type) {
	//case map[interface{}]interface{}:
	//	t.General = make(GeneralTransactionMetadata)
	//	//for k, v := range value {
	//	//	if intK, ok := k.(uint64); ok {
	//	//		t.General[TransactionMetadatumLabel(intK)] = v
	//	//	}
	//	//	t.General[k.(TransactionMetadatumLabel)] = v.(TransactionMetadatum)
	//	//}
	//default:
	//	return errors.New("unexpected bytes in cbor")
	//}
	return nil
}

// TransactionMetadataFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#768
func (t *TransactionMetadata) UnmarshalCBOR(bytes []byte) error {
	var gtm GeneralTransactionMetadata
	err := cbor.Unmarshal(bytes, &gtm)
	if err != nil {
		var full []cbor.RawMessage
		err := cbor.Unmarshal(bytes, &full)
		if err != nil {
			return err
		}
		if len(full) != 2 {
			return errors.New("unexpected array")
		}
		err = cbor.Unmarshal(full[0], &t.General)
		if err != nil {
			return err
		}
		var nativeScripts []cbor.RawMessage
		err = cbor.Unmarshal(full[1], &nativeScripts)
		if err != nil {
			return err
		}
		for _, rawNS := range nativeScripts {
			ns, err := common.DeserializeNativeScript(rawNS)
			if err != nil {
				return err
			}
			t.Native = append(t.Native, ns)
		}
	} else {
		t.General = gtm
		return nil
	}
	return nil
}

// TransactionMetadataFromBytes implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/metadata.rs#768
func TransactionMetadataFromBytes(bytes []byte) (TransactionMetadata, error) {
	var res TransactionMetadata
	err := cbor.Unmarshal(bytes, &res)
	return res, err
}

// HashMetadata implements https://github.com/Emurgo/cardano-serialization-lib/blob/0e89deadf9183a129b9a25c0568eed177d6c6d7c/rust/src/utils.rs#610
func HashMetadata(metadata *TransactionMetadata) (types.MetadataHash, error) {
	metadataBytes, err := metadata.ToBytes()
	if err != nil {
		return nil, err
	}
	b2b := crypto.Blake2b256(metadataBytes)
	return b2b[:], nil
}
