package recs

import (
	"github.com/fncodr/godbase/cols"
	"testing"
)

func TestId(t *testing.T) {
	r1, r2 := New(nil), New(nil)

	if r1.Id() == r2.Id() {
		t.Errorf("equal ids")
	}
}

func TestGetSet(t *testing.T) {
	c := cols.NewInt64("foo")

	r := New(nil)

	if v, ok := r.Find(c); ok {
		t.Errorf("invalid get res from empty rec: %v", v)
	}

	SetInt64(r, c, 1)
	if v := Int64(r, c); v != 1 {
		t.Errorf("invalid int64 res from rec: %v", v)
	}

	SetInt64(r, c, 3)
	if v := Int64(r, c); v != 3 {
		t.Errorf("invalid int64 res from updated rec: %v", v)
	}

	if !r.Delete(c) {
		t.Errorf("delete failed")
	}

	if v, ok := r.Find(c); ok {
		t.Errorf("invalid get res from emptied rec: %v", v)
	}
}
