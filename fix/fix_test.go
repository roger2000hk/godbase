package fix

import (
	//"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	var res Val
	res.AddInt64(*New(1234, 100), 1234, 10)

	if tr, fr := res.Frac(); tr != 135 || fr != 74 {
		t.Errorf("invalid frac: %v/%v", tr, fr)
	}

	if s := res.String(); s != "135.74" {
		t.Errorf("invalid str: %v", s)
	}
}

func TestCmp(t *testing.T) {
	if res := New(1234, 10).Cmp(*New(12340, 100)); res != 0 {
		t.Errorf("invalid cmp res: %v", res)
	}

	if res := New(12340, 1000).Cmp(*New(1234, 100)); res != 0 {
		t.Errorf("invalid cmp res: %v", res)
	}
}

func TestFloat64(t *testing.T) {
	var v Val
	v.AddFloat64(*New(1234, 100), 12.34)

	if res := v.Float64(); res != 24.68 {
		t.Errorf("invalid float64 res: %v", res)
	}
}

func TestScale(t *testing.T) {
	v := New(1234, 100)

	if res := v.Cmp(*v.Scale(1000)); res != 0 {
		t.Errorf("invalid scale res: %v", res)
	}

	if res := v.Cmp(*v.Scale(10)); res != 0 {
		t.Errorf("invalid scale res: %v", res)
	}
}

func TestSub(t *testing.T) {
	var res Val
	res.SubInt64(*New(1234, 10), 1234, 100)

	if tr, fr := res.Frac(); tr != 111 || fr != 0 {
		t.Errorf("invalid trunc: %v/%v", tr, fr)
	}

	if s := res.String(); s != "111.0" {
		t.Errorf("invalid str: %v", s)
	}
}
