package godbase

type Set interface {
	// returns clone of set
	Clone() Set

	// deletes key from start,
	// returns next idx, or -1 if not found

	Delete(start int, key Key) int

	// deletes keys from start to end (exclusive),
	// returns next idx, or -1 if not found; and nr of deleted elems

	DeleteAll(start, end int, key Key) (int, int64)

	// returns first index of key, from start; or -1 if not found
	First(start int, key Key) int

	// returns elem at index i, within slot for hash sets; key is ignored for sorted sets
	Get(key Key, i int) Key

	// returns last index of key, from start to end (exclusive); or -1 if not found
	Last(start, end int, key Key) int

	// inserts key into set, from start; rejects dup keys if multi=false
	// returns updated set and final index, or org set and -1 if dup

	Insert(start int, key Key, multi bool) (int, bool)

	// returns number of elems in set
	Len() int64

	// calls fn with successive elems until it returns false; returns false on early exit
	While(IKTestFn) bool
}
