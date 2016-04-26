package cols

import (
	"encoding/binary"
	"github.com/fncodr/godbase"
	"io"
)

type Any interface {
	godbase.Def
	ReadVal(io.Reader) error
	WriteVal(interface{}, io.Writer) error
}

type BasicCol struct {
	godbase.BasicDef
}

type Int64Col struct {
	BasicCol
}

type UInt64Col struct {
	BasicCol
}

func NewInt64Col(n string) *Int64Col {
	return new(Int64Col).Init(n)
}

func NewUInt64Col(n string) *UInt64Col {
	return new(UInt64Col).Init(n)
}

func (c *Int64Col) Init(n string) *Int64Col {
	c.BasicDef.Init(n)
	return c
}

func (c *UInt64Col) Init(n string) *UInt64Col {
	c.BasicDef.Init(n)
	return c
}

func (c *Int64Col) ReadVal(r io.Reader) (interface{}, error) {
	var v int64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *UInt64Col) ReadVal(r io.Reader) (interface{}, error) {
	var v uint64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *Int64Col) WriteVal(_v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBin(&v, w)
}


func (c *UInt64Col) WriteVal(_v interface{}, w io.Writer) error {
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
