package maps

import (
	"github.com/fncodr/godbase"
	"math/rand"
	"unsafe"
)

const (
	// Nr of slots for embedded hash maps
	testESlots = 20000

	// Nr of levels for hashed maps
	testHashLevels = 1

	// Nr of levels for non-hashed maps
	testLevels = 14

	// Nr of elems / reps
	testReps = 20000

	// Size of slabs
	testSlabSize = 500

	// Nr of slots for non-embedded hash maps
	testSlots = 30000
)

var testItemOffs = unsafe.Offsetof(new(testItem).node)
var testAlloc = NewSlabAlloc(testSlabSize)

type testItem struct {
	node ENode
}

type testItems []testItem
type testKey int

func (k testKey) Less(other godbase.Key) bool {
	return k < other.(testKey)
}

func genHash(k godbase.Key) uint64 { return uint64(k.(testKey)) }

func genMapHash(k godbase.Key) interface{} { return k }

func toTestItem(node *ENode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

func sortedItems(n int) testItems {
	res := make(testItems, n)

	for i := 0; i < n; i++ {
		res[i].node.Init(testKey(i))
	}

	return res
}

func reverseItems(n int) testItems {
	res := make(testItems, n)

	for i := n-1; i >= 0; i-- {
		res[i].node.Init(testKey(i))
	}

	return res
}

func randItems(n int) testItems {
	res := sortedItems(n)

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		res[i].node.key, res[j].node.key = res[j].node.key, res[i].node.key
	}

	return res
}

func allocESort(_ godbase.Key) godbase.Map {
	return NewESort()
}

func allocSlab(_ godbase.Key) godbase.Map {
	return NewSlab(testAlloc, testHashLevels)
}
