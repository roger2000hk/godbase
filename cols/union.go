package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
)

type UnionCol struct {
	Basic
}

type UnionType struct {
	BasicType
	typeFn UnionTypeFn
}

func NewUnion(n string, fn UnionTypeFn) *UnionCol {
	return new(UnionCol).Init(n, fn)
}

func (c *UnionCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (t *UnionType) AsKey(r godbase.Rec, v interface{}) godbase.Key {
	return t.typeFn(r).AsKey(r, v)
}

func (t *UnionType) Hash(r godbase.Rec, v interface{}, h hash.Hash64) {
	t.typeFn(r).Hash(r, v, h)
}

func (c *UnionCol) Init(n string, fn UnionTypeFn) *UnionCol {
	c.Basic.Init(n, Union(n, fn))
	return c
}

func (t *UnionType) Init(n string, fn UnionTypeFn) *UnionType {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	t.typeFn = fn
	return t
}

func (t *UnionType) Read(rec godbase.Rec, s godbase.ValSize, r io.Reader) (interface{}, error) {
	return t.typeFn(rec).Read(rec, s, r)
}

func (t *UnionType) Write(r godbase.Rec, v interface{}, w io.Writer) error {
	return t.typeFn(r).Write(r, v, w)
}
