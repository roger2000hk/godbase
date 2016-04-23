package maps

import (
	"bufio"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func (m *ESkip) testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.Insert(key, &val.(*testItem).skipNode, allowMulti)
	return toTestItem(res.(*ESkipNode)), ok
}

func (m *SkipHash) testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.Insert(key, val, allowMulti)
	return res, ok
}

func (m *ESkipHash) testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.Insert(key, &val.(*testItem).skipNode, allowMulti)
	return toTestItem(res.(*ESkipNode)), ok
}

func (m Map) testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.Insert(key, val, allowMulti)
	return res, ok
}

func (m *Skip) testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.Insert(key, val, allowMulti)
	return res, ok
}

func testMapBasics(label string, m testMap, its testItems) {
	start := time.Now()
	for i, it := range its {
		m.testInsert(it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert1", len(its), label)

	if l := m.Len(); l != int64(len(its)) {
		log.Printf("invalid Len() after Insert(): %v / %v", l, len(its))
	}

	start = time.Now()
	for i := 0; i < len(its) / 2; i++ {
		if res := m.Delete(its[i].skipNode.key, nil); res != 1 {
			log.Printf("invalid Delete(%v) res: %v", its[i].skipNode.key, res)
		}
	}
	PrintTime(start, "%v * %v.Delete1", len(its), label)

	start = time.Now()
	for i, it := range its {
		//fmt.Printf("%v\n", m)
		m.testInsert(it.skipNode.key, &its[i], false)
	}
	PrintTime(start, "%v * %v.Insert2", len(its), label)


	start = time.Now()
	for _, it := range its {
		m.Delete(it.skipNode.key, nil)
	}
	PrintTime(start, "%v * %v.Delete2", len(its), label)
}

func RunBasicTests() {
	its := randItems(100000)

	mm := NewMap()
	testMapBasics("Map", mm, its) 

	a := NewSkipNodeAlloc(55)
	//ssm := NewSkip(a, 1)
	//testMapBasics("List", ssm, its) 

	file, err := os.Create("test.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(bufio.NewWriter(file))
	
	sm := NewSkip(a, 14)
	testMapBasics("Skip", sm, its) 

	esm := NewESkip()
	testMapBasics("ESkip", esm, its) 

	hm := NewSkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 80000, a, 1)
	testMapBasics("SkipHash", hm, its)

	ehm := NewESkipHash(func(k Cmp) uint64 { return uint64(k.(testKey)) }, 50000)
	testMapBasics("ESkipHash", ehm, its)

	pprof.StopCPUProfile()
	file.Close()
}