package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
)

type StrCol struct {
	Basic
}

type StrType struct {
	BasicType
}

func NewStr(n string) *StrCol {
	return new(StrCol).Init(n)
}

func (c *StrCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *StrType) AsKey(_ godbase.Rec, v interface{}) godbase.Key {
	return godbase.StrKey(v.(string))
}

func (_ *StrType) Hash(_ godbase.Rec, v interface{}, h hash.Hash64) {
	h.Write([]byte(v.(godbase.StrKey)))
}

func (c *StrCol) Init(n string) *StrCol {
	c.Basic.Init(n, Str())
	return c
}

func (t *StrType) Init(n string) *StrType {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	return t
}

func (_ *StrType) Read(_ godbase.Rec, s godbase.ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (_ *StrType) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	return WriteBytes([]byte(_v.(string)), w)
}
