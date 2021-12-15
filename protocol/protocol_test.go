package protocol_test

import (
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/fivebinaries/go-cardano-serialization/protocol"
)

func TestProtocolLoadFromFile(t *testing.T) {
	_, rnFp, _, _ := runtime.Caller(0)
	fp := filepath.Join(path.Dir(path.Dir(rnFp)), "testdata", "protocol", "json", "protocol_parameters.json")
	t.Log(fp)
	p, err := protocol.LoadProtocolFromFile(fp)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(*p, protocol.Protocol{}) {
		t.Fatalf("Got nil protocol, %+v", p)
	}
	t.Logf("LinearFee: %d(tx) + %d", p.TxFeePerByte, p.TxFeeFixed)
}
