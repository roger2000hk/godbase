package godbase

type Map interface {
	// clears all elems from map
	// deallocates nodes for maps that use allocators

	Clear()

	// cuts elems from start to end for which fn returns non nil key into new set
	// start, end & fn are all optional
	// when fn is specified, the returned key/val replaces the original; 
	// except for maps with embedded nodes, where the returned val replaces the
	// entire node
	// no safety checks are provided; if you mess up the ordering, you're on your own
	// circular cuts, with start/end on opposite sides of root; are supported 
	// returns a cut from the start slot for hash maps

	Cut(start, end Iter, fn KVMapFn) Map

	// deletes elems from start to end, matching key/val
	// start, end, key & val are all optional, nil means all elems 
	// specifying iters for hash maps only works within the same slot 
	// circular deletes, with start/end on opposite sides of root; are supported
	// deallocates nodes for slab allocated maps
	// returns an iter to next elem and number of deleted elems

	Delete(start, end Iter, key Key, val interface{}) (Iter, int)

	// returns iter for first elem after start matching key and ok
	// start & val are optional, 
	// specifying a start iter for hash maps only works within the same slot

	Find(start Iter, key Key, val interface{}) (Iter, bool)
	
	// returns iter to first elem
	// not supported by hash maps

	First() Iter

	// returns val for key and ok
	Get(key Key) (interface{}, bool)

	// inserts key/val into map after start
	// start & val are both optional,
	// dup checks can be disabled by setting multi to false
	// returns iter to inserted val & true on success, iter to existing val & false on dup 
	// specifying a start iter for hash maps only works within the same slot

	Insert(start Iter, key Key, val interface{}, multi bool) (Iter, bool)

	// returns a new, empty map of the same type as the receiver
	New() Map

	// returns the number of elems in map
	Len() int64

	// inserts/updates key to val and returns true on insert
	Set(key Key, val interface{}) bool

	// returns string rep of map
	String() string
	
	// calls fn with successive elems until false; returns false on early exit
	While(KVTestFn) bool
}
