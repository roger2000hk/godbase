package godbase

import (
	"bytes"
	"github.com/fncodr/godbase/fix"
	"strings"
	"time"
)

type Key interface {
	Less(Key) bool
}

type BoolKey bool
type FixKey fix.Val
type Int64Key int64
type StringKey string
type TimeKey time.Time
type UIdKey UId

func (k BoolKey) Less(other Key) bool {
	return !bool(k) && bool(other.(BoolKey))
}

func (_k FixKey) Less(_other Key) bool {
	k := fix.Val(_k)
	other := fix.Val(_other.(FixKey))
	return k.Cmp(other) < 0
}

func (k Int64Key) Less(other Key) bool {
	return k < other.(Int64Key)
}

func (k StringKey) Less(other Key) bool {
	return strings.Compare(string(k), string(other.(StringKey))) < 0
}

func (k TimeKey) Less(other Key) bool {
	return time.Time(k).Before(time.Time(other.(TimeKey)))
}

func (k UIdKey) Less(_other Key) bool {
	other := _other.(UIdKey)
	return bytes.Compare(k[:], other[:]) < 0
}
