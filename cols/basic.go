package cols

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/defs"
	"hash"
	"io"
)

type Basic struct {
	defs.Basic
	colType godbase.ColType
}

type BasicType struct {
	name string
}

func (c *Basic) AsKey(r godbase.Rec, v interface{}) godbase.Key {
	return c.colType.AsKey(r, v)
}
	
func (t *BasicType) AsKey(_ godbase.Rec, _ interface{}) godbase.Key {
	panic(fmt.Sprintf("AsKey() not supported for %v!", t))
}

func (c *Basic) CloneVal(v interface{}) interface{} {
	return c.colType.CloneVal(v)
}

func (_ *BasicType) CloneVal(v interface{}) interface{} {
	return v
}

func (c *Basic) Decode(v interface{}) interface{} {
	return c.colType.Decode(v)
}

func (_ *BasicType) Decode(v interface{}) interface{} {
	return v
}

func (c *Basic) Encode(v interface{}) interface{} {
	return c.colType.Encode(v)
}

func (_ *BasicType) Encode(v interface{}) interface{} {
	return v
}

func (_ *BasicType) Eq(l, r interface{}) bool {
	return l == r
}

func (c *Basic) Eq(l, r interface{}) bool {
	return c.colType.Eq(l, r)
}

func (c *Basic) Hash(r godbase.Rec, v interface{}, h hash.Hash64) {
	c.colType.Hash(r, v, h)
}

func (c *Basic) Init(n string, ct godbase.ColType) *Basic {
	c.Basic.Init(n)
	c.colType = ct
	return c
}

func (t *BasicType) Init(n string) *BasicType {
	t.name = n
	return t
}

func (t *BasicType) Name() string {
	return t.name
}

func (c *Basic) Read(rec godbase.Rec, s godbase.ValSize, r io.Reader) (interface{}, error) {
	return c.colType.Read(rec, s, r)
}

func (c *Basic) Type() godbase.ColType {
	return c.colType
}

func (c *Basic) Write(r godbase.Rec, v interface{}, w io.Writer) error {
	return c.colType.Write(r, v, w)
}
