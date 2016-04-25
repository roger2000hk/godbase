package maps

import (
	//"fmt"
	"testing"
)

const multiReps = 20000
var multiSkipAlloc = NewSkipAlloc(100)
var multiIts1 = sortedItems(multiReps)
var multiIts2 = sortedItems(multiReps)
var multiIts3 = sortedItems(multiReps)

func runMultiTests(t *testing.B, label string, m testAny, its1, its2, its3 []testItem) {
	for i, it := range its1 {
		m.testInsert(nil, it.skipNode.key, &its1[i], true)
		m.testInsert(nil, it.skipNode.key, &its2[i], true)
		m.testInsert(nil, it.skipNode.key, &its3[i], true)
	}

	if l := m.Len(); l != int64(len(its1) * 3) {
		t.Errorf("%v invalid len after insert: %v / %v", label, l, len(its1) * 3)
	}

	for i, it := range its1 {
		if res, cnt := m.testDelete(nil, nil, it.skipNode.key, &its1[i]); cnt != 1 {
			t.Errorf("%v invalid multi delete1 (%v) res: %v/%v", label, it.skipNode.key, res, cnt)
		}
	}

	if l := m.Len(); l != int64(len(its1) * 2) {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, len(its1) * 2)
	}

	for i, it := range its1 {
		if res, ok := m.testFind(nil, it.skipNode.key, nil); !ok {
			t.Errorf("%v invalid find res0: %v", label, res)
		}
		if res, ok := m.testFind(nil, it.skipNode.key, &its1[i]); ok {
			t.Errorf("%v invalid find res1: %v", label, res)
		} 
		if res, ok := m.testFind(nil, it.skipNode.key, &its2[i]); !ok {
			t.Errorf("%v invalid find res2: %v", label, res)
		} 
		if res, ok := m.testFind(nil, it.skipNode.key, &its3[i]); !ok {
			t.Errorf("%v invalid find res3: %v", label, res)
		} 
	}

	for _, it := range its1 {
		if res, cnt := m.testDelete(nil, nil, it.skipNode.key, nil); cnt != 2 {
			t.Errorf("%v invalid multi delete2 (%v) res: %v/%v", 
				label, it.skipNode.key, res, cnt)
		}
	}

	if l := m.Len(); l != 0 {
		t.Errorf("%v invalid len after delete1: %v / %v", label, l, 0)
	}
}

func BenchmarkMultiSkip(t *testing.B) {
	runMultiTests(t, "Skip", NewSkip(nil, 14), multiIts1, multiIts2, multiIts3) 
}

func BenchmarkMultiSkipSlab(t *testing.B) {
	runMultiTests(t, "Skip/Slab", NewSkip(multiSkipAlloc, 14), multiIts1, multiIts2, multiIts3) 
}

func BenchmarkMultiESkip(t *testing.B) {
	runMultiTests(t, "ESkip", NewESkip(), multiIts1, multiIts2, multiIts3) 
}

func BenchmarkMultiSkipHash(t *testing.B) {
	runMultiTests(t, "SkipHash", NewSkipHash(genHash, 80000, multiSkipAlloc, 1), 
		multiIts1, multiIts2, multiIts3)
}

func BenchmarkMultiESkipHash(t *testing.B) {
	runMultiTests(t, "ESkipHash", NewESkipHash(genHash, 50000), multiIts1, multiIts2, multiIts3)
}
