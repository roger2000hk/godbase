package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
)

type UIdCol struct {
	Basic
}

type UIdType struct {
	BasicType
}

func NewUId(n string) *UIdCol {
	return new(UIdCol).Init(n)
}

func (c *UIdCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *UIdType) AsKey(_ godbase.Rec, v interface{}) godbase.Key {
	return godbase.UIdKey(v.(godbase.UId))
}

func (_ *UIdType) Hash(_ godbase.Rec, _v interface{}, h hash.Hash64) {
	v := _v.(godbase.UIdKey)
	h.Write(v[:])
}

func (c *UIdCol) Init(n string) *UIdCol {
	c.Basic.Init(n, UId())
	return c
}

func (t *UIdType) Init(n string) *UIdType {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	return t
}

func (_ *UIdType) Read(_ godbase.Rec, _ godbase.ValSize, r io.Reader) (interface{}, error) {
	return godbase.ReadUId(r)
}

func (_ *UIdType) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	v := _v.(godbase.UId)
	return WriteBytes(v[:], w)
}
