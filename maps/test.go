package maps

import (
	"math/rand"
	"unsafe"
)

const (
	// Nr of slots for embedded hash maps
	testESlots = 20000

	// Nr of levels for hashed skip maps
	testHashLevels = 1

	// Nr of levels for non-hashed skip maps
	testLevels = 14

	// Nr of elems / reps
	testReps = 20000

	// Size of skip node slabs
	testSlabSize = 500

	// Nr of slots for non-embedded hash maps
	testSlots = 30000
)

var testItemOffs = unsafe.Offsetof(new(testItem).skipNode)
var testSkipAlloc = NewSkipAlloc(testSlabSize)

type testAny interface {
	Any
	testDelete(start, end Iter, key Key, val interface{}) (Iter, int)
	testFind(start Iter, key Key, val interface{}) (Iter, bool)
	testInsert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool)
}

type testItem struct {
	skipNode ESkipNode
}

type testItems []testItem
type testKey int

func (k testKey) Less(other Key) bool {
	return k < other.(testKey)
}

func genHash(k Key) uint64 { return uint64(k.(testKey)) }

func genMapHash(k Key) interface{} { return k }

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

func allocESkip(_ Key) Any {
	return NewESkip()
}

func allocSkip(_ Key) Any {
	return NewSkip(testSkipAlloc, testHashLevels)
}
