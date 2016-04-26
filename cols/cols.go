package cols

import (
	"encoding/binary"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/defs"
	"io"
)

type Any interface {
	defs.Any
	ReadVal(io.Reader) error
	WriteVal(interface{}, io.Writer) error
}

type BasicCol struct {
	defs.Basic
}

type Int64 struct {
	BasicCol
}

type UInt64 struct {
	BasicCol
}

func NewInt64(n string) *Int64 {
	return new(Int64).Init(n)
}

func NewUInt64(n string) *UInt64 {
	return new(UInt64).Init(n)
}

func (c *Int64) Init(n string) *Int64 {
	c.Basic.Init(n)
	return c
}

func (c *UInt64) Init(n string) *UInt64 {
	c.Basic.Init(n)
	return c
}

func (c *Int64) ReadVal(r io.Reader) (interface{}, error) {
	var v int64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *UInt64) ReadVal(r io.Reader) (interface{}, error) {
	var v uint64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *Int64) WriteVal(_v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBin(&v, w)
}


func (c *UInt64) WriteVal(_v interface{}, w io.Writer) error {
	v := _v.(uint64)
	return WriteBin(&v, w)
}


func ReadBin(ptr interface{}, r io.Reader) error {
	if err := binary.Read(r, godbase.ByteOrder, ptr); err != nil {
		return err
	}

	return nil
}

func WriteBin(ptr interface{}, w io.Writer) error {
	return binary.Write(w, godbase.ByteOrder, ptr)
}
