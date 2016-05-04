package sets

import (
	"github.com/fncodr/godbase"
	"testing"
)

type testKey int

func (k testKey) Less(other godbase.Key) bool {
	return k < other.(testKey)
}

func genHash(k godbase.Key) uint64 { return uint64(k.(testKey)) }

func TestConstructors(t *testing.T) {
	// Map is mostly meant as a reference for performance comparisons,
	// it only supports enough of the api to run basic tests on top of a native map
	
	NewMap(100)

	// a sorted set
	NewSort()

	// note that the zero value works fine as well
	// var sortSet Sort

	// hashed set with 1000 slots
	NewSortHash(1000, genHash)
}
