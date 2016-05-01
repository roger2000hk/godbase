package idxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/recs"
	"testing"
	"time"
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
	foobarIdx := NewHash("idx", []godbase.Col{foo, bar}, true, 100, nil, 1)

	r := recs.New(nil)
	recs.SetInt64(r, foo, 1)
	recs.SetString(r, bar, "abc")

	if _, err := foobarIdx.Insert(nil, r); err != nil {
		t.Errorf("insert failed: %v", err)
	}

	if _, err := foobarIdx.Insert(nil, r); err != nil {
		t.Errorf("dup insert of same rec not allowed")
	}

	rr := recs.New(nil)
	recs.SetInt64(rr, foo, 1)
	recs.SetString(rr, bar, "abc")

	if _, err := foobarIdx.Insert(nil, rr); err == nil {
		t.Errorf("dup insert allowed")
	}

	if err := foobarIdx.Delete(nil, r); err != nil {
		t.Errorf("del failed: %v", err)
	}

	if err := foobarIdx.Delete(nil, rr); err == nil {
		t.Errorf("dup del allowed")
	}
}

func TestMultiSort(t *testing.T) {
	date := cols.NewTime("date")
	amount := cols.NewFix("amount", 1000)
	orderIdx := NewSort("idx", []godbase.Col{date, amount}, false, nil, 1)
	d := time.Now().Truncate(time.Hour * 24)

	o1 := recs.New(nil)
	recs.SetTime(o1, date, d)
	recs.SetFix(o1, amount, *fix.New(200, 1))
	orderIdx.Insert(nil, o1)

	o2 := recs.New(nil)
	recs.SetTime(o2, date, d)
	recs.SetFix(o2, amount, *fix.New(300, 1))
	orderIdx.Insert(nil, o2)

	o3 := recs.New(nil)
	recs.SetTime(o3, date, d.AddDate(0, 0, 1))
	recs.SetFix(o3, amount, *fix.New(100, 1))
	orderIdx.Insert(nil, o3)

	i, _ := orderIdx.Find(nil, orderIdx.Key(recs.Time(o3, date)), nil)
	i = i.Next()
	if i.Val() != o3.Id() {
		t.Errorf("invalid find res: %v", i.Key())
	}
}
