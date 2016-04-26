package cols

import (
	"encoding/binary"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/defs"
	"io"
)

type ValSize uint64

type Any interface {
	defs.Any
	ReadVal(ValSize, io.Reader) (interface{}, error)
	WriteVal(interface{}, io.Writer) error
}

type BasicCol struct {
	defs.Basic
}

type Int64 struct {
	BasicCol
}

type String struct {
	BasicCol
}

type UInt64 struct {
	BasicCol
}

func NewInt64(n string) *Int64 {
	return new(Int64).Init(n)
}

func NewString(n string) *String {
	return new(String).Init(n)
}

func NewUInt64(n string) *UInt64 {
	return new(UInt64).Init(n)
}

func (c *Int64) Init(n string) *Int64 {
	c.Basic.Init(n)
	return c
}

func (c *String) Init(n string) *String {
	c.Basic.Init(n)
	return c
}

func (c *UInt64) Init(n string) *UInt64 {
	c.Basic.Init(n)
	return c
}

func (c *Int64) ReadVal(s ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *String) ReadVal(s ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (c *UInt64) ReadVal(s ValSize, r io.Reader) (interface{}, error) {
	var v uint64

	if err := ReadBin(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *Int64) WriteVal(_v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}

func (c *String) WriteVal(_v interface{}, w io.Writer) error {
	v := []byte(_v.(string))
	WriteSize(ValSize(len(v)), w)
	_, err := w.Write(v)
	return err
}

func (c *UInt64) WriteVal(_v interface{}, w io.Writer) error {
	v := _v.(uint64)
	return WriteBinVal(8, &v, w)
}

func ReadBin(ptr interface{}, r io.Reader) error {
	if err := binary.Read(r, godbase.ByteOrder, ptr); err != nil {
		return err
	}

	return nil
}

func ReadSize(r io.Reader) (s ValSize, err error) {
	return s, ReadBin(&s, r)
}

func WriteBin(ptr interface{}, w io.Writer) error {
	return binary.Write(w, godbase.ByteOrder, ptr)
}

func WriteBinVal(s ValSize, ptr interface{}, w io.Writer) error {
	if err := WriteSize(s, w); err != nil {
		return err
	}

	return binary.Write(w, godbase.ByteOrder, ptr)
}

func WriteSize(s ValSize, w io.Writer) error {
	return WriteBin(&s, w)
}
