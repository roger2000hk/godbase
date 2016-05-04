package godbase

// iters are circular and cheap, 
// since they are nothing but a common interface on top of actual nodes

type Iter interface {
	// returns key for elem or nil if root
	Key() Key

	// returns iter to next elem
	Next() Iter

	// returns iter to next elem
	Prev() Iter

	// returns val for elem
	Val() interface{}

	// returns true if not root
	Valid() bool
}
