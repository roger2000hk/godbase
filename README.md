# godbase
### a hacker friendly Golang DB

still cooking, but here's a tiny teaser:

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

	// Create temp tbl named "foos" backed by a hashed 1-level skip map without allocator 

	foos := tbls.New("foos", 100, nil, 1)
	bar := cols.NewInt64("bar")
	foos.Add(bar)

	// Fill the table with recs
	rs := make([]recs.Any, nrecs)

	for i, _ := range rs {
		r := recs.New(nil)
		r.SetInt64(bar, int64(i))
		
		var ok bool
		if rs[i], ok = foos.Upsert(r); !ok {
			t.Errorf("upsert failed: %v", i)
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
		// Reset() updates specified rec to stored val for all cols in tbl
		// recs.Init() creates a new rec from existing id
		// Eq() compares values for all cols in receiver

		if rr, ok := foos.Reset(recs.Init(r.Id(), nil)); !ok || !r.Eq(rr) {
			t.Errorf("invalid loaded rec: %v/%v", rr, r)
		}
	}
}

```