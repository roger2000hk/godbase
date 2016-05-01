package godbase

// Iters are circular and cheap, since they are nothing but a common 
// interface on top of actual nodes. 

type Iter interface {
	// Returns key for elem or nil if root
	Key() Key

	// Returns iter to next elem
	Next() Iter

	// Returns val for elem
	Val() interface{}

	// Returns true if not root
	Valid() bool
}
