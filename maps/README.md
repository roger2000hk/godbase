# beating The Map
### a glorious quest to write a faster, more flexible map

## why?
Good question. This rabbit hole was a bit deeper than anticipated, but I still feel it was worth the effort given the results. To cut a long story short, I was itching for more flexible maps to implement in-memory indexing in godbase. I quickly found that my lofty goal of adding multi-capability and sorting, while matching the performance of native maps; was far from a walk in the park. Which is probably part of the reason I couldn't find anyone else trying.

## how?
I ended up with two designs based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes; the other with embedded nodes, constant number of levels, using one node per value. These do pretty good by themselves, especially considering that they add sorting to the mix; both currently hovering around 2-3 times slower than a native map. The embedded flavor usually wins the allocation race by a slim margin but pays the price of having a fixed number of levels for tiny / huge datasets. Still, somewhere along a line; not separately allocating nodes affects overall performance positively.

Once sorted maps were working properly, I had the crazy idea to put a hash on top just to see what happens. It turns out that dividing the dataset into a tuned number of ordered sets helps puts us consistently ahead of native maps in the synthetic performance game for millions of items. That's far from the end of the story though. I still haven't had enough time to ponder the consequences of having access to both hashed and ordered aspects of the data simultaneously, but I have a hunch it will bend the rules to my advantage in a number of tricky scenarios. Additionally; the hash adapter supports any kind of map for slot chains, which opens the door for multi level hashing where each chain is another hash that further divides the dataset.

## status
Basic functionality and testing in place; bells, whistles & polish are still on the stove.

## benchmarks
go test -bench=.*

## license
NOP

## code
I trust you'll find godbase more RISC/Lispy than your everyday set/map api. Providing an optimal api is part of implementing an optimal algorithm, and there's more low hanging fruit in the garden of set/map apis than most places. It's obvious to me that academic dogmatics and software (or life in general, for that matter) isn't really the match made in heaven it's being sold as.

### interfaces

```go

// All map keys are requred to support Cmp
type Cmp interface {
	Less(Cmp) bool
}

// All hash maps require a hash fn
type HashFn func (Cmp) uint64

// Iters are circular and cheap, since they are nothing but a common 
// interface on top of actual nodes. They are positioned before start
// on return, so you need to call Next() to get the first elem.

type Iter interface {
	// Returns true if next elem is not root
	HasNext() bool

	// Returns true if prev elem is not root
	HasPrev() bool

	// Returns key for elem or nil if root
	Key() Cmp

	// Returns iter to next elem
	Next() Iter

	// Returns iter to prev elem
	Prev() Iter

	// Returns val for elem or nil if root
	Val() interface{}
}

// Basic map ops supported by all implementations
type Any interface {
	// Cuts elems from start to end for which fn returns true into new set;
	// start, end & fn are all optional. Circular cuts, with start/end on
	// opposite sides of root; are supported. Returns a cut from the start slot
	// for hash maps.

	Cut(start, end Iter, fn TestFn) Any

	// Deletes all elems after start matching key/val;
	// start, end, key & val are all optional, nil deletes all elems. Specifying 
	// iters for hash maps only works within the same slot. Circular deletes,
	// with start/end on opposite sides of root; are supported. Returns an iter to next 
	// elem and the number of deleted elems.

	Delete(start, end Iter, key Cmp, val interface{}) (Iter, int)

	// Returns iter for first elem after start matching key and ok;
	// start & val are optional, specifying a start iter for hash maps only works within the 
	// same slot.

	Find(start Iter, key Cmp, val interface{}) (Iter, bool)
	
	// Inserts key/val into map after start;
	// start & val are both optional, dup checks can be disabled by setting allowMulti to false. 
	// Returns iter to inserted val & true on success, or iter to existing val & false on dup. 
	// Specifying a start iter for hash maps only works within the same slot.

	Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool)

	// Returns the number of elems in map
	Len() int64
}

type TestFn func (Cmp, interface{}) bool


```

### constructors

```go

type testKey int

func (k testKey) Less(other Cmp) bool {
	return k < other.(testKey)
}

func genHash(k Cmp) uint64 { return uint64(k.(testKey)) }

func TestConstructors(t *testing.T) {
	// Map is mostly meant as a reference for performance comparisons,
	// it only supports enough of the api to run basic tests on top of 
	// a native map.
	NewMap()
	
	// 10 level skip map with separately allocated nodes
	NewSkip(nil, 10)

	// slab allocator with 50 nodes per slab
	a := NewSkipAlloc(50)

	// 20 level skip map with slab allocated nodes
	NewSkip(a, 20)

	// skip map with embedded nodes
	NewESkip()

	// 1000 hash slots backed by 2 level skip maps with slab allocated nodes
	ss := NewSkipSlots(1000, genHash, a, 2)

	// hash map based on skip slots
	NewHash(ss)

	// 500 hash slots backed by an embedded skip map
	ess := NewESkipSlots(500, genHash)

	// hash map based on embedded skip slots
	NewHash(ess)
}

```

### embedded nodes
I picked up the idea of embedding node infrastructure into elems from the Linux kernel, but I'm sure the idea is at least as old as the C language. It's a nice tool to reduce memory allocation which bends the rules enough for the previously undoable to become possible. If you don't mind keeping a reference per collection in your type, or sprinkling a bit of unsafe magic on top; this might be for you. lists.EDouble contains a double-linked list implementation based on the same idea.

```go

// pretend this is your value type
type testItem struct {
	skipNode ESkipNode

	// additional fields...
}

// calculate offset of node within struct
var testItemOffs = unsafe.Offsetof(new(testItem).skipNode)

// helper to get pointer to item from node
func toTestItem(node *ESkipNode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

func TestEmbedded(t *testing.T) {
	m := NewESkip()
	
	const n = 100
	its := make([]testItem, n)

	for i := 0; i < n; i++ {
		k := testKey(i)

		// the map only deals with nodes,
		// translation to/from values is left to client code
		m.Insert(nil, k, &its[i].skipNode, false)

		res, ok := m.Find(nil, k, nil)
		res = res.Next()

		if !ok || res.Key() != k || res.Val().(*ESkipNode) != &its[i].skipNode {
			t.Errorf("invalid iter: %v/%v/%v", i, res.Key(), res.Val())
		} else if toTestItem(res.Val().(*ESkipNode)) != &its[i] {
			t.Errorf("invalid value: %v", i)
		}
	}
}

```