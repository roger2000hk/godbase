package godbase

type Set interface {
	// returns clone of set
	Clone() Set

	// deletes key from start,
	// returns next idx, or -1 if not found

	Delete(start int, key Key) int

	// deletes keys from start to end (exclusive),
	// specify end=0 for rest of set
	// returns next idx and nr of deleted elems

	DeleteAll(start, end int, key Key) (int, int64)

	// returns first index of key, from start; regardless of actually finding key
	// second result is true if key was found
	First(start int, key Key) (int, bool)

	// returns elem at index i, within slot for hash sets; key is ignored for sorted sets
	Get(key Key, i int) Key

	// returns last index of key, from start to end (exclusive); 
	// regardless of actually finding key
	// specify end=0 for rest of set
	// second result is true if key was found

	Last(start, end int, key Key) (int, bool)


	// loads keys into set, from start
	// keys are assumed to be in order and no dup checks are performed
	// for hashed sets; first key decides slot
	Load(start int, keys...Key)

	// inserts key into set, from start; rejects dup keys if multi=false
	// returns updated set and final index, or org set and -1 if dup

	Insert(start int, key Key, multi bool) (int, bool)

	// returns number of elems in set
	Len() int64

	// calls fn with successive elems until it returns false; returns false on early exit
	While(IKTestFn) bool
}
