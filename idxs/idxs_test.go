package idxs

import (
	//"fmt"
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
	bar := cols.NewStr("bar")
	foobarIdx := NewHash("idx", []godbase.Col{foo, bar}, true, 100, nil, 1)

	r := recs.New(godbase.NewUId())
	r.SetInt64(foo, 1)
	r.SetStr(bar, "abc")

	if _, err := foobarIdx.Insert(nil, r); err != nil {
		t.Errorf("insert failed: %v", err)
	}

	if _, err := foobarIdx.Insert(nil, r); err != nil {
		t.Errorf("dup insert of same rec not allowed")
	}

	rr := recs.New(godbase.NewUId())
	rr.SetInt64(foo, 1)
	rr.SetStr(bar, "abc")

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

	o1 := recs.New(godbase.NewUId())
	o1.SetTime(date, d)
	o1.SetFix(amount, *fix.New(200, 1))
	orderIdx.Insert(nil, o1)

	o2 := recs.New(godbase.NewUId())
	o2.SetTime(date, d)
	o2.SetFix(amount, *fix.New(300, 1))
	orderIdx.Insert(nil, o2)

	o3 := recs.New(godbase.NewUId())
	o3.SetTime(date, d.AddDate(0, 0, 1))
	o3.SetFix(amount, *fix.New(100, 1))
	orderIdx.Insert(nil, o3)

	i, _ := orderIdx.Find(nil, orderIdx.Key(o3.Time(date)), nil)
	i = i.Next()
	if i.Val() != o3.Id() {
		t.Errorf("invalid find res: %v", i.Key())
	}
}
