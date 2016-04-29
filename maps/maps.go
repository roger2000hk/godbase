package maps

import (
	"bytes"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/decimal"
	"strings"
	"time"
)

// All map keys are requred to support Key
type Key interface {
	Less(Key) bool
}

type BoolKey bool
type DecimalKey decimal.Value
type Int64Key int64
type StringKey string
type TimeKey time.Time
type UIdKey godbase.UId

// Iters are circular and cheap, since they are nothing but a common 
// interface on top of actual nodes. 

type Iter interface {
	// Returns key for elem or nil if root
	Key() Key

	// Returns iter to next elem
	Next() Iter

	// Returns iter to prev elem
	Prev() Iter

	// Returns val for elem
	Val() interface{}

	// Returns true if not root
	Valid() bool
}

// Basic map ops supported by all implementations
type Any interface {
	// Clears all elems from map. Deallocates nodes for maps that use allocators.
	Clear()

	// Cuts elems from start to end for which fn returns non nil key into new set;
	// start, end & fn are all optional. When fn is specified, the returned key/val replaces
	// the original; except for maps with embedded nodes, where the returned val replaces the
	// entire node. No safety checks are provided; if you mess up the ordering, you're on your
	// own. Circular cuts, with start/end on opposite sides of root; are supported. 
	// Returns a cut from the start slot for hash maps.

	Cut(start, end Iter, fn MapFn) Any

	// Deletes elems from start to end, matching key/val;
	// start, end, key & val are all optional, nil means all elems. Specifying 
	// iters for hash maps only works within the same slot. Circular deletes,
	// with start/end on opposite sides of root; are supported. Deallocates nodes
	// for maps that use allocators. Returns an iter to next 
	// elem and number of deleted elems.

	Delete(start, end Iter, key Key, val interface{}) (Iter, int)

	// Returns iter for first elem after start matching key and ok;
	// start & val are optional, specifying a start iter for hash maps only works within the 
	// same slot.

	Find(start Iter, key Key, val interface{}) (Iter, bool)
	
	// Returns iter to first elem; not supported by hash maps
	First() Iter

	// Returns val for key and ok
	Get(key Key) (interface{}, bool)

	// Inserts key/val into map after start;
	// start & val are both optional, dup checks can be disabled by setting allowMulti to false. 
	// Returns iter to inserted val & true on success, or iter to existing val & false on dup. 
	// Specifying a start iter for hash maps only works within the same slot.

	Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool)

	// Returns a new, empty map of the same type as the receiver
	New() Any

	// Returns the number of elems in map
	Len() int64

	// Inserts/updates key to val and returns val
	Set(key Key, val interface{}) bool

	// Returns string repr for printing
	String() string
	
	// Calls fn with successive elems until it returns false; returns false on early exit
	While(TestFn) bool
}

// map allocator interface
type Alloc func () Any

type MapFn func (Key, interface{}) (Key, interface{})
type TestFn func (Key, interface{}) bool

func (k BoolKey) Less(other Key) bool {
	return !bool(k) && bool(other.(BoolKey))
}

func (_k DecimalKey) Less(_other Key) bool {
	k := decimal.Value(_k)
	other := decimal.Value(_other.(DecimalKey))
	return k.Cmp(&other) < 0
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
