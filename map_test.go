package godbase

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type testKey int

func (k testKey) Less(other Cmp) bool {
	return k < other.(testKey)
}

type testItem struct {
	key testKey
}

type testItems []testItem

func randItems(n int) testItems {
	res := make(testItems, n)

	for i := 0; i < n; i++ {
		res[i].key = testKey(i)
	}

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		res[i].key, res[j].key = res[j].key, res[i].key
	}

	return res
}

func PrintTime(start time.Time, msg string, args...interface{}) {
	elapsed := time.Since(start)
	log.Printf("%s:\t%s", fmt.Sprintf(msg, args...), elapsed)
}

func testMapBasics(label string, m Map, its testItems, t *testing.T) {
	start := time.Now()
	for i, it := range its {
		m.Insert(it.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert1", len(its), label)
	
	if l := m.Len(); l != int64(len(its)) {
		t.Errorf("invalid Len() after Insert(): %v / %v", l, len(its))
	}

	start = time.Now()
	for i := 0; i < len(its) / 2; i++ {
		if res := m.Delete(its[i].key, nil); res != 1 {
			t.Errorf("invalid Delete(%v) res: %v", its[i].key, res)
		}
	}
	PrintTime(start, "%v * %v.Delete1", len(its), label)

	start = time.Now()
	for i, it := range its {
		m.Insert(it.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert2", len(its), label)

	start = time.Now()
	for _, it := range its {
		m.Delete(it.key, nil)
	}
	PrintTime(start, "%v * %v.Delete2", len(its), label)

}

func TestMapBasics(t *testing.T) {
	its := randItems(1000000)

	mm := NewMapMap()
	testMapBasics("MapMap", mm, its, t) 

	a := NewSkipNodeAlloc(55)
	sm := NewSkipMap(a, 14)
	testMapBasics("SkipMap", sm, its, t) 

	hm := NewHashMap(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 100, a, 10)
	testMapBasics("HashMap", hm, its, t) 
}
