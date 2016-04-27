package tbls

import (
	"bytes"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	foos := New("foos", 100, nil, 1)

	foos.Add(cols.NewInt64("int64"))
	foos.Add(cols.NewString("string"))
	foos.Add(cols.NewUId("uid"))

	if c := foos.Col("string"); c.Name() != "string" {
		t.Errorf("invalid col: %v", c)
	}

	i := foos.Cols()
	if c := i.Val().(cols.Any); c != cols.CreatedAt() {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(cols.Any); c != cols.RecId() {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(cols.Any); c.Name() != "int64" {
		t.Errorf("invalid col: %v", c)
	}
}

func TestReadWriteRec(t *testing.T) {
	foos := New("foos", 100, nil, 1)

	int64Col := foos.Add(cols.NewInt64("int64")).(*cols.Int64Col)
	stringCol := foos.Add(cols.NewString("string")).(*cols.StringCol)
	timeCol := foos.Add(cols.NewTime("time")).(*cols.TimeCol)
	uidCol := foos.Add(cols.NewUId("uid")).(*cols.UIdCol)
	
	r := recs.New(nil)
	r.SetInt64(int64Col, 1)
	r.SetString(stringCol, "abc")
	r.SetTime(timeCol, time.Now())
	r.SetUId(uidCol, godbase.NewUId())

	var buf bytes.Buffer
	if err := foos.Write(r, &buf); err != nil {
		panic(err)
	}

	rr, err := foos.Read(recs.New(nil), &buf);
	if err != nil {
		panic(err)
	}
	
	if !rr.Eq(r) {
		t.Errorf("invalid loaded val: %v/%v", rr, r)
	}
}

func TestUpsert(t *testing.T) {
	foos := New("foos", 100, nil, 1)
	r, _ := foos.Upsert(recs.New(nil))

	if l := foos.Len(); l != 1 {
		t.Errorf("invalid len after upsert: %v", l)	
	}

	if rr, ok := foos.Reset(recs.Init(r.Id(), nil)); !ok || !rr.Eq(r) {
		t.Errorf("invalid loaded rec: %v/%v", rr, r)
	}
}

func TestDumpClearSlurp(t *testing.T) {
	const nrecs = 1000

	foos := New("foos", 100, nil, 1)


	rs := make([]recs.Any, nrecs)

	for i, _ := range rs {
		rs[i], _ = foos.Upsert(recs.New(nil))
	}

	var buf bytes.Buffer
	if err := foos.Dump(&buf); err != nil {
		panic(err)
	}

	foos.Clear()

	if l := foos.Len(); l != 0 {
		t.Errorf("wrong len after clear: %v", l)
	}

	if err := foos.Slurp(&buf); err != nil {
		panic(err)
	}

	if l := foos.Len(); l != nrecs {
		t.Errorf("wrong len after slurp: %v", l)
	}

	for _, r := range rs {
		if rr, ok := foos.Reset(recs.Init(r.Id(), nil)); !ok || !r.Eq(rr) {
			t.Errorf("invalid loaded rec: %v/%v", rr, r)
		}
	}
}
