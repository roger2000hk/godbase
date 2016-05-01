package tbls

import (
	"bytes"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cxs"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/recs"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	foos := New("foos", 100, nil, 1)

	AddBool(foos, "bool")
	AddFix(foos, "fix", 10)
	AddInt64(foos, "int64")
	AddString(foos, "string")
	AddUId(foos, "uid")

	if c := foos.Col("string"); c.Name() != "string" {
		t.Errorf("invalid col: %v", c)
	}

	i := foos.Cols()

	if c := i.Val().(godbase.Col); c.Name() != "bool" {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(godbase.Col); c.Name() != "fix" {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(godbase.Col); c.Name() != "foos/revision" {
		t.Errorf("invalid col: %v", c)
	}
}

func TestReadWriteRec(t *testing.T) {
	foos := New("foos", 100, nil, 1)

	boolCol := AddBool(foos, "bool")
	fixCol := AddFix(foos, "fix", 10)
	int64Col := AddInt64(foos, "int64")
	stringCol := AddString(foos, "string")
	timeCol := AddTime(foos, "time")
	uidCol := AddUId(foos, "uid")
	
	r := recs.New(nil)
	recs.SetFix(r, fixCol, *fix.New(123, 10))
	recs.SetBool(r, boolCol, true)
	recs.SetInt64(r, int64Col, 1)
	recs.SetString(r, stringCol, "abc")
	recs.SetTime(r, timeCol, time.Now())
	recs.SetUId(r, uidCol, godbase.NewUId())

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
	cx := cxs.New(100)

	foos := New("foos", 100, nil, 1)

	r := recs.New(nil)

	if v := foos.Revision(r); v != -1 {
		t.Errorf("invalid revision before upsert: %v", v)	
	}

	var zt time.Time
	if v := foos.UpsertedAt(r); v != zt {
		t.Errorf("invalid upsert time before upsert: %v", v)	
	}

	foos.Upsert(cx, r)

	if l := foos.Len(); l != 1 {
		t.Errorf("invalid len after upsert: %v", l)	
	}

	if v := foos.Revision(r); v != 0 {
		t.Errorf("invalid revision after upsert: %v", v)	
	}

	if v := foos.UpsertedAt(r); !v.After(zt) {
		t.Errorf("invalid upsert time after upsert: %v", v)	
	}

	if rr, err := foos.Reset(cx.InitRecId(new(recs.Basic), r.Id())); err != nil {
		panic(err)
	} else if !rr.Eq(r) {
		t.Errorf("invalid rec found: %v/%v", rr, r)
	}
}

func TestDumpClearSlurp(t *testing.T) {
	const nrecs = 1000

	cx := cxs.New(100)
	foos := New("foos", 100, nil, 1)

	rs := make([]godbase.Rec, nrecs)

	for i, _ := range rs {
		var err error

		if rs[i], err = foos.Upsert(cx, recs.New(nil)); err != nil {
			panic(err)
		}
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
		if rr, err := foos.Reset(cx.InitRecId(new(recs.Basic), r.Id())); err != nil {
			panic(err)
		} else if !r.Eq(rr) {
			t.Errorf("invalid rec found: %v/%v", rr, r)
		}
	}
}
