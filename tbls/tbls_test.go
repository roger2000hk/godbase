package tbls

import (
	"bytes"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
	"testing"
)

func TestCreate(t *testing.T) {
	foos := New("foos")

	foos.Add(cols.NewInt64("int64"))
	foos.Add(cols.NewString("string"))
	foos.Add(cols.NewUId("uid"))

	if c := foos.Col("string"); c.Name() != "string" {
		t.Errorf("invalid col: %v", c)
	}

	i := foos.Cols()
	if c := i.Val().(cols.Any); c != recs.IdCol() {
		t.Errorf("invalid col: %v", c)
	}

	i = i.Next()
	if c := i.Val().(cols.Any); c.Name() != "int64" {
		t.Errorf("invalid col: %v", c)
	}
}

func TestWriteRec(t *testing.T) {
	foos := New("foos")

	int64Col := foos.Add(cols.NewInt64("int64")).(*cols.Int64)
	stringCol := foos.Add(cols.NewString("string")).(*cols.String)
	uidCol := foos.Add(cols.NewUId("uid")).(*cols.UId)
	
	r := recs.New(nil)
	r.SetInt64(int64Col, 1)
	r.SetString(stringCol, "abc")
	r.SetUId(uidCol, godbase.NewUId())

	var buf bytes.Buffer
	if err := foos.Write(r, &buf); err != nil {
		t.Errorf("write error: %v", err)
	}
}
