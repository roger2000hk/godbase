package maps

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"testing"
)

func runCutTests(t *testing.T, m godbase.Map) {
	its := sortedItems(100)

	for i, it := range its {
		k := it.node.key
		m.Insert(nil, k, &its[i].node, false)
	}

	start, _ := m.Find(nil, its[90].node.key, nil)
	if k := start.Key(); int(k.(testKey)) != 90 {
		t.Errorf("invalid start: %v", start)
	}

	end, _ := m.Find(nil, its[10].node.key, nil)
	if k := end.Key(); int(k.(testKey)) != 10 {
		t.Errorf("invalid end: %v", end)
	}

	cm := m.Cut(start, end, 
		func (k godbase.Key, v interface{}) (godbase.Key, interface{}) { 
			return testKey(int(k.(testKey)) * 2), v
		})

	if l := cm.Len(); l != 20 {
		t.Errorf("invalid cut target len: %v", l)
	}

	if l := m.Len(); l != 80 {
		t.Errorf("invalid cut src len: %v", l)
	}

}


func TestCut(t *testing.T) {
	runCutTests(t, NewSort(3))
	runCutTests(t, NewESort())
}

func TestEmbedded(t *testing.T) {
	m := NewESort()
	
	const n = 100
	its := make([]testItem, n)

	for i := 0; i < n; i++ {
		k := testKey(i)
		m.Insert(nil, k, &its[i].node, false)

		res, ok := m.Find(nil, k, nil)
		if !ok || res.Key() != k || res.Val().(*ENode) != &its[i].node {
			t.Errorf("invalid find res: %v/%v/%v", i, res.Key(), res.Val())
		} else if toTestItem(res.Val().(*ENode)) != &its[i] {
			t.Errorf("invalid find res: %v/%v", res.Key(), res.Val())
		}
	}
}

func TestConstructors(t *testing.T) {
	// Map is mostly meant as a reference for performance comparisons,
	// it only supports enough of the api to run basic tests on top of 
	// a native map.
	
	NewMap()
	
	// 10 level sorted map
	NewSort(10)

	// slab allocator with 50 nodes per slab
	a := NewSlabAlloc(50)

	// 20 level sorted map with slab allocated nodes
	NewSlab(a, 20)

	// sorted map with embedded nodes
	NewESort()

	// 1000 hash slots backed by 2 level maps with slab allocated nodes
	NewSlabHash(1000, genHash, a, 2)

	// 1000 hash slots backed by 2 level maps with separately allocated nodes
	NewSortHash(1000, genHash, 2)

	// 1000 hash slots backed by maps with embedded nodes
	NewESortHash(1000, genHash)
}

var basicIts = randItems(testReps)

func runBasicTests(t *testing.B, label string, m godbase.Map, its []testItem) {
	for i, it := range its {
		if res, ok := m.Insert(nil, it.node.key, &its[i].node, false); !ok {
			t.Errorf("invalid insert res: %v", res)
		}
	}

	if l := m.Len(); l != int64(len(its)) {
		t.Errorf("invalid len after insert: %v / %v", l, len(its))
	}

	for i, it := range its {
		k := it.node.key
		v := &its[i].node

		res, ok := m.Find(nil, k, nil)

		if !ok || (res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid find(%v) res: %v/%v/%v/%v", label, k, ok, res.Key() == k, res.Val().(*ENode).key, v.key)
		}

		res, ok = m.Find(nil, k, v); 
		
		if !ok || (res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid find(%v) res: %v", label, k, res)		
		}
	}

	for i := 0; i < len(its) / 2; i++ {
		k := its[i].node.key

		res, cnt := m.Delete(nil, nil, k, nil);
		if cnt != 1 || (res != nil && res.Key() != nil && !k.Less(res.Key())) {
			t.Errorf("%v invalid delete(%v) res: %v", label, k, res)
		}
	}

	for i, it := range its {
		k := it.node.key
		v := &its[i].node

		res, ok := m.Insert(nil, it.node.key, v, false)
		if (((i < len(its) / 2) && !ok) || ((i >= len(its) / 2) && ok)) &&
			(res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid insert(%v) res: %v", label, k, res)		
		}
	}

	for _, it := range its {
		if res, cnt := m.Delete(nil, nil, it.node.key, nil); cnt != 1 {
			t.Errorf("invalid delete res: %v", res)
		}
	}
}

func BenchmarkBasicMap(t *testing.B) {
	runBasicTests(t, "Map", NewMap(), basicIts) 
}

func BenchmarkBasicSort(t *testing.B) {
	runBasicTests(t, "Sort", NewSort(testLevels), basicIts) 
}

func BenchmarkBasicSlab(t *testing.B) {
	runBasicTests(t, "Slab", NewSlab(testAlloc, testLevels), basicIts) 
}

func BenchmarkBasicESort(t *testing.B) {
	runBasicTests(t, "ESort", NewESort(), basicIts) 
}

func BenchmarkBasicSlabHash(t *testing.B) {
	runBasicTests(t, "SlabHash", 
		NewSlabHash(testSlots, genHash, testAlloc, testHashLevels), 
		basicIts) 
}

func BenchmarkBasicSortHash(t *testing.B) {
	runBasicTests(t, "SortHash", NewSortHash(testSlots, genHash, testHashLevels), basicIts) 
}

func BenchmarkBasicESortHash(t *testing.B) {
	runBasicTests(t, "ESortHash", NewESortHash(testESlots, genHash), basicIts) 
}
