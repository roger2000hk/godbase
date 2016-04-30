package maps

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"testing"
)

var fullIts1 = sortedItems(testReps)
var fullIts2 = sortedItems(testReps)
var fullIts3 = sortedItems(testReps)

func runFullTests(t *testing.B, label string, m godbase.Map, its1, its2, its3 []testItem) {
	for i, it := range its1 {
		m.Insert(nil, it.node.key, &its1[i].node, true)
		m.Insert(nil, it.node.key, &its2[i].node, true)
		m.Insert(nil, it.node.key, &its3[i].node, true)
	}

	if l := m.Len(); l != int64(len(its1) * 3) {
		t.Errorf("%v invalid len after insert: %v / %v", label, l, len(its1) * 3)
	}

	for i, it := range its1 {
		k := it.node.key
		v := &its1[i].node
		if res, cnt := m.Delete(nil, nil, k, v); cnt != 1 {
			t.Errorf("%v invalid full delete1 (%v) res: %v/%v", label, it.node.key, res, cnt)
		}
	}

	if l := m.Len(); l != int64(len(its1) * 2) {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, len(its1) * 2)
	}

	for i, it := range its1 {
		k := it.node.key

		if res, ok := m.Find(nil, k, nil); !ok {
			t.Errorf("%v invalid find res0: %v", label, res)
		}

		if res, ok := m.Find(nil, k, &its1[i].node); ok {
			t.Errorf("%v invalid find res1: %v", label, res)
		} 

		if res, ok := m.Find(nil, k, &its2[i].node); 
		!ok || res.Val() != &its2[i].node {
			t.Errorf("%v invalid find res2: %v", label, res)
		} 

		if res, ok := m.Find(nil, k, &its3[i].node); 
		!ok  || res.Val() != &its3[i].node {
			t.Errorf("%v invalid find res3: %v", label, res)
		} 
	}

	for _, it := range its1 {
		k := it.node.key
		if res, cnt := m.Delete(nil, nil, k, nil); cnt != 2 {
			t.Errorf("%v invalid full delete2 (%v) res: %v/%v", 
				label, it.node.key, res, cnt)
		}
	}

	if l := m.Len(); l != 0 {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, 0)
	}
}

func BenchmarkFullSort(t *testing.B) {
	runFullTests(t, "Sort", NewSort(testLevels), fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSlab(t *testing.B) {
	runFullTests(t, "Slab", NewSlab(testAlloc, testLevels), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESort(t *testing.B) {
	runFullTests(t, "ESort", NewESort(), fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSlabHash(t *testing.B) {
	runFullTests(t, "SlabHash", 
		NewHash(NewSlabSlots(testSlots, genHash, testAlloc, testHashLevels)),
		fullIts1, fullIts2, fullIts3)
}

func BenchmarkFullESortHash(t *testing.B) {
	runFullTests(t, "ESortHash", NewHash(NewESortSlots(testESlots, genHash)), 
		fullIts1, fullIts2, fullIts3)
}

func BenchmarkFullSortBasicHash(t *testing.B) {
	runFullTests(t, "SortBasicHash", NewHash(NewSlots(testSlots, genHash, allocSlab)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESortBasicHash(t *testing.B) {
	runFullTests(t, "ESortBasicHash", NewHash(NewSlots(testESlots, genHash, allocESort)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullSortMapHash(t *testing.B) {
	runFullTests(t, "SortMapHash", NewHash(NewMapSlots(testSlots, genMapHash, allocSlab)), 
		fullIts1, fullIts2, fullIts3) 
}

func BenchmarkFullESortMapHash(t *testing.B) {
	runFullTests(t, "ESortMapHash", NewHash(NewMapSlots(testESlots, genMapHash, allocESort)), 
		fullIts1, fullIts2, fullIts3)
}
