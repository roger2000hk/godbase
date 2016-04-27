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
	if c := i.Val().(cols.Any); c != recs.CreatedAtCol() {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(cols.Any); c != recs.IdCol() {
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

	for i := foos.Cols(); i.Valid(); i = i.Next() {
		c := i.Val().(cols.Any)
		if rr.Get(c) != r.Get(c) {
			t.Errorf("invalid loaded val: %v/%v", rr.Get(c), r.Get(c))
		}
	}
}

func TestUpsert(t *testing.T) {
	foos := New("foos", 100, nil, 1)
	r := foos.Upsert(recs.New(nil))

	if l := foos.Len(); l != 1 {
		t.Errorf("invalid len after upsert: %v", l)	
	}

	if rr, ok := foos.Reset(recs.Get(r.Id(), nil)); !ok || !rr.Eq(r) {
		t.Errorf("invalid loaded rec: %v/%v", rr, r)
	}
}
