package maps

// All map keys are requred to support Cmp.

type Cmp interface {
	Less(Cmp) bool
}

// Iters are circular and cheap, since they are nothing but an interface on 
// top of actual nodes.

type Iter interface {
	// Returns true if next elem is not root
	HasNext() bool

	// Returns true if prev elem is not root
	HasPrev() bool

	// Returns key for elem or nil if root
	Key() Cmp

	// Returns iter to next elem
	Next() Iter

	// Returns iter to prev elem
	Prev() Iter

	// Returns val for elem or nil if root
	Val() interface{}
}

// Basic map operations

type Any interface {
	// Deletes all elems matching key/val;
	// val is optional and can be set to nil to delete all elems with key.
	// Returns the number of deleted elems.

	Delete(key Cmp, val interface{}) int

	// Inserts key/val into map;
	// val is optional and can be set to nil, and dup checks can be disabled 
	// by setting allowMulti to false. Returns the inserted val & true on 
	// success, or existing val & false on dup.

	Insert(key Cmp, val interface{}, allowMulti bool) (Iter, bool)

	// Returns the number of elems in map
	Len() int64
}

type TestFn func (Cmp, interface{}) bool

// Operations specific to sorted implementations

type Sorted interface {
	Any

	// Cuts elems from start to end for which fn returns true into new set;
	// start, end & fn are optional. Circular cutting, with start/end on
	// opposite sides of root; is supported.

	Cut(start Iter, end Iter, fn TestFn) Sorted

	// Returns iter for first elem after start >= key;
	// start & key are optional

	First(start Iter, key Cmp) Iter
}
