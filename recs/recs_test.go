package recs

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"testing"
)

func TestId(t *testing.T) {
	r1, r2 := New(godbase.NewUId()), New(godbase.NewUId())

	if r1.Id() == r2.Id() {
		t.Errorf("equal ids")
	}
}

func TestGetSet(t *testing.T) {
	c := cols.NewInt64("foo")

	r := New(godbase.NewUId())

	if v, ok := r.Find(c); ok {
		t.Errorf("invalid get res from empty rec: %v", v)
	}

	r.SetInt64(c, 1)

	if v := r.Int64(c); v != 1 {
		t.Errorf("invalid int64 res from rec: %v", v)
	}

	r.SetInt64(c, 3)
	if v := r.Int64(c); v != 3 {
		t.Errorf("invalid int64 res from updated rec: %v", v)
	}

	if !r.Delete(c) {
		t.Errorf("delete failed")
	}

	if v, ok := r.Find(c); ok {
		t.Errorf("invalid get res from emptied rec: %v", v)
	}
}
