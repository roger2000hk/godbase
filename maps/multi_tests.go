package maps

import (
	//"fmt"
	"log"
	"time"
)

func runMultiTests(label string, m testAny, its1, its2, its3 testItems) {
	start := time.Now()
	for i, it := range its1 {
		m.testInsert(nil, it.skipNode.key, &its1[i], true)
		m.testInsert(nil, it.skipNode.key, &its2[i], true)
		m.testInsert(nil, it.skipNode.key, &its3[i], true)
	}
	PrintTime(start, "%v * %v multi insert", len(its1), label)

	if l := m.Len(); l != int64(len(its1) * 3) {
		log.Panicf("%v invalid len after insert: %v / %v", 
			label, l, len(its1) * 3)
	}

	start = time.Now()
	for i, it := range its1 {
		if res, cnt := m.testDelete(nil, nil, it.skipNode.key, &its1[i]); cnt != 1 {
			log.Panicf("%v invalid multi delete (%v) res: %v/%v", label, it.skipNode.key, res, cnt)
		}
	}
	PrintTime(start, "%v * %v multi delete1", len(its1), label)

	if l := m.Len(); l != int64(len(its1) * 2) {
		log.Panicf("%v invalid len after delete1: %v / %v", 
			label, l, len(its1) * 2)
	}

	start = time.Now()
	for i, it := range its1 {
		if res, ok := m.testFind(nil, it.skipNode.key, &its1[i]); ok {
			log.Panicf("%v invalid find res1: %v", label, res)
		} 
		if res, ok := m.testFind(nil, it.skipNode.key, &its2[i]); !ok {
			log.Panicf("%v invalid find res2: %v", label, res)
		} 
		if res, ok := m.testFind(nil, it.skipNode.key, &its3[i]); !ok {
			log.Panicf("%v invalid find res3: %v", label, res)
		} 
	}
	PrintTime(start, "%v * %v multi find", len(its1), label)
}

func RunMultiTests() {
	const nreps = 10000
	its1 := sortedItems(nreps)
	its2 := sortedItems(nreps)
	its3 := sortedItems(nreps)

	a := NewSkipNodeAlloc(55)
	//ssm := NewSkip(a, 1)
	//runMultiTests("List", ssm, its1, its2, its3) 

	sm := NewSkip(a, 14)
	runMultiTests("Skip", sm, its1, its2, its3) 

	esm := NewESkip()
	runMultiTests("ESkip", esm, its1, its2, its3) 

	hm := NewSkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 80000, a, 1)
	runMultiTests("SkipHash", hm, its1, its2, its3)

	ehm := NewESkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 50000)
	runMultiTests("ESkipHash", ehm, its1, its2, its3)
}
