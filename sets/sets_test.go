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

var its = randits(5000)

func BenchmarkSortBasics(b *testing.B) {
	var ok bool
 	var _s Sort
	s := godbase.Set(_s)

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

