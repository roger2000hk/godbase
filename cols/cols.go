package cols

import (
	"encoding/binary"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/decimal"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/maps"
	"hash"
	"io"
	"math/big"
	"time"
)

type NameSize uint8
type ValSize uint32


type Any interface {
	defs.Any
	AsKey(interface{}) maps.Key
	CloneVal(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(interface{}, hash.Hash64)
	Read(ValSize, io.Reader) (interface{}, error)
	Write(interface{}, io.Writer) error
}

type Type interface {
	Name() string
	AsKey(interface{}) maps.Key
	CloneVal(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(interface{}, hash.Hash64)
	Read(ValSize, io.Reader) (interface{}, error)
	Write(interface{}, io.Writer) error	
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

type DecimalCol struct {
	Basic
}

type DecimalType struct {
	BasicType
	denom big.Int
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
}

type TimeCol struct {
	Basic
}

type TimeType struct {
	BasicType
}

type UIdCol struct {
	Basic
}

type UIdType struct {
	BasicType
}

var (
	boolType BoolType
	decimalType DecimalType
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

func Decimal(d int64) Type {
	return new(DecimalType).Init(d)
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

func NewDecimal(n string, m int64) *DecimalCol {
	return new(DecimalCol).Init(n, m)
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
	return c.colType.AsKey(v)
}
	
func (t *BasicType) AsKey(_ interface{}) maps.Key {
	panic(fmt.Sprintf("AsKey() not supported for %v!", t))
}

func (_ *BoolType) AsKey(v interface{}) maps.Key {
	return maps.BoolKey(v.(bool))
}

func (t *DecimalType) AsKey(_v interface{}) maps.Key {
	if v, ok := _v.(decimal.Value); ok {
		return maps.DecimalKey(v)
	}

	v := _v.(big.Int)
	var kv decimal.Value
	kv.Init(&v, &t.denom)
	return maps.DecimalKey(kv)
}

func (_ *Int64Type) AsKey(v interface{}) maps.Key {
	return maps.Int64Key(v.(int64))
}

func (_ *StringType) AsKey(v interface{}) maps.Key {
	return maps.StringKey(v.(string))
}

func (_ *TimeType) AsKey(v interface{}) maps.Key {
	return maps.TimeKey(v.(time.Time))
}

func (_ *UIdType) AsKey(v interface{}) maps.Key {
	return maps.UIdKey(v.(godbase.UId))
}
		
func (c *Basic) CloneVal(v interface{}) interface{} {
	return c.colType.CloneVal(v)
}

func (_ *BasicType) CloneVal(v interface{}) interface{} {
	return v
}

func (c *DecimalCol) Denom() big.Int {
	return c.colType.(*DecimalType).denom
}

func (c *Basic) Encode(v interface{}) interface{} {
	return c.colType.Encode(v)
}

func (_ *BasicType) Encode(v interface{}) interface{} {
	return v
}

func (t *DecimalType) Encode(_v interface{}) interface{} {
	if v, ok := _v.(decimal.Value); ok {
		return v.Scale(t.denom.Int64()).Num()
	}

	return _v
}

func (_ *BasicType) Eq(l, r interface{}) bool {
	return l == r
}

func (_ *DecimalType) Eq(_l, _r interface{}) bool {
	l, r := _l.(big.Int), _r.(big.Int)
	return l.Cmp(&r) == 0
}

func (c *Basic) Eq(l, r interface{}) bool {
	return c.colType.Eq(l, r)
}

func (c *Basic) Hash(v interface{}, h hash.Hash64) {
	c.colType.Hash(v, h)
}

func (_ *BoolType) Hash(_v interface{}, h hash.Hash64) {
	v := _v.(maps.BoolKey)
	godbase.Write(&v, h)
}

func (_ *DecimalType) Hash(_v interface{}, h hash.Hash64) {
	v := decimal.Value(_v.(maps.DecimalKey))
	d := v.Num()
	h.Write(d.Bytes())
}

func (_ *Int64Type) Hash(_v interface{}, h hash.Hash64) {
	v := _v.(maps.Int64Key)
	godbase.Write(&v, h)
}

func (_ *StringType) Hash(v interface{}, h hash.Hash64) {
	h.Write([]byte(v.(maps.StringKey)))
}

func (_ *TimeType) Hash(_v interface{}, h hash.Hash64) {
	v := time.Time(_v.(maps.TimeKey)).Unix()
	godbase.Write(&v, h)
}

func (_ *UIdType) Hash(_v interface{}, h hash.Hash64) {
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

func (c *DecimalCol) Init(n string, d int64) *DecimalCol {
	c.Basic.Init(n, new(DecimalType).Init(d))
	return c
}

func (t *DecimalType) Init(d int64) *DecimalType {
	t.BasicType.Init(fmt.Sprintf("Decimal(%v)", d))
	t.denom.SetInt64(d)
	return t
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
	return c.colType.Read(s, r)
}

func (_ *BoolType) Read(_ ValSize, r io.Reader) (interface{}, error) {
	var v byte

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v == 1, nil
}

func (_ *DecimalType) Read(s ValSize, r io.Reader) (interface{}, error) {
	bs := make([]byte, s)

	if _, err := io.ReadFull(r, bs); err != nil {
		return nil, err
	}
	
	var v big.Int
	v.SetBytes(bs)
	return v, nil
}

func (_ *Int64Type) Read(_ ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *StringType) Read(s ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (_ *TimeType) Read(s ValSize, r io.Reader) (interface{}, error) {
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

func (_ *UIdType) Read(_ ValSize, r io.Reader) (interface{}, error) {
	var v godbase.UId

	if _, err := io.ReadFull(r, v[:]); err != nil {
		return nil, err
	}

	return godbase.UId(v), nil
}

func (c *Basic) Write(v interface{}, w io.Writer) error {
	return c.colType.Write(v, w)
}

func (_ *BoolType) Write(_v interface{}, w io.Writer) error {
	v := byte(0)
	if _v.(bool) {
		v = 1
	} 

	return WriteBinVal(1, &v, w)
}

func (_ *DecimalType) Write(_v interface{}, w io.Writer) error {
	v := _v.(big.Int)
	return WriteBytes(v.Bytes(), w)
}

func (_ *Int64Type) Write(_v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}

func (_ *StringType) Write(_v interface{}, w io.Writer) error {
	return WriteBytes([]byte(_v.(string)), w)
}

func (_ *TimeType) Write(_v interface{}, w io.Writer) error {
	bs, err := _v.(time.Time).MarshalBinary()

	if err != nil {
		return err
	}

	return WriteBytes(bs, w)
}

func (_ *UIdType) Write(_v interface{}, w io.Writer) error {
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
