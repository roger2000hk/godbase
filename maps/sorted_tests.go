package maps

import (
	"log"
	"time"
)

type sortedTestMap interface {
	testMap
}

func runSortedTests(label string, m sortedTestMap, its testItems) {
	start := time.Now()
	for i, it := range its {
		m.testInsert(it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert", len(its), label)

	if l := m.Len(); l != int64(len(its)) {
		log.Printf("invalid Len() after Insert(): %v / %v", l, len(its))
	}
}

func RunSortedTests() {
	its := reverseItems(100000)

	mm := NewMap()
	runSortedTests("Map", mm, its) 

	a := NewSkipNodeAlloc(55)
	//ssm := NewSkip(a, 1)
	//runSortedTests("List", ssm, its) 

	sm := NewSkip(a, 14)
	runSortedTests("Skip", sm, its) 

	esm := NewESkip()
	runSortedTests("ESkip", esm, its) 

	hm := NewSkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 80000, a, 1)
	runSortedTests("SkipHash", hm, its)

	ehm := NewESkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 50000)
	runSortedTests("ESkipHash", ehm, its)
}
