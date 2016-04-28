package cols

import (
	"encoding/binary"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/maps"
	"hash"
	"io"
	"time"
)

type NameSize uint8
type ValSize uint32


type Any interface {
	defs.Any
	AsKey(interface{}) maps.Key
	CloneVal(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(interface{}, hash.Hash64)
	Read(ValSize, io.Reader) (interface{}, error)
	Write(interface{}, io.Writer) error
}

type Type interface {
	Name() string
	AsKey(Any, interface{}) maps.Key
	CloneVal(Any, interface{}) interface{}
	Eq(Any, interface{}, interface{}) bool
	Hash(Any, interface{}, hash.Hash64)
	Read(Any, ValSize, io.Reader) (interface{}, error)
	Write(Any, interface{}, io.Writer) error	
}

type Basic struct {
	defs.Basic
	colType Type
}

type BasicType struct {
	name string
}

type BoolCol struct {
	Basic
}

type BoolType struct {
	BasicType
}

type Int64Col struct {
	Basic
}

type Int64Type struct {
	BasicType
}

type StringCol struct {
	Basic
}

type StringType struct {
	BasicType
	hash hash.Hash64
}

type TimeCol struct {
	Basic
}

type TimeType struct {
	BasicType
}

type UIdCol struct {
	Basic
	hash hash.Hash64
}

type UIdType struct {
	BasicType
}

var (
	boolType BoolType
	int64Type Int64Type
	stringType StringType
	timeType TimeType
	uidType UIdType
	
	recId UIdCol
	createdAt TimeCol
)

func init() {
	boolType.Init("Bool")
	int64Type.Init("Int64")
	stringType.Init("String")
	timeType.Init("Time")
	uidType.Init("UId")

	recId.Init("godbase/id")
	createdAt.Init("godbase/createdAt")
}

func Bool() Type {
	return &boolType
}

func CreatedAt() *TimeCol {
	return &createdAt
}

func Int64() Type {
	return &int64Type
}

func RecId() *UIdCol {
	return &recId
}

func String() Type {
	return &stringType
}

func Time() Type {
	return &timeType
}

func UId() Type {
	return &uidType
}

func NewBool(n string) *BoolCol {
	return new(BoolCol).Init(n)
}

func NewInt64(n string) *Int64Col {
	return new(Int64Col).Init(n)
}

func NewString(n string) *StringCol {
	return new(StringCol).Init(n)
}

func NewTime(n string) *TimeCol {
	return new(TimeCol).Init(n)
}

func NewUId(n string) *UIdCol {
	return new(UIdCol).Init(n)
}

func (c *Basic) AsKey(v interface{}) maps.Key {
	return c.colType.AsKey(c, v)
}
	
func (_ *BasicType) AsKey(c Any, _ interface{}) maps.Key {
	panic(fmt.Sprintf("AsKey() not supported for %v!", c))
}

func (_ *BoolType) AsKey(_ Any, v interface{}) maps.Key {
	return maps.BoolKey(v.(bool))
}

func (_ *Int64Type) AsKey(_ Any, v interface{}) maps.Key {
	return maps.Int64Key(v.(int64))
}

func (_ *StringType) AsKey(_ Any, v interface{}) maps.Key {
	return maps.StringKey(v.(string))
}

func (_ *TimeType) AsKey(_ Any, v interface{}) maps.Key {
	return maps.TimeKey(v.(time.Time))
}

func (_ *UIdType) AsKey(_ Any, v interface{}) maps.Key {
	return maps.UIdKey(v.(godbase.UId))
}
		
func (c *Basic) CloneVal(v interface{}) interface{} {
	return c.colType.CloneVal(c, v)
}

func (_ *BasicType) CloneVal(_ Any, v interface{}) interface{} {
	return v
}

func (_ *BasicType) Eq(_ Any, l, r interface{}) bool {
	return l == r
}

func (c *Basic) Eq(l, r interface{}) bool {
	return c.colType.Eq(c, l, r)
}

func (c *Basic) Hash(v interface{}, h hash.Hash64) {
	c.colType.Hash(c, v, h)
}

func (_ *BoolType) Hash(_ Any, _v interface{}, h hash.Hash64) {
	v := _v.(maps.BoolKey)
	godbase.Write(&v, h)
}

func (_ *Int64Type) Hash(_ Any, _v interface{}, h hash.Hash64) {
	v := _v.(maps.Int64Key)
	godbase.Write(&v, h)
}

func (_ *StringType) Hash(_ Any, v interface{}, h hash.Hash64) {
	h.Write([]byte(v.(maps.StringKey)))
}

func (_ *TimeType) Hash(_ Any, _v interface{}, h hash.Hash64) {
	v := time.Time(_v.(maps.TimeKey)).Unix()
	godbase.Write(&v, h)
}

func (_ *UIdType) Hash(_ Any, _v interface{}, h hash.Hash64) {
	v := _v.(maps.UIdKey)
	h.Write(v[:])
}

func (c *Basic) Init(n string, ct Type) *Basic {
	c.Basic.Init(n)
	c.colType = ct
	return c
}

func (t *BasicType) Init(n string) *BasicType {
	t.name = n
	return t
}

func (c *BoolCol) Init(n string) *BoolCol {
	c.Basic.Init(n, Bool())
	return c
}

func (c *Int64Col) Init(n string) *Int64Col {
	c.Basic.Init(n, Int64())
	return c
}

func (c *StringCol) Init(n string) *StringCol {
	c.Basic.Init(n, String())
	return c
}

func (c *TimeCol) Init(n string) *TimeCol {
	c.Basic.Init(n, Time())
	return c
}

func (c *UIdCol) Init(n string) *UIdCol {
	c.Basic.Init(n, UId())
	return c
}

func (t *BasicType) Name() string {
	return t.name
}

func (c *Basic) Read(s ValSize, r io.Reader) (interface{}, error) {
	return c.colType.Read(c, s, r)
}

func (_ *BoolType) Read(_ Any, _ ValSize, r io.Reader) (interface{}, error) {
	var v byte

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v == 1, nil
}

func (_ *Int64Type) Read(_ Any, _ ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *StringType) Read(_ Any, s ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (_ *TimeType) Read(_ Any, s ValSize, r io.Reader) (interface{}, error) {
	bs := make([]byte, s)

	if _, err := io.ReadFull(r, bs); err != nil {
		return nil, err
	}

	var v time.Time
	
	if err := v.UnmarshalBinary(bs); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *UIdType) Read(_ Any, _ ValSize, r io.Reader) (interface{}, error) {
	var v godbase.UId

	if _, err := io.ReadFull(r, v[:]); err != nil {
		return nil, err
	}

	return godbase.UId(v), nil
}

func (c *Basic) Write(v interface{}, w io.Writer) error {
	return c.colType.Write(c, v, w)
}

func (_ *BoolType) Write(_ Any, _v interface{}, w io.Writer) error {
	v := byte(0)
	if _v.(bool) {
		v = 1
	} 

	return WriteBinVal(1, &v, w)
}

func (_ *Int64Type) Write(_ Any, _v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}

func (_ *StringType) Write(_ Any, _v interface{}, w io.Writer) error {
	return WriteBytes([]byte(_v.(string)), w)
}

func (_ *TimeType) Write(_ Any, _v interface{}, w io.Writer) error {
	bs, err := _v.(time.Time).MarshalBinary()

	if err != nil {
		return err
	}

	return WriteBytes(bs, w)
}

func (_ *UIdType) Write(_ Any, _v interface{}, w io.Writer) error {
	v := _v.(godbase.UId)
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
