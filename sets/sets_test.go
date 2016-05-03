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
 	var s Sort
	runBasicTests(b, &s, randits(10000))
}

func BenchmarkSortHashBasics(b *testing.B) {
	var s SortHash
	s.Init(10000, func(k godbase.Key) uint64 { return uint64(k.(godbase.Int64Key)) })
	runBasicTests(b, &s, randits(200000))
}

func BenchmarkMapHashBasics(b *testing.B) {
	var s MapHash
	s.Init(10000, func(k godbase.Key) interface{} { return k.(godbase.Int64Key) % 10000 }, 
		func(_ godbase.Key) godbase.Set { return new(Sort) })
	runBasicTests(b, &s, randits(300000))
}

func BenchmarkMapBasics(b *testing.B) {
	its := randits(300000)
	m := make(map[int64]bool)

	for _, it := range its {
		m[it] = true
	}

	for _, it := range its {
		if _, ok := m[it]; !ok {
			b.Errorf("not found: %v", it)
		}
	}

	for _, it := range its {
		delete(m, it)
	}
}

