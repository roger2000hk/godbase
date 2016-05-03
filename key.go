package godbase

import (
	"bytes"
	"github.com/fncodr/godbase/fix"
	"strings"
)

type Key interface {
	Less(Key) bool
}

type BoolKey bool
type FixKey fix.Val
type SetTestFn func (int, Key) bool
type Int64Key int64
type KVMapFn func (Key, interface{}) (Key, interface{})
type KVTestFn func (Key, interface{}) bool
type StrKey string
type StringsKey []string
type TimeKey struct {
	Secs, NSecs int64
}
type UIdKey UId

func (k BoolKey) Less(other Key) bool {
	return !bool(k) && bool(other.(BoolKey))
}

func (_k FixKey) Less(_other Key) bool {
	k, other := fix.Val(_k), fix.Val(_other.(FixKey))
	return k.Cmp(other) < 0
}

func (k Int64Key) Less(other Key) bool {
	return k < other.(Int64Key)
}

func (k StrKey) Less(other Key) bool {
	return strings.Compare(string(k), string(other.(StrKey))) < 0
}

func (k StringsKey) Less(_other Key) bool {
	other := _other.(StringsKey)

	for i, v := range k {
		if res := strings.Compare(v, other[i]); res != 0 {
			return res < 0
		}
	}

	return false
}

func (k TimeKey) Less(_other Key) bool {
	other := _other.(TimeKey)
	return k.Secs < other.Secs || (k.Secs == other.Secs && k.NSecs < other.NSecs)
}

func (k UIdKey) Less(_other Key) bool {
	other := _other.(UIdKey)
	return bytes.Compare(k[:], other[:]) < 0
}
