package maps

import (
	"fmt"
	"log"
	"time"
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

func PrintTime(start time.Time, msg string, args...interface{}) {
	elapsed := time.Since(start)
	log.Printf("%s:\t%s", fmt.Sprintf(msg, args...), elapsed)
}

func (m *ESkip) testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}

	res, cnt := m.Delete(start, end, key, val)
	return res, cnt
}

func (m *SkipHash) testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	res, cnt := m.Delete(start, end, key, val)
	return res, cnt
}

func (m *ESkipHash) testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}
	
	res, cnt := m.Delete(start, end, key, val)
	return res, cnt
}

func (m Map) testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	res, cnt := m.Delete(start, end, key, val)
	return res, cnt
}

func (m *Skip) testDelete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	res, cnt := m.Delete(start, end, key, val)
	return res, cnt
}

func (m *ESkip) testFind(start Iter, key Cmp, val interface{}) (Iter, bool) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}

	res, ok := m.Find(start, key, val)
	return res, ok
}

func (m *SkipHash) testFind(start Iter, key Cmp, val interface{}) (Iter, bool) {
	res, ok := m.Find(start, key, val)
	return res, ok
}

func (m *ESkipHash) testFind(start Iter, key Cmp, val interface{}) (Iter, bool) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}
	
	res, ok := m.Find(start, key, val)
	return res, ok
}

func (m Map) testFind(start Iter, key Cmp, val interface{}) (Iter, bool) {
	res, ok := m.Find(start, key, val)
	return res, ok
}

func (m *Skip) testFind(start Iter, key Cmp, val interface{}) (Iter, bool) {
	res, ok := m.Find(start, key, val)
	return res, ok
}

func (m *ESkip) testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}

	res, ok := m.Insert(start, key, val, allowMulti)
	return res, ok
}

func (m *SkipHash) testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	res, ok := m.Insert(start, key, val, allowMulti)
	return res, ok
}

func (m *ESkipHash) testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	if val != nil {
		val = &val.(*testItem).skipNode
	}

	res, ok := m.Insert(start, key, val, allowMulti)
	return res, ok
}

func (m Map) testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	res, ok := m.Insert(start, key, val, allowMulti)
	return res, ok
}

func (m *Skip) testInsert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	res, ok := m.Insert(start, key, val, allowMulti)
	return res, ok
}
