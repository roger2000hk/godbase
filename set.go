package godbase

type Set interface {
	// deletes key from set
	// returns updated set and ok

	Delete(key Key) (Set, bool)

	// returns val for key and ok
	HasKey(key Key) bool

	// inserts key into set
	// returns updated set and ok

	Insert(key Key) (Set, bool)

	// returns number of elems in set
	Len() int
}
