package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
	"time"
)

type TimeCol struct {
	Basic
}

type TimeType struct {
	BasicType
}

func NewTime(n string) *TimeCol {
	return new(TimeCol).Init(n)
}

func (c *TimeCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *TimeType) AsKey(_ godbase.Rec, v interface{}) godbase.Key {
	return godbase.TimeKey(v.(time.Time))
}

func (_ *TimeType) Hash(_ godbase.Rec, _v interface{}, h hash.Hash64) {
	v := time.Time(_v.(godbase.TimeKey))
	s, ns := v.Unix(), v.UnixNano()
	godbase.Write(&s, h)
	godbase.Write(&ns, h)
}

func (c *TimeCol) Init(n string) *TimeCol {
	c.Basic.Init(n, Time())
	return c
}

func (t *TimeType) Init(n string) *TimeType {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	return t
}

func (_ *TimeType) Read(_ godbase.Rec, s godbase.ValSize, r io.Reader) (interface{}, error) {
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

func (_ *TimeType) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	bs, err := _v.(time.Time).MarshalBinary()

	if err != nil {
		return err
	}

	return WriteBytes(bs, w)
}
