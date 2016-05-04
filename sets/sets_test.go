package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"math/rand"
	"testing"
)

func sortelems(n int) []int64 {
	res := make([]int64, n)

	for i := 0; i < n; i++ {
		res[i] = int64(i)
	}

	return res
}

func sortkeys(n int) []godbase.Key {
	res := make([]godbase.Key, n)

	for i := 0; i < n; i++ {
		res[i] = godbase.Int64Key(i)
	}

	return res
}

func randelems(n int) []int64 {
	elems := sortelems(n)

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		elems[i], elems[j] = elems[j], elems[i]
	}

	return elems
}

func runBasicTests(b *testing.B, s godbase.Set, elems []int64) {
	for _, e := range elems {
		if _, ok := s.Insert(0, godbase.Int64Key(e), false); !ok {
			b.Errorf("insert failed: %v", e)
		}
	}

	for _, e := range elems {
		if _, ok := s.First(0, godbase.Int64Key(e)); !ok {
			b.Errorf("not found: %v", e)
		}
	}

	for _, e := range elems {
		if i := s.Delete(0, godbase.Int64Key(e)); i == -1 {
			b.Errorf("delete failed: %v", e)
		}
	}
}

func BenchmarkSortBasics(b *testing.B) {
	const nreps = 20000
	runBasicTests(b, new(Sort).Resize(nreps), randelems(nreps))
}

var hashslots = 400000
var hashelems = randelems(200000)

func BenchmarkSortHashBasics(b *testing.B) {
	var s SortHash
	s.Init(hashslots, func(k godbase.Key) uint64 { return uint64(k.(godbase.Int64Key)) })
	runBasicTests(b, &s, hashelems)
}

func BenchmarkMapBasics(b *testing.B) {
	runBasicTests(b, NewMap(len(hashelems)), hashelems)
}

func runCloneTests(b *testing.B, s godbase.Set, elems []int64) {
	for _, e := range elems {
		if _, ok := s.Insert(0, godbase.Int64Key(e), false); !ok {
			b.Errorf("insert failed: %v", e)
		}
	}
	
	for i := 0; i < 5000; i++ {
		s = s.Clone()
	}
}

var cloneelems = sortelems(1000)

func BenchmarkSortClone(b *testing.B) {
	runCloneTests(b, new(Sort).Resize(len(cloneelems)), cloneelems)
}

func BenchmarkMapClone(b *testing.B) {
	runCloneTests(b, NewMap(len(cloneelems)), cloneelems)
}

func runLoadTests(b *testing.B, s godbase.Set, elems []godbase.Key) {
	s.Load(0, elems...)
}

var loadelems = sortkeys(1000000)

func BenchmarkSortLoad(b *testing.B) {
	runLoadTests(b, new(Sort), loadelems)
}

func BenchmarkMapLoad(b *testing.B) {
	runLoadTests(b, NewMap(len(loadelems)), loadelems)
}

func TestMulti(t *testing.T) {
	var s Sort

	s.Insert(0, godbase.Int64Key(1), false)
	s.Insert(0, godbase.Int64Key(2), false)
	s.Insert(0, godbase.Int64Key(3), false)
	s.Insert(0, godbase.Int64Key(3), true)
	s.Insert(0, godbase.Int64Key(3), true)
	s.Insert(0, godbase.Int64Key(4), false)
	s.Insert(0, godbase.Int64Key(5), false)

	if l := s.Len(); l != 7 {
		t.Errorf("wrong len after multi insert: %v", l)
	}

	if i, ok := s.First(0, godbase.Int64Key(3)); !ok || i != 2 {
		t.Errorf("wrong first res for multi: %v", i)
	}

	if i, ok := s.Last(0, 0, godbase.Int64Key(3)); !ok || i != 4 {
		t.Errorf("wrong last res for multi: %v", i)
	}

	i, ok := s.DeleteAll(0, 4, godbase.Int64Key(3))

	if i != 2 {
		t.Errorf("wrong res from multi delete: %v", ok)
	}

	if ok != 2 {
		t.Errorf("wrong res from multi delete: %v", ok)
	}

	if l := s.Len(); l != 5 {
		t.Errorf("wrong len after multi delete: %v", l)
	}
}
