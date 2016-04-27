package cols

import (
	"encoding/binary"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/defs"
	"io"
)

type NameSize uint8
type ValSize uint32

type Any interface {
	defs.Any
	Read(ValSize, io.Reader) (interface{}, error)
	Write(interface{}, io.Writer) error
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

type UId struct {
	BasicCol
}

func NewInt64(n string) *Int64 {
	return new(Int64).Init(n)
}

func NewString(n string) *String {
	return new(String).Init(n)
}

func NewUId(n string) *UId {
	return new(UId).Init(n)
}

func (c *Int64) Init(n string) *Int64 {
	c.Basic.Init(n)
	return c
}

func (c *String) Init(n string) *String {
	c.Basic.Init(n)
	return c
}

func (c *UId) Init(n string) *UId {
	c.Basic.Init(n)
	return c
}

func (c *Int64) Read(s ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *String) Read(s ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (c *UId) Read(s ValSize, r io.Reader) (interface{}, error) {
	var v [16]byte

	if _, err := io.ReadFull(r, v[:]); err != nil {
		return nil, err
	}

	return godbase.UId(v), nil
}

func (c *Int64) Write(_v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}

func (c *String) Write(_v interface{}, w io.Writer) error {
	return WriteBytes([]byte(_v.(string)), w)
}

func (c *UId) Write(_v interface{}, w io.Writer) error {
	v := [16]byte(_v.(godbase.UId))
	return WriteBytes(v[:], w)
}

func Read(c Any, r io.Reader) (interface{}, error) {
	var s ValSize
	var err error

	if s, err = ReadSize(r); err != nil {
		return nil, err
	}

	var v interface{}

	if v, err = c.Read(s, r); err != nil {
		return nil, err
	}

	return v, nil
}

func ReadName(r io.Reader) (string, error) {
	var s NameSize

	if err := godbase.Read(&s, r); err != nil {
		return "", err
	}

	v := make([]byte, s)
	if _, err := io.ReadFull(r, v); err != nil {
		return "", err
	}

	return string(v), nil
}

func ReadSize(r io.Reader) (s ValSize, err error) {
	return s, godbase.Read(&s, r)
}

func WriteBinVal(s ValSize, ptr interface{}, w io.Writer) error {
	if err := WriteSize(s, w); err != nil {
		return err
	}

	return binary.Write(w, godbase.ByteOrder, ptr)
}

func WriteBytes(v []byte, w io.Writer) error {
	if err := WriteSize(ValSize(len(v)), w); err != nil {
		return err
	}

	_, err := w.Write(v)
	return err
}

func WriteSize(s ValSize, w io.Writer) error {
	return godbase.Write(&s, w)
}

func Write(c Any, v interface{}, w io.Writer) error {
	n := []byte(c.Name())
	s := NameSize(len(n))

	if err := godbase.Write(&s, w); err != nil {
		return err
	}

	if _, err := w.Write(n); err != nil {
		return err
	}

	return c.Write(v, w)
}
