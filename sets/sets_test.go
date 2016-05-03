package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"math/rand"
	"testing"
)

func sortits(n int) []int64 {
	res := make([]int64, n)

	for i := 0; i < n; i++ {
		res[i] = int64(i)
	}

	return res
}

func randits(n int) []int64 {
	its := sortits(n)

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		its[i], its[j] = its[j], its[i]
	}

	return its
}

func runBasicTests(b *testing.B, s godbase.Set, its []int64) {
	for _, it := range its {
		if _, ok := s.Insert(0, godbase.Int64Key(it), false); !ok {
			b.Errorf("insert failed: %v", it)
		}
	}

	for _, it := range its {
		if i := s.First(0, godbase.Int64Key(it)); i == -1 {
			b.Errorf("not found: %v", it)
		}
	}

	for _, it := range its {
		if i := s.Delete(0, godbase.Int64Key(it)); i == -1 {
			b.Errorf("delete failed: %v", it)
		}
	}
}

func BenchmarkSortBasics(b *testing.B) {
	const nreps = 20000
	runBasicTests(b, new(Sort).Resize(nreps), randits(nreps))
}

var hashslots = 400000
var hashits = randits(200000)

func BenchmarkSortHashBasics(b *testing.B) {
	var s SortHash
	s.Init(hashslots, func(k godbase.Key) uint64 { return uint64(k.(godbase.Int64Key)) })
	runBasicTests(b, &s, hashits)
}

func BenchmarkMapHashBasics(b *testing.B) {
	var s MapHash
	s.Init(hashslots, func(k godbase.Key) interface{} { return int64(k.(godbase.Int64Key)) % int64(hashslots) }, 
		func(_ godbase.Key) godbase.Set { return new(Sort) })
	runBasicTests(b, &s, hashits)
}

func BenchmarkMapBasics(b *testing.B) {
	runBasicTests(b, NewMap(0), hashits)
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

	if i := s.First(0, godbase.Int64Key(3)); i != 2 {
		t.Errorf("wrong first res for multi: %v", i)
	}

	if i := s.Last(0, -1, godbase.Int64Key(3)); i != 4 {
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
