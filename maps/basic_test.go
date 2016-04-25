package maps

import (
	//"fmt"
	"testing"
)

func TestEmbedded(t *testing.T) {
	m := NewESkip()
	
	const n = 100
	its := make([]testItem, n)

	for i := 0; i < n; i++ {
		k := testKey(i)
		m.Insert(nil, k, &its[i].skipNode, false)

		res, ok := m.Find(nil, k, nil)
		res = res.Next()

		if !ok || res.Key() != k || res.Val().(*ESkipNode) != &its[i].skipNode {
			t.Errorf("invalid find res: %v/%v/%v", i, res.Key(), res.Val())
		} else if toTestItem(res.Val().(*ESkipNode)) != &its[i] {
			t.Errorf("invalid find res: %v/%v", res.Key(), res.Val())
		}
	}
}

func TestConstructors(t *testing.T) {
	// Map is mostly meant as a reference for performance comparisons,
	// it only supports enough of the api to run basic tests on top of 
	// a native map.
	NewMap()
	
	// 10 level skip map with separately allocated nodes
	NewSkip(nil, 10)

	// slab allocator with 50 nodes per slab
	a := NewSkipAlloc(50)

	// 20 level skip map with slab allocated nodes
	NewSkip(a, 20)

	// skip map with embedded nodes
	NewESkip()

	// 2 level hashed skip map with 1000 slots and slab allocated nodes
	NewSkipHash(genHash, 1000, a, 2)

	// hashed skip map with 10000 slots and embedded nodes
	NewESkipHash(genHash, 10000)
}

const basicReps = 50000
var basicIts = randItems(basicReps)
var basicSkipAlloc = NewSkipAlloc(100)

func runBasicTests(t *testing.B, label string, m testAny, its []testItem) {
	for i, it := range its {
		m.testInsert(nil, it.skipNode.key, &its[i], false)
	}

	if l := m.Len(); l != int64(len(its)) {
		t.Errorf("invalid len after insert: %v / %v", l, len(its))
	}

	for i, it := range its {
		k := it.skipNode.key
		v := &its[i]

		res, ok := m.testFind(nil, k, nil)
		if res != nil {
			res = res.Next()
		}

		if !ok || (res != nil && (res.Key() != k || 
			(res.Val() != v && res.Val() != &v.skipNode))) {
			t.Errorf("%v invalid find(%v) res: %v", label, k, res.Key())		
		}

		res, ok = m.testFind(nil, k, v); 
		if res != nil {
			res = res.Next()
		}
		
		if !ok || (res != nil && (res.Key() != k || 
			(res.Val() != v && res.Val() != &v.skipNode))) {
			t.Errorf("%v invalid find(%v) res: %v", label, k, res.Key())		
		}
	}

	for i := 0; i < len(its) / 2; i++ {
		k := its[i].skipNode.key

		if res, cnt := m.testDelete(nil, nil, k, nil); 
		cnt != 1 || (res != nil && res.Key() != nil && !k.Less(res.Key())) {
			t.Errorf("%v invalid delete(%v) res: %v", label, k, res.Key())
		}
	}

	for i, it := range its {
		k := it.skipNode.key
		v := &its[i]

		if res, ok := m.testInsert(nil, it.skipNode.key, v, false); 
		(((i < len(its) / 2) && !ok) || ((i >= len(its) / 2) && ok)) &&
			(res != nil && (res.Key() != k || 
			(res.Val() != v && res.Val() != v.skipNode))) {
			t.Errorf("%v invalid insert(%v) res: %v", label, k, res)		
		}
	}

	for _, it := range its {
		m.testDelete(nil, nil, it.skipNode.key, nil)
	}
}

func BenchmarkBasicMap(t *testing.B) {
	runBasicTests(t, "Map", NewMap(), basicIts) 
}

func BenchmarkBasicSkip(t *testing.B) {
	runBasicTests(t, "Skip", NewSkip(nil, 14), basicIts) 
}

func BenchmarkBasicSkipSlab(t *testing.B) {
	runBasicTests(t, "Skip/Slab", NewSkip(basicSkipAlloc, 14), basicIts) 
}

func BenchmarkBasicESkip(t *testing.B) {
	runBasicTests(t, "ESkip", NewESkip(), basicIts) 
}

func BenchmarkBasicSkipHash(t *testing.B) {
	runBasicTests(t, "SkipHash", NewSkipHash(genHash, 80000, basicSkipAlloc, 1), basicIts) 
}

func BenchmarkBasicESkipHash(t *testing.B) {
	runBasicTests(t, "ESkipHash", NewESkipHash(genHash, 50000), basicIts) 
}