package maps

import (
	"log"
	"time"
)

func runSortedTests(label string, m testAny, its testItems) {
	start := time.Now()
	for i, it := range its {
		m.testInsert(nil, it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v sorted insert", len(its), label)

	if l := m.Len(); l != int64(len(its)) {
		log.Panicf("invalid len after insert: %v / %v", l, len(its))
	}
}

func RunSortedTests() {
	its := reverseItems(100000)
	a := NewSkipNodeAlloc(55)

	ssm := NewSkip(a, 1)
	runSortedTests("List", ssm, its) 

	sm := NewSkip(a, 14)
	runSortedTests("Skip", sm, its) 

	esm := NewESkip()
	runSortedTests("ESkip", esm, its) 
}
