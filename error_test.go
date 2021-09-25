package mft

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestErrorShow(t *testing.T) {
	err := ErrorSf("abc %v", 5)
	err = err.AppendList(ErrorE(fmt.Errorf("fmt.Errorf")),
		ErrorCSf(222, "ErrorCSf %v", 99).AppendList(ErrorCSf(333, "ErrorCSf 2 %v", 88)))

	s := err.Error()

	sCheck := `abc 5
  [
    fmt.Errorf;
    [222] ErrorCSf 99
      [
        [333] ErrorCSf 2 88
      ]
  ]`

	if s != sCheck {
		t.Fatalf("`%v`\n!=\n`%v`", s, sCheck)
	}
}
