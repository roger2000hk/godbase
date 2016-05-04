# godbase
#### a hacker friendly Golang DB

Still cooking, but here's a tiny teaser:

```go

import (
	"bytes"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
	"github.com/fncodr/godbase/tbls"
	"testing"
)

func TestDumpClearSlurp(t *testing.T) {
	const nrecs = 1000

	// create new context
	cx := godbase.NewCx()

	// create tmp tbl named "foos" backed by a 100 slot hashed 1-level map without allocator 

	foos := tbls.New("foos", nil, 100, nil, 1)
	bar := tbls.AddInt64(foos, "bar")

	// fill table with recs

	rs := make([]godbase.Rec, nrecs)

	for i, _ := range rs {
		r := recs.New(godbase.NewUId())
		r.SetInt64(bar, int64(i))
		
		var err error
		if rs[i], err = foos.Upsert(cx, r); err != nil {
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

	if err := foos.Slurp(cx, &buf); err != nil {
		panic(err)
	}

	// make sure everything is back to normal

	if l := foos.Len(); l != nrecs {
		t.Errorf("wrong len after slurp: %v", l)
	}

	for _, r := range rs {
		// Get() returns rec for id or err
		// Eq() compares vals for all cols in receiver with param

		if rr, err := foos.Reset(recs.New(r.Id())); err != nil {
			panic(err)
		} else if !r.Eq(rr) {
			t.Errorf("invalid loaded rec: %v/%v", rr, r)
		}
	}
}

```

Or, if you prefer your examples less toyish; have a look at the reservation system I'm building on top here: https://github.com/fncodr/remento.

### license
NOP