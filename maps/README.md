# maps
#### sorted & hashed maps based on deterministic skip lists

### design
This package contains two implementations of sorted maps based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes; the other with embedded nodes, constant number of levels, using one node per value. Two corresponding hashed implementations are also provided.

### status
Basic functionality and testing in place, and evolving on a daily basis. The rest of godbase is currently being built on top, there are plenty of examples in the other packages.

### performance
The short story is that sorted implementations are around 3-5 times slower than native maps for 100k elems; a properly tuned hashed one is comparable to a native map. Several parameters are available for tuning the benchmarks, they are defined in test.go

```
	go test -bench=.*
```

### interfaces

```go

// all map keys are requred to support the key interface
type Key interface {
	Less(Key) bool
}

// all hash maps require a hash fn
type HashFn func (Key) uint64


// iters are circular and cheap, 
// since they are nothing but a common interface on top of actual nodes

type Iter interface {
	// returns key for elem or nil if root
	Key() Key

	// returns iter to next elem
	Next() Iter

	// returns iter to next elem
	Prev() Iter

	// returns val for elem
	Val() interface{}

	// returns true if not root
	Valid() bool
}

// the map interface is supported by all implementations

type Map interface {
	// clears all elems from map
	// deallocates nodes for maps that use allocators

	Clear()

	// cuts elems from start to end for which fn returns non nil key into new set
	// start, end & fn are all optional
	// when fn is specified, the returned key/val replaces the original; 
	// except for maps with embedded nodes, where the returned val replaces the
	// entire node
	// no safety checks are provided; if you mess up the ordering, you're on your own
	// circular cuts, with start/end on opposite sides of root; are supported 
	// returns a cut from the start slot for hash maps

	Cut(start, end Iter, fn KVMapFn) Map

	// deletes elems from start to end, matching key/val
	// start, end, key & val are all optional, nil means all elems 
	// specifying iters for hash maps only works within the same slot 
	// circular deletes, with start/end on opposite sides of root; are supported
	// deallocates nodes for slab allocated maps
	// returns an iter to next elem and number of deleted elems

	Delete(start, end Iter, key Key, val interface{}) (Iter, int)

	// returns iter for first elem after start matching key and ok
	// start & val are optional, 
	// specifying a start iter for hash maps only works within the same slot

	Find(start Iter, key Key, val interface{}) (Iter, bool)
	
	// returns iter to first elem
	// not supported by hash maps

	First() Iter

	// returns val for key and ok
	Get(key Key) (interface{}, bool)

	// rnserts key/val into map after start
	// start & val are both optional,
	// dup checks can be disabled by setting multi to false
	// returns iter to inserted val & true on success, iter to existing val & false on dup 
	// specifying a start iter for hash maps only works within the same slot

	Insert(start Iter, key Key, val interface{}, multi bool) (Iter, bool)

	// returns a new, empty map of the same type as the receiver
	New() Map

	// returns the number of elems in map
	Len() int64

	// inserts/updates key to val and returns true on insert
	Set(key Key, val interface{}) bool

	// returns string rep of map
	String() string
	
	// calls fn with successive elems until false; returns false on early exit
	While(KVTestFn) bool
}

// callbacks

type KVMapFn func (Key, interface{}) (Key, interface{})

type KVTestFn func (Key, interface{}) bool


```

### constructors

```go

type testKey int

func (k testKey) Less(other godbase.Key) bool {
	return k < other.(testKey)
}

func genHash(k godbase.Key) uint64 { return uint64(k.(testKey)) }

func TestConstructors(t *testing.T) {
	// Map is mostly meant as a reference for performance comparisons,
	// it only supports enough of the api to run basic tests on top of a native map
	
	NewMap()
	
	// 10 level sorted map
	NewSort(10)

	// slab allocator with 50 nodes per slab
	a := NewSlabAlloc(50)

	// 20 level sorted map with slab allocated nodes
	NewSlab(a, 20)

	// sorted map with embedded nodes
	NewESort()

	// 1000 hash slots backed by 2 level maps with slab allocated nodes
	NewSlabHash(1000, genHash, a, 2)

	// 1000 hash slots backed by 2 level maps with separately allocated nodes
	NewSortHash(1000, genHash, 2)

	// 1000 hash slots backed by maps with embedded nodes
	NewESortHash(1000, genHash)
}

```

