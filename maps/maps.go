package maps

// All map keys are requred to support Cmp

type Cmp interface {
	Less(Cmp) bool
}

// Iters are circular and cheap, since they are nothing but a common 
// interface on top of actual nodes. Iters are positioned before start
// on return, call Next() to get first elem.

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

// Basic map ops supported by all implementations

type Any interface {
	// Cuts elems from start to end for which fn returns true into new set;
	// start, end & fn are all optional. Circular cuts, with start/end on
	// opposite sides of root; are supported. Returns a cut from the start slot
	// for hash maps.

	Cut(start, end Iter, fn TestFn) Any

	// Deletes all elems after start matching key/val;
	// start, end, key & val are all optional, nil deletes all elems. Specifying 
	// iters for hash maps only works within the same slot. Circular deletes,
	// with start/end on opposite sides of root; are supported. Returns an iter to next 
	// elem and the number of deleted elems.

	Delete(start, end Iter, key Cmp, val interface{}) (Iter, int)

	// Returns iter for first elem after start matching key and ok;
	// start & val are optional, specifying a start iter for hash maps only works within the 
	// same slot.

	Find(start Iter, key Cmp, val interface{}) (Iter, bool)
	
	// Inserts key/val into map after start;
	// start & val are both optional, dup checks can be disabled by setting allowMulti to false. 
	// Returns iter to inserted val & true on success, or iter to existing val & false on dup. 
	// Specifying a start iter for hash maps only works within the same slot.

	Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool)

	// Returns the number of elems in map
	Len() int64
}

type TestFn func (Cmp, interface{}) bool
