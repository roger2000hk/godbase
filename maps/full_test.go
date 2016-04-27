package maps

import (
	//"fmt"
	"testing"
)

var fullIts1 = sortedItems(testReps)
var fullIts2 = sortedItems(testReps)
var fullIts3 = sortedItems(testReps)

func runFullTests(t *testing.B, label string, m Any, its1, its2, its3 []testItem) {
	for i, it := range its1 {
		m.Insert(nil, it.skipNode.key, &its1[i].skipNode, true)
		m.Insert(nil, it.skipNode.key, &its2[i].skipNode, true)
		m.Insert(nil, it.skipNode.key, &its3[i].skipNode, true)
	}

	if l := m.Len(); l != int64(len(its1) * 3) {
		t.Errorf("%v invalid len after insert: %v / %v", label, l, len(its1) * 3)
	}

	for i, it := range its1 {
		k := it.skipNode.key
		v := &its1[i].skipNode
		if res, cnt := m.Delete(nil, nil, k, v); cnt != 1 {
			t.Errorf("%v invalid full delete1 (%v) res: %v/%v", label, it.skipNode.key, res, cnt)
		}
	}

	if l := m.Len(); l != int64(len(its1) * 2) {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, len(its1) * 2)
	}

	for i, it := range its1 {
		k := it.skipNode.key

		if res, ok := m.Find(nil, k, nil); !ok {
			t.Errorf("%v invalid find res0: %v", label, res)
		}

		if res, ok := m.Find(nil, k, &its1[i].skipNode); ok {
			t.Errorf("%v invalid find res1: %v", label, res)
		} 

		if res, ok := m.Find(nil, k, &its2[i].skipNode); 
		!ok || res.Val() != &its2[i].skipNode {
			t.Errorf("%v invalid find res2: %v", label, res)
		} 

		if res, ok := m.Find(nil, k, &its3[i].skipNode); 
		!ok  || res.Val() != &its3[i].skipNode {
			t.Errorf("%v invalid find res3: %v", label, res)
		} 
	}

	for _, it := range its1 {
		k := it.skipNode.key
		if res, cnt := m.Delete(nil, nil, k, nil); cnt != 2 {
			t.Errorf("%v invalid full delete2 (%v) res: %v/%v", 
				label, it.skipNode.key, res, cnt)
		}
	}

	if l := m.Len(); l != 0 {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, 0)
	}
}

func BenchmarkFullSkip(t *testing.B) {
	runFullTests(t, "Skip", NewSkip(nil, testLevels), fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSkipSlab(t *testing.B) {
	runFullTests(t, "Skip/Slab", NewSkip(testSkipAlloc, testLevels), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESkip(t *testing.B) {
	runFullTests(t, "ESkip", NewESkip(), fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSkipHash(t *testing.B) {
	runFullTests(t, "SkipHash", 
		NewHash(NewSkipSlots(testSlots, genHash, testSkipAlloc, testHashLevels)),
		fullIts1, fullIts2, fullIts3)
}

func BenchmarkFullESkipHash(t *testing.B) {
	runFullTests(t, "ESkipHash", NewHash(NewESkipSlots(testESlots, genHash)), 
		fullIts1, fullIts2, fullIts3)
}

func BenchmarkFullSkipAnyHash(t *testing.B) {
	runFullTests(t, "SkipAnyHash", NewHash(NewSlots(testSlots, genHash, allocSkip)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESkipAnyHash(t *testing.B) {
	runFullTests(t, "ESkipAnyHash", NewHash(NewSlots(testESlots, genHash, allocESkip)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSkipMapHash(t *testing.B) {
	runFullTests(t, "SkipMapHash", NewHash(NewMapSlots(testSlots, genMapHash, allocSkip)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESkipMapHash(t *testing.B) {
	runFullTests(t, "ESkipMapHash", NewHash(NewMapSlots(testESlots, genMapHash, allocESkip)), 
		fullIts1, fullIts2, fullIts3)
}
