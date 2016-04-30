package cols

import (
	"encoding/binary"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/defs"
	"hash"
	"io"
	"math/big"
	"time"
)

type Basic struct {
	defs.Basic
	colType godbase.ColType
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

type FixCol struct {
	Basic
}

type FixType struct {
	BasicType
	denom big.Int
}

type Int64Col struct {
	Basic
}

type Int64Type struct {
	BasicType
}

type RefCol struct {
	UIdCol
}

type RefType struct {
	UIdType
	tbl godbase.Tbl
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
	fixType FixType
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

func Bool() godbase.ColType {
	return &boolType
}

func Fix(d int64) godbase.ColType {
	return new(FixType).Init(d)
}

func CreatedAt() *TimeCol {
	return &createdAt
}

func Int64() godbase.ColType {
	return &int64Type
}

func RecId() *UIdCol {
	return &recId
}

func Ref(tbl godbase.Tbl) godbase.ColType {
	return new(RefType).Init(tbl)
}

func String() godbase.ColType {
	return &stringType
}

func Time() godbase.ColType {
	return &timeType
}

func UId() godbase.ColType {
	return &uidType
}

func NewBool(n string) *BoolCol {
	return new(BoolCol).Init(n)
}

func NewFix(n string, m int64) *FixCol {
	return new(FixCol).Init(n, m)
}

func NewInt64(n string) *Int64Col {
	return new(Int64Col).Init(n)
}

func NewRef(n string, t godbase.Tbl) *RefCol {
	return new(RefCol).Init(n, t)
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

func (c *Basic) AsKey(v interface{}) godbase.Key {
	return c.colType.AsKey(v)
}
	
func (t *BasicType) AsKey(_ interface{}) godbase.Key {
	panic(fmt.Sprintf("AsKey() not supported for %v!", t))
}

func (_ *BoolType) AsKey(v interface{}) godbase.Key {
	return godbase.BoolKey(v.(bool))
}

func (t *FixType) AsKey(_v interface{}) godbase.Key {
	if v, ok := _v.(fix.Val); ok {
		return godbase.FixKey(v)
	}

	var kv fix.Val
	kv.Init(_v.(big.Int), t.denom)
	return godbase.FixKey(kv)
}

func (_ *Int64Type) AsKey(v interface{}) godbase.Key {
	return godbase.Int64Key(v.(int64))
}

func (_ *StringType) AsKey(v interface{}) godbase.Key {
	return godbase.StringKey(v.(string))
}

func (_ *TimeType) AsKey(v interface{}) godbase.Key {
	return godbase.TimeKey(v.(time.Time))
}

func (_ *UIdType) AsKey(v interface{}) godbase.Key {
	return godbase.UIdKey(v.(godbase.UId))
}
		
func (c *Basic) CloneVal(v interface{}) interface{} {
	return c.colType.CloneVal(v)
}

func (_ *BasicType) CloneVal(v interface{}) interface{} {
	return v
}

func (c *FixCol) Denom() big.Int {
	return c.colType.(*FixType).denom
}

func (c *Basic) Decode(v interface{}) interface{} {
	return c.colType.Decode(v)
}

func (c *Basic) Encode(v interface{}) interface{} {
	return c.colType.Encode(v)
}

func (_ *BasicType) Decode(v interface{}) interface{} {
	return v
}

func (t *FixType) Decode(v interface{}) interface{} {
	var res fix.Val
	res.Init(v.(big.Int), t.denom)
	return res
}

func (_ *BasicType) Encode(v interface{}) interface{} {
	return v
}

func (t *FixType) Encode(_v interface{}) interface{} {
	if v, ok := _v.(fix.Val); ok {
		return v.Scale(t.denom.Int64()).Num()
	}

	return _v
}

func (_ *RefType) Encode(v interface{}) interface{} {
	return godbase.UId(v.(godbase.Rec).Id())
}

func (_ *BasicType) Eq(l, r interface{}) bool {
	return l == r
}

func (_ *FixType) Eq(_l, _r interface{}) bool {
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
	v := _v.(godbase.BoolKey)
	godbase.Write(&v, h)
}

func (_ *FixType) Hash(_v interface{}, h hash.Hash64) {
	v := fix.Val(_v.(godbase.FixKey))
	d := v.Num()
	h.Write(d.Bytes())
}

func (_ *Int64Type) Hash(_v interface{}, h hash.Hash64) {
	v := _v.(godbase.Int64Key)
	godbase.Write(&v, h)
}

func (_ *StringType) Hash(v interface{}, h hash.Hash64) {
	h.Write([]byte(v.(godbase.StringKey)))
}

func (_ *TimeType) Hash(_v interface{}, h hash.Hash64) {
	v := time.Time(_v.(godbase.TimeKey)).Unix()
	godbase.Write(&v, h)
}

func (_ *UIdType) Hash(_v interface{}, h hash.Hash64) {
	v := _v.(godbase.UIdKey)
	h.Write(v[:])
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

func (c *BoolCol) Init(n string) *BoolCol {
	c.Basic.Init(n, Bool())
	return c
}

func (c *FixCol) Init(n string, d int64) *FixCol {
	c.Basic.Init(n, new(FixType).Init(d))
	return c
}

func (t *FixType) Init(d int64) *FixType {
	t.BasicType.Init(fmt.Sprintf("Fix(%v)", d))
	t.denom.SetInt64(d)
	return t
}

func (c *Int64Col) Init(n string) *Int64Col {
	c.Basic.Init(n, Int64())
	return c
}

func (c *RefCol) Init(n string, tbl godbase.Tbl) *RefCol {
	c.Basic.Init(n, Ref(tbl))
	return c
}

func (t *RefType) Init(tbl godbase.Tbl) *RefType {
	t.tbl = tbl
	return t
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

func (c *Basic) Read(s godbase.ValSize, r io.Reader) (interface{}, error) {
	return c.colType.Read(s, r)
}

func (_ *BoolType) Read(_ godbase.ValSize, r io.Reader) (interface{}, error) {
	var v byte

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v == 1, nil
}

func (_ *FixType) Read(s godbase.ValSize, r io.Reader) (interface{}, error) {
	bs := make([]byte, s)

	if _, err := io.ReadFull(r, bs); err != nil {
		return nil, err
	}
	
	var v big.Int
	v.SetBytes(bs)
	return v, nil
}

func (_ *Int64Type) Read(_ godbase.ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *StringType) Read(s godbase.ValSize, r io.Reader) (interface{}, error) {
	v := make([]byte, s)

	if _, err := io.ReadFull(r, v); err != nil {
		return nil, err
	}

	return string(v), nil
}

func (_ *TimeType) Read(s godbase.ValSize, r io.Reader) (interface{}, error) {
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

func (_ *UIdType) Read(_ godbase.ValSize, r io.Reader) (interface{}, error) {
	return godbase.ReadUId(r)
}

func (c *RefCol) Tbl() godbase.Tbl {
	return c.colType.(*RefType).tbl
}

func (c *Basic) Type() godbase.ColType {
	return c.colType
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

func (_ *FixType) Write(_v interface{}, w io.Writer) error {
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

func Read(c godbase.Col, r io.Reader) (interface{}, error) {
	var s godbase.ValSize
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
	var s godbase.NameSize

	if err := godbase.Read(&s, r); err != nil {
		return "", err
	}

	v := make([]byte, s)
	if _, err := io.ReadFull(r, v); err != nil {
		return "", err
	}

	return string(v), nil
}

func ReadSize(r io.Reader) (s godbase.ValSize, err error) {
	return s, godbase.Read(&s, r)
}

func WriteBinVal(s godbase.ValSize, ptr interface{}, w io.Writer) error {
	if err := WriteSize(s, w); err != nil {
		return err
	}

	return binary.Write(w, godbase.ByteOrder, ptr)
}

func WriteBytes(v []byte, w io.Writer) error {
	if err := WriteSize(godbase.ValSize(len(v)), w); err != nil {
		return err
	}

	_, err := w.Write(v)
	return err
}

func WriteSize(s godbase.ValSize, w io.Writer) error {
	return godbase.Write(&s, w)
}

func Write(c godbase.Col, v interface{}, w io.Writer) error {
	n := []byte(c.Name())
	s := godbase.NameSize(len(n))

	if err := godbase.Write(&s, w); err != nil {
		return err
	}

	if _, err := w.Write(n); err != nil {
		return err
	}

	return c.Write(v, w)
}
