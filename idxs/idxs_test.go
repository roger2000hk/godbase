package idxs

import (
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
	"testing"
)

// array equality is pretty cool
// as long as the length can be determined at compile time,
// any comparable values can be packed up in a generic array and compared in one op

type testKey2 [2]interface{}
type testKey3 [3]interface{}

func TestKeyEq(t *testing.T) {
	if l, r := (testKey2{"abc", "def"}), (testKey2{"abc", "def"}); l != r {
		t.Errorf("not equal")
	}

	if l, r := (testKey2{"abc", "def"}), (testKey2{"abc", "ghi"}); l == r {
		t.Errorf("equal")
	}

	if l, r := (testKey2{1, 2}), (testKey2{1, 2}); l != r {
		t.Errorf("not equal")
	}

	// mixed types works just as well
	if l, r := (testKey2{"abc", "def"}), (testKey2{1, 2}); l == r {
		t.Errorf("equal")
	}

	if l, r := (testKey2{"abc", 1}), (testKey2{"abc", 1}); l != r {
		t.Errorf("not equal")
	}

	// or how about mixed multi-key indexing?

	m := make(map[interface{}]interface{})

	k2 := testKey2{"abc", "def"}
	m[k2] = 42

	if v := m[k2]; v != 42 {
		t.Errorf("not equal")
	}

	// mixed keys
	k3 := testKey3{1, 2, 3}
	m[k3] = "any value"

	if v := m[k3]; v != "any value" {
		t.Errorf("not equal")
	}
}

func TestUniqueInsertDelete(t *testing.T) {
	foo := cols.NewInt64("foo")
	bar := cols.NewString("bar")
	foobarIdx := NewHash([]cols.Any{foo, bar}, true, 100, nil, 1)

	r := recs.New(nil)
	r.SetInt64(foo, 1)
	r.SetString(bar, "abc")

	if _, err := foobarIdx.Insert(r); err != nil {
		t.Errorf("insert failed: %v", err)
	}

	rr := recs.New(nil)
	rr.SetInt64(foo, 1)
	rr.SetString(bar, "abc")

	if _, err := foobarIdx.Insert(rr); err == nil {
		t.Errorf("dup insert allowed")
	}
}