### embedded nodes
I picked up the idea of embedding node infrastructure into elems from the Linux kernel, but I'm sure the idea is at least as old as the C language. It's a nice tool to reduce memory allocation which bends the rules enough for the previously undoable to become possible. If you don't mind keeping a reference per collection in your type, or sprinkling a pinch of unsafe magic on top; this might be for you. lists.EDouble contains a double-linked list implementation based on the same idea.

```go

// pretend this is your value type
type testItem struct {
	node ENode

	// additional fields...
}

// calculate offset of node within struct
var testItemOffs = unsafe.Offsetof(new(testItem).node)

// helper to get pointer to item from node
func toTestItem(node *ENode) *testItem {
	return (*testItem)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - testItemOffs))
}

func TestEmbedded(t *testing.T) {
	m := NewESort()
	
	const n = 100
	its := make([]testItem, n)

	for i := 0; i < n; i++ {
		k := testKey(i)

		// the map only deals with nodes,
		// translation to/from values is left to client code
		m.Insert(nil, k, &its[i].node, false)

		res, ok := m.Find(nil, k, nil)
		res = res.Next()

		if !ok || res.Key() != k || res.Val().(*ENode) != &its[i].node {
			t.Errorf("invalid iter: %v/%v/%v", i, res.Key(), res.Val())
		} else if toTestItem(res.Val().(*ENode)) != &its[i] {
			t.Errorf("invalid value: %v", i)
		}
	}
}

```

### extending
Extending the map api is as simple as embedding one of the implementations in your struct and optionally overriding parts of the api. maps.Suffix implements a suffix map on top of maps.Sort:

```go

type Suffix struct {
	Sort
}

func NewSuffix(a *SlabAlloc, ls int) *Suffix {
	res := new(Suffix)
	res.Sort.Init(a, ls)
	return res
}

// override to delete all suffixes
func (self *Suffix) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	sk := key.(godbase.StrKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := self.Sort.Delete(start, end, godbase.StrKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := self.Sort.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

// override to insert all suffixes
func (self *Suffix) Insert(start godbase.Iter, key godbase.Key, val interface{}, multi bool) (godbase.Iter, bool) {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Insert(start, godbase.StrKey(sk[i:]), val, multi)
	}

	return self.Sort.Insert(start, key, val, multi)
}

func TestSuffix(t *testing.T) {
	m := NewSuffix(nil, 3)

	// keys must be of type godbase.StrKey
	// per key dup check control is inherited from the map api

	m.Insert(nil, godbase.StrKey("abc"), "abc", true)
	m.Insert(nil, godbase.StrKey("abcdef"), "abcdef", true)
	m.Insert(nil, godbase.StrKey("abcdefghi"), "abcdefghi", true)

	// find first suffix starting with "de" using wrapped Find()
	i, ok := m.Find(nil, godbase.StrKey("de"), nil)
	
	// we shouldn't get a clean find on "de"
	if ok {
		t.Errorf("found: %v", i)		
	}

	// since we're prefix searching, iter needs to be stepped once
	i = i.Next()

	// then we get all matching suffixes in order
	if i.Key().(godbase.StrKey) != "def" || i.Val().(string) != "abcdef" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	i = i.Next()

	if i.Key().(godbase.StrKey) != "defghi" || i.Val().(string) != "abcdefghi" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	// check that Delete removes all suffixes for specified val
	if res, cnt := m.Delete(nil, nil, godbase.StrKey("bcdef"), "abcdef"); 
	cnt != 4 || res.Next().Key().(godbase.StrKey) != "cdefghi" {
		t.Errorf("invalid delete res: %v", res.Next().Key())	
	}
}

```

### wrapping it up
godbase provides scaffolding for trivial ad-hoc extension of the maps api in form of a Wrap struct. maps.Trace serves well as an introduction:

```go

// embedding Wrap gives you default delegation to wrapped map

type Trace struct {
	Wrap

	// we're adding an id for logging
	id string
}

func NewTrace(m godbase.Map, id string) *Trace {
	res := new(Trace)
	res.Init(m)
	res.id = id
	return res
}

// override to log actions before updating wrapped map

func (self *Trace) Delete(start, end godbase.Iter, key godbase.Key, 
	val interface{}) (godbase.Iter, int) {
	log.Printf("%v.Delete '%v': '%v'", self.id, key, val)
	return self.wrapped.Delete(start, end, key, val)
}

func (self *Trace) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	multi bool) (godbase.Iter, bool) {
	log.Printf("%v.Insert/%v '%v': '%v'", self.id, multi, key, val)
	return self.wrapped.Insert(start, key, val, multi)
}

```