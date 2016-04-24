package maps

import (
	//"fmt"
	"log"
	"time"
)

func runBasicTests(label string, m testAny, its testItems) {
	start := time.Now()
	for i, it := range its {
		m.testInsert(nil, it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert1", len(its), label)

	if l := m.Len(); l != int64(len(its)) {
		log.Panicf("invalid Len() after Insert(): %v / %v", l, len(its))
	}

	start = time.Now()
	for i := 0; i < len(its) / 2; i++ {
		//fmt.Printf("%v\n", m)
		k := its[i].skipNode.key
		if res, cnt := m.testDelete(nil, nil, k, nil); cnt != 1 {
			log.Panicf("%v invalid Delete (%v) res: %v", label, its[i].skipNode.key, res)
		}
	}
	PrintTime(start, "%v * %v.Delete1", len(its), label)

	start = time.Now()
	for i, it := range its {
		m.testInsert(nil, it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert2", len(its), label)


	start = time.Now()
	for _, it := range its {
		m.testDelete(nil, nil, it.skipNode.key, nil)
	}
	PrintTime(start, "%v * %v.Delete2", len(its), label)
}

func RunBasicTests() {
	its := randItems(10000)

	mm := NewMap()
	runBasicTests("Map", mm, its) 

	a := NewSkipNodeAlloc(55)
	//ssm := NewSkip(a, 1)
	//runBasicTests("List", ssm, its) 

	sm := NewSkip(a, 14)
	runBasicTests("Skip", sm, its) 

	esm := NewESkip()
	runBasicTests("ESkip", esm, its) 

	hm := NewSkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 80000, a, 1)
	runBasicTests("SkipHash", hm, its)

	ehm := NewESkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 50000)
	runBasicTests("ESkipHash", ehm, its)
}