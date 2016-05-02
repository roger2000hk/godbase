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
	var ok bool

	for _, it := range its {
		if s, ok = s.Insert(godbase.Int64Key(it)); !ok {
			b.Errorf("insert failed: %v", it)
		}
	}

	for _, it := range its {
		if !s.HasKey(godbase.Int64Key(it)) {
			b.Errorf("not found: %v", it)
		}
	}

	for _, it := range its {
		if s, ok = s.Delete(godbase.Int64Key(it)); !ok {
			b.Errorf("delete failed: %v", it)
		}
	}
}

func BenchmarkSortBasics(b *testing.B) {
 	var s Sort

	runBasicTests(b, godbase.Set(s), randits(5000))
}

func BenchmarkHashBasics(b *testing.B) {
	var s Hash
	s.Init(100000, func(k godbase.Key) uint64 { return uint64(k.(godbase.Int64Key)) })
	runBasicTests(b, s, randits(100000))
}
