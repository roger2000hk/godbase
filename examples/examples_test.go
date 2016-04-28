package examples

import (
	"bytes"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
	"github.com/fncodr/godbase/tbls"
	"testing"
)

func TestDumpClearSlurp(t *testing.T) {
	const nrecs = 1000

	// create tmp tbl named "foos" backed by a hashed 1-level map without allocator 

	foos := tbls.New("foos", 100, nil, 1)
	bar := cols.NewInt64("bar")
	foos.Add(bar)

	// fill table with recs

	rs := make([]recs.Any, nrecs)

	for i, _ := range rs {
		r := recs.New(nil)
		r.SetInt64(bar, int64(i))
		
		var err error
		if rs[i], err = foos.Upsert(r); err != nil {
			panic(err)
		}
	}

	// dump tbl to buffer

	var buf bytes.Buffer
	if err := foos.Dump(&buf); err != nil {
		panic(err)
	}

	// clear recs from tbl

	foos.Clear()

	if l := foos.Len(); l != 0 {
		t.Errorf("wrong len after clear: %v", l)
	}

	// slurp recs from buffer

	if err := foos.Slurp(&buf); err != nil {
		panic(err)
	}

	// make sure everything is back to normal

	if l := foos.Len(); l != nrecs {
		t.Errorf("wrong len after slurp: %v", l)
	}

	for _, r := range rs {
		// Get() returns rec for id or err
		// Eq() compares vals for all cols in receiver with param

		if rr, err := foos.Get(r.Id()); err != nil {
			panic(err)
		} else if !r.Eq(rr) {
			t.Errorf("invalid loaded rec: %v/%v", rr, r)
		}
	}
}
