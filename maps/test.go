package maps

import (
	"bufio"
	"fmt"
	"log"
	"time"
	"math/rand"
	"os"
	"runtime/pprof"
	"unsafe"
)

type testKey int

func (k testKey) Less(other Cmp) bool {
	return k < other.(testKey)
}

type testMap interface {
	Delete(key Cmp, val interface{}) int
	testInsert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool)
	Len() int64
}

type testItem struct {
	skipNode ESkipNode
}

var testItemOffs = unsafe.Offsetof(new(testItem).skipNode)

func toTestItem(node *ESkipNode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

type testItems []testItem

func sortedItems(n int) testItems {
	res := make(testItems, n)

	for i := 0; i < n; i++ {
		res[i].skipNode.Init(testKey(i))
	}

	return res
}

func reverseItems(n int) testItems {
	res := make(testItems, n)

	for i := n-1; i >= 0; i-- {
		res[i].skipNode.Init(testKey(i))
	}

	return res
}

func randItems(n int) testItems {
	res := sortedItems(n)

	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		res[i].skipNode.key, res[j].skipNode.key = res[j].skipNode.key, res[i].skipNode.key
	}

	return res
}

func PrintTime(start time.Time, msg string, args...interface{}) {
	elapsed := time.Since(start)
	log.Printf("%s:\t%s", fmt.Sprintf(msg, args...), elapsed)
}

func RunTests() {
	file, err := os.Create("test.prof")
	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(bufio.NewWriter(file))
	RunBasicTests()
	RunMultiTests()
	RunSortedTests()
	pprof.StopCPUProfile()
}
