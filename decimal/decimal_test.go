package decimal

import (
	//"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	var res Value
	res.AddInt64(New(1234, 100), 1234, 10)

	if tv := res.Trunc(); tv != 135 {
		t.Errorf("invalid trunc: %v", tv)
	}

	if fv := res.Frac(); fv != 74 {
		t.Errorf("invalid frac: %v", fv)
	}

	if s := res.String(); s != "135.74" {
		t.Errorf("invalid str: %v", s)
	}
}

func TestSub(t *testing.T) {
	var res Value
	res.SubInt64(New(1234, 10), 1234, 100)

	if tv := res.Trunc(); tv != 111 {
		t.Errorf("invalid trunc: %v", tv)
	}

	if fv := res.Frac(); fv != 1 {
		t.Errorf("invalid frac: %v", fv)
	}

	if s := res.String(); s != "111.1" {
		t.Errorf("invalid str: %v", s)
	}
}
