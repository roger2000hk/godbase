package maps

import (
	"math/rand"
	"unsafe"
)

const (
	// Nr of slots for embedded hash maps
	testESlots = 5000

	// Nr of levels for hashed skip maps
	testHashLevels = 3

	// Nr of levels for non-hashed skip maps
	testLevels = 14

	// Nr of elems / reps
	testReps = 50000

	// Size of skip node slabs
	testSlabSize = 100

	// Nr of slots for non-embedded hash maps
	testSlots = 10000
)

var testItemOffs = unsafe.Offsetof(new(testItem).skipNode)
var testSkipAlloc = NewSkipAlloc(testSlabSize)

type testAny interface {
	Any
	testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int)
	testFind(start Iter, key Cmp, val interface{}) (Iter, bool)
	testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool)
}

type testItem struct {
	skipNode ESkipNode
}

type testItems []testItem
type testKey int

func (k testKey) Less(other Cmp) bool {
	return k < other.(testKey)
}

func genHash(k Cmp) uint64 { return uint64(k.(testKey)) }

func toTestItem(node *ESkipNode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

func sortedItems(n int) testItems {
	res := make(testItems, n)

	for i := 0; i < n; i++ {
		res[i].skipNode.Init(testKey(i))
	}

	return res
}

func reverseItems(n int) testItems {
	res := make(testItems, n)

	for i := n-1; i >= 0; i-- {
		res[i].skipNode.Init(testKey(i))
	}

	return res
}

func randItems(n int) testItems {
	res := sortedItems(n)

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		res[i].skipNode.key, res[j].skipNode.key = res[j].skipNode.key, res[i].skipNode.key
	}

	return res
}

func allocESkip(_ Cmp) Any {
	return NewESkip()
}

func allocSkip(_ Cmp) Any {
	return NewSkip(testSkipAlloc, testHashLevels)
}
