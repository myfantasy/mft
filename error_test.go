package mft

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	errV0 := errors.New("ttt")

	errV1 := ErrorNew("msg1", errV0)
	errV2 := ErrorNew2("msg2", errV0, errV1)

	b1, err := json.Marshal(errV1)
	if err != nil {
		t.Fatal("error v1 fail", err)
	}
	b2, err := json.Marshal(errV2)
	if err != nil {
		t.Fatal("error v1 fail", err)
	}

	if `{"code":50000,"msg":"msg1","iet":"ttt"}` != string(b1) {
		t.Fatal("Error 1 fail", string(b1))
	}

	if `{"code":50000,"msg":"msg2","ie":{"code":50000,"msg":"ttt","ie":{"code":50000,"msg":"msg1","iet":"ttt"}}}` != string(b2) {
		t.Fatal("Error 2 fail", "  ", string(b2))
	}
}
