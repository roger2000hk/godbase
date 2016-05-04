package recs

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/sets"
	"time"
)

type Vals []interface{}

type Basic struct {
	cols sets.Sort
	vals Vals
}

func New(id godbase.UId) *Basic {
	return new(Basic).Init(id)
}

func (self *Basic) Bool(col *cols.BoolCol) bool {
	return self.Get(col).(bool)
}

func (self *Basic) CreatedAt() time.Time {
	return self.Time(cols.CreatedAt())
}

func (self *Basic) Clear() {
	self.cols = sets.Sort{}
	self.vals = Vals{}
}

func (self *Basic) Clone() godbase.Rec {
	return &Basic {
		cols: *self.cols.Clone().(*sets.Sort),
		vals: cloneVals(self.vals) }
}

func (self *Basic) Delete(col godbase.Col) bool {
	if i := self.cols.Delete(0, col); i != -1 {
		self.vals = delVal(self.vals, i)
		return true
	}

	return false
}

func (self *Basic) Eq(_other godbase.Rec) bool {
	other := _other.(*Basic)

	return self.cols.While(func (i int, _c godbase.Key) bool {
		c := _c.(godbase.Col)

		if v, ok := other.Find(c); !ok || !c.Eq(c.Encode(v), self.vals[i]) {
			return false
		}
		
		return true
	})
}


func (self *Basic) Find(col godbase.Col) (interface{}, bool) {
	if i := self.cols.First(0, col); int64(i) < self.cols.Len() && self.cols.Get(nil, i) == col {
		return col.Decode(self.vals[i]), true
	}

	return nil, false
}

func (self *Basic) Fix(col *cols.FixCol) fix.Val {
	return self.Get(col).(fix.Val)
}

func (self *Basic) Get(col godbase.Col) interface{} {
	if v, ok := self.Find(col); ok {
		return v
	}

	panic(fmt.Sprintf("col not found: %v", col.Name()))	
}

func (self *Basic) Id() godbase.UId {
	return self.UId(cols.RecId())	
}

func (self *Basic) Init(id godbase.UId) *Basic {
	Init(self, id)
	return self
}

func (self *Basic) Int64(col *cols.Int64Col) int64 {
	return self.Get(col).(int64)
}

func (self *Basic) Len() int {
	return len(self.vals)
}

func (self *Basic) Ref(col *cols.RefCol, res godbase.Rec) (godbase.Rec, error) {
	res.Set(cols.RecId(), self.Get(col).(godbase.UId)) 
	return col.Tbl().Reset(res)
}

func (self *Basic) Set(col godbase.Col, val interface{}) interface{} {
	i, ok := self.cols.Insert(0, col, false)

	if ok {
		self.vals = insertVal(self.vals, i, col.Encode(val))
	} else {
		self.vals[i] = val
	}

	return val
}

func (self *Basic) SetBool(col *cols.BoolCol, v bool) bool {
	return self.Set(col, v).(bool)
}

func (self *Basic) SetFix(col *cols.FixCol, v fix.Val) fix.Val {
	return self.Set(col, v).(fix.Val)
}

func (self *Basic) SetInt64(col *cols.Int64Col, v int64) int64 {
	return self.Set(col, v).(int64)
}

func (self *Basic) SetRef(col *cols.RefCol, v *Basic) *Basic {
	return self.Set(col, v).(*Basic)
}

func (self *Basic) SetStr(col *cols.StrCol, v string) string {
	return self.Set(col, v).(string)
}

func (self *Basic) SetTime(col *cols.TimeCol, v time.Time) time.Time {
	return self.Set(col, v).(time.Time)
}

func (self *Basic) SetUId(col *cols.UIdCol, v godbase.UId) godbase.UId {
	return self.Set(col, v).(godbase.UId)
}

func (self *Basic) Str(c *cols.StrCol) string {
	return self.Get(c).(string)
}

func (self *Basic) String() string {
	var buf bytes.Buffer
	sep := ""

	self.cols.While(func (i int, c godbase.Key) bool {
		fmt.Fprintf(&buf, "%v%v: %v", sep, c.(godbase.Col).Name(), self.vals[i])
		sep = ", "
		return true
	})

	return buf.String()
}

func (self *Basic) Time(c *cols.TimeCol) time.Time {
	return self.Get(c).(time.Time)
}

func (self *Basic) UId(c *cols.UIdCol) godbase.UId {
	return self.Get(c).(godbase.UId)
}

func (self *Basic) While(fn godbase.ColValTestFn) bool {
	return self.cols.While(func (i int, c godbase.Key) bool {
		if !fn(c.(godbase.Col), self.vals[i]) {
			return false
		}

		return true
	})
}

func cloneVals(in Vals) Vals {
	out := make(Vals, len(in))
	copy(out, in)
	return out
}

func delVal(in Vals, i int) Vals {
	copy(in[i:], in[i+1:])
	return in[:len(in)-1]
}

func insertVal(in Vals, i int, val interface{}) Vals {
	if i < len(in) {
		out := append(in, nil)
		copy(out[i+1:], out[i:])
		out[i] = val
		return out
	}

	return append(in, val) 
}
