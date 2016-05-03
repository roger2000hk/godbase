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

	// returns first index of key, starting at start; or -1 if not found
	First(start int, key Key) int

	// returns last index of key, between start and end (exclusive); or -1 if not found
	Last(start, end int, key Key) int

	// inserts key into set, starting at start
	// rejects duplicate keys if multi=false and returns updated set and final index; 
	// or org set and -1 if not found

	Insert(start int, key Key, multi bool) (int, bool)

	// returns number of elems in set
	Len() int64
}
