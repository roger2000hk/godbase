package godbase

type Set interface {
	// deletes key from set, starting at offs
	// returns updated set and ok

	Delete(int, Key) (Set, bool)

	// returns true if key exists, starting at offs
	HasKey(int, Key) bool

	// returns index of key, starting at offs; or length if not found
	Index(int, Key) int

	// inserts key into set, starting at offs
	// returns updated set and ok

	Insert(int, Key) (Set, bool)

	// returns number of elems in set
	Len() int64
}
