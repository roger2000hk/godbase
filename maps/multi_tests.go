package maps

import (
	"log"
	"time"
)

func runMultiTests(label string, m testMap, its1, its2 testItems) {
	start := time.Now()
	for i, it := range its1 {
		m.testInsert(it.skipNode.key, &its1[i], true)
		m.testInsert(it.skipNode.key, &its2[i], true)
	}
	PrintTime(start, "%v * %v multi insert", len(its1), label)

	if l := m.Len(); l != int64(len(its1) * 2) {
		log.Printf("%v invalid len after insert: %v / %v", 
			label, l, len(its1) * 2)
	}
}

func RunMultiTests() {
	its1 := sortedItems(100000)
	its2 := sortedItems(100000)

	a := NewSkipNodeAlloc(55)
	//ssm := NewSkip(a, 1)
	//runMultiTests("List", ssm, its1, its2) 

	sm := NewSkip(a, 14)
	runMultiTests("Skip", sm, its1, its2) 

	esm := NewESkip()
	runMultiTests("ESkip", esm, its1, its2) 

	hm := NewSkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 80000, a, 1)
	runMultiTests("SkipHash", hm, its1, its2)

	ehm := NewESkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 50000)
	runMultiTests("ESkipHash", ehm, its1, its2)
}
