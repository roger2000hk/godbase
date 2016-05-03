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
	var i int

	for _, it := range its {
		if i = s.Insert(0, godbase.Int64Key(it)); i == -1 {
			b.Errorf("insert failed: %v", it)
		}
	}

	for _, it := range its {
		if i := s.Index(0, godbase.Int64Key(it)); i == -1 {
			b.Errorf("not found: %v", it)
		}
	}

	for _, it := range its {
		if i = s.Delete(0, godbase.Int64Key(it)); i == -1 {
			b.Errorf("delete failed: %v", it)
		}
	}
}

func BenchmarkSortBasics(b *testing.B) {
	const nreps = 20000
 	var s Sort
	s.Resize(nreps)
	runBasicTests(b, &s, randits(nreps))
}

var hashslots = 100000
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

