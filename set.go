package godbase

type Set interface {
	// deletes key from set, starting at offs
	// returns updated set and idx, or org set and -1 if not found

	Delete(int, Key) (Set, int)

	// returns index of key, starting at offs; or -1 if not found
	Index(int, Key) int

	// inserts key into set, starting at offs
	// returns updated set and final index; or org set and -1 if not found

	Insert(int, Key) (Set, int)

	// returns number of elems in set
	Len() int64
}
