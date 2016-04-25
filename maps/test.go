package maps

import (
	"math/rand"
	"unsafe"
)

type testAny interface {
	Any
	testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int)
	testFind(start Iter, key Cmp, val interface{}) (Iter, bool)
	testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool)
}

type testKey int

func (k testKey) Less(other Cmp) bool {
	return k < other.(testKey)
}

func genHash(k Cmp) uint64 { return uint64(k.(testKey)) }

type testItem struct {
	skipNode ESkipNode
}

var testItemOffs = unsafe.Offsetof(new(testItem).skipNode)

func toTestItem(node *ESkipNode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

type testItems []testItem

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
