package common

import (
	"errors"
	"github.com/fivebinaries/go-cardano-serialization/lib"
	"github.com/fivebinaries/go-cardano-serialization/types"
	"github.com/fxamacker/cbor/v2"
)

//func (t *lib.TimelockStart) UnmarshalCBOR(bytes []byte) error {
//	var tmp []uint32
//	err := cbor.Unmarshal(bytes, &tmp)
//	if err != nil {
//		return err
//	}
//	if len(tmp) != 2 {
//		return errors.New("unexpected array length")
//	}
//	if tmp[0] != 4 {
//		return errors.New("fixed value mismatch")
//	}
//	t.Slot = lib.Slot(tmp[1])
//	return nil
//}
//
//func (t *lib.TimelockExpiry) UnmarshalCBOR(bytes []byte) error {
//	panic("implement me")
//}
//
func DeserializeNativeScript(bytes []byte) (types.NativeScript, error) {
	var ts lib.TimelockStart
	err := cbor.Unmarshal(bytes, &ts)
	if err == nil {
		return types.NativeScript{}, nil
	}

	return types.NativeScript{}, errors.New("unexpected native script")
}

//
//func (s *lib.ScriptPubkey) MarshalCBOR() ([]byte,error) {
//	panic("implement me")
//}
//
//func (s *lib.ScriptAll) MarshalCBOR() ([]byte,error) {
//	panic("implement me")
//}
//
//func (s *lib.ScriptAny) MarshalCBOR() ([]byte,error) {
//	panic("implement me")
//}
//
//func (s *lib.ScriptNOfK) MarshalCBOR() ([]byte,error) {
//	panic("implement me")
//}
//
//func (t *lib.TimelockStart) MarshalCBOR() ([]byte,error) {
//	return cbor.Marshal([]uint64{4, uint64(t.Slot)})
//}
//
//func (t *lib.TimelockExpiry) MarshalCBOR() ([]byte,error) {
//	panic("implement me")
//}
