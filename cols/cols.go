package cols

import (
	"encoding/binary"
	"fmt"
	"github.com/fncodr/godbase"
	"io"
)

type TypeRegistry map[string]godbase.ColType
type UnionTypeFn func(godbase.Rec) godbase.ColType

var (
	typeRegistry TypeRegistry

	boolType BoolType
	fixType FixType
	int64Type Int64Type
	stringType StrType
	timeType TimeType
	uidType UIdType
	
	recId UIdCol
	createdAt TimeCol
)

func init() {
	typeRegistry = make(TypeRegistry)

	boolType.Init("Bool")
	int64Type.Init("Int64")
	stringType.Init("Str")
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

func GetType(n string) godbase.ColType {
	if t, ok := typeRegistry[n]; ok {
		return t
	}
	
	panic(fmt.Sprintf("invalid col type: %v", n))
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

func Str() godbase.ColType {
	return &stringType
}

func Time() godbase.ColType {
	return &timeType
}

func UId() godbase.ColType {
	return &uidType
}

func Union(n string, fn UnionTypeFn) godbase.ColType {
	return new(UnionType).Init(n, fn)
}

func Read(rec godbase.Rec, c godbase.Col, r io.Reader) (interface{}, error) {
	var s godbase.ValSize
	var err error

	if s, err = ReadSize(r); err != nil {
		return nil, err
	}

	var v interface{}

	if v, err = c.Read(rec, s, r); err != nil {
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

func Write(r godbase.Rec, c godbase.Col, v interface{}, w io.Writer) error {
	n := []byte(c.Name())
	s := godbase.NameSize(len(n))

	if err := godbase.Write(&s, w); err != nil {
		return err
	}

	if _, err := w.Write(n); err != nil {
		return err
	}

	return c.Write(r, v, w)
}
