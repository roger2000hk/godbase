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

	// 1000 slots backed by generic map allocator
	// generic slots could be used in any of the following examples,
	// but specializing the slot type allows allocating all slots at once and
	// accessing by value which makes a difference in some scenarios.
	// the allocator receives the key as param which enables choosing
	// differend kinds of slot chains for different keys.

	gs := NewSlots(1000, genHash, func (_ Cmp) Any { return NewSkip(nil, 2) })
	NewHash(gs)

	// 1000 skip slots backed by 2 level skip maps with slab allocated nodes
	ss := NewSkipSlots(1000, genHash, a, 2)
	NewHash(ss)

	// 1000 hash slots backed by embedded skip maps
	ess := NewESkipSlots(1000, genHash)
	NewHash(ess)

	// 1000 hash slots backed by hash maps with 100 embedded skip slots
	hs := NewHashSlots(1000, genHash, func (_ Cmp) Slots { return NewESkipSlots(100, genHash) })
	NewHash(hs)
}

var basicIts = randItems(testReps)

func runBasicTests(t *testing.B, label string, m Any, its []testItem) {
	for i, it := range its {
		m.Insert(nil, it.skipNode.key, &its[i].skipNode, false)
	}

	if l := m.Len(); l != int64(len(its)) {
		t.Errorf("invalid len after insert: %v / %v", l, len(its))
	}

	for i, it := range its {
		k := it.skipNode.key
		v := &its[i].skipNode

		res, ok := m.Find(nil, k, nil)
		if res != nil {
			res = res.Next()
		}

		if !ok || (res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid find(%v) res: %v/%v/%v/%v", label, k, ok, res.Key() == k, res.Val().(*ESkipNode).key, v.key)
		}

		res, ok = m.Find(nil, k, v); 
		if res != nil {
			res = res.Next()
		}
		
		if !ok || (res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid find(%v) res: %v", label, k, res)		
		}
	}

	for i := 0; i < len(its) / 2; i++ {
		k := its[i].skipNode.key

		if res, cnt := m.Delete(nil, nil, k, nil); 
		cnt != 1 || (res != nil && res.Key() != nil && !k.Less(res.Key())) {
			t.Errorf("%v invalid delete(%v) res: %v", label, k, res)
		}
	}

	for i, it := range its {
		k := it.skipNode.key
		v := &its[i].skipNode

		if res, ok := m.Insert(nil, it.skipNode.key, v, false); 
		(((i < len(its) / 2) && !ok) || ((i >= len(its) / 2) && ok)) &&
			(res != nil && (res.Key() != k || res.Val() != v)) {
			t.Errorf("%v invalid insert(%v) res: %v", label, k, res)		
		}
	}

	for _, it := range its {
		m.Delete(nil, nil, it.skipNode.key, nil)
	}
}

func BenchmarkBasicMap(t *testing.B) {
	runBasicTests(t, "Map", NewMap(), basicIts) 
}

func BenchmarkBasicSkip(t *testing.B) {
	runBasicTests(t, "Skip", NewSkip(nil, testLevels), basicIts) 
}

func BenchmarkBasicSkipSlab(t *testing.B) {
	runBasicTests(t, "Skip/Slab", NewSkip(testSkipAlloc, testLevels), basicIts) 
}

func BenchmarkBasicESkip(t *testing.B) {
	runBasicTests(t, "ESkip", NewESkip(), basicIts) 
}

func BenchmarkBasicSkipHash(t *testing.B) {
	runBasicTests(t, "SkipHash", 
		NewHash(NewSkipSlots(testSlots, genHash, testSkipAlloc, testHashLevels)), 
		basicIts) 
}

func BenchmarkBasicESkipHash(t *testing.B) {
	runBasicTests(t, "ESkipHash", NewHash(NewESkipSlots(testESlots, genHash)), basicIts) 
}

func BenchmarkBasicSkipAnyHash(t *testing.B) {
	runBasicTests(t, "SkipAnyHash", NewHash(NewSlots(testSlots, genHash, allocSkip)), 
		basicIts) 
}

func BenchmarkBasicESkipAnyHash(t *testing.B) {
	runBasicTests(t, "ESkipAnyHash", NewHash(NewSlots(testESlots, genHash, allocESkip)), 
		basicIts) 
}
