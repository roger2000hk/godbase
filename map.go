package godbase

type Map interface {
	// Clears all elems from map. Deallocates nodes for maps that use allocators.
	Clear()

	// Cuts elems from start to end for which fn returns non nil key into new set;
	// start, end & fn are all optional. When fn is specified, the returned key/val replaces
	// the original; except for maps with embedded nodes, where the returned val replaces the
	// entire node. No safety checks are provided; if you mess up the ordering, you're on your
	// own. Circular cuts, with start/end on opposite sides of root; are supported. 
	// Returns a cut from the start slot for hash maps.

	Cut(start, end Iter, fn KVMapFn) Map

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
	New() Map

	// Returns the number of elems in map
	Len() int64

	// Inserts/updates key to val and returns true on insert
	Set(key Key, val interface{}) bool

	// Returns string repr for printing
	String() string
	
	// Calls fn with successive elems until it returns false; returns false on early exit
	While(KVTestFn) bool
}
