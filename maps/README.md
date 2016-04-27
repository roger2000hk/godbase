# beating The Map
### a glorious quest to write a faster, more flexible map

## why?
Good question. This rabbit hole was deeper than planned, but I still feel it was worth the effort given the results. To cut a long story short, I was itching for more flexible maps to implement in-memory indexing in godbase. I quickly found that my lofty goal of adding multi-capability and sorting, while matching the performance of native maps; was far from a walk in the park. Which is probably part of the reason I couldn't find anyone else trying.

## how?
I ended up with two designs based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes; the other with embedded nodes, constant number of levels, using one node per value. These do pretty good by themselves, especially considering that they add sorting to the mix; both currently hovering around 2-5 times slower than a native map. The embedded flavor usually wins the allocation race by a slim margin but pays the price of having a fixed number of levels for tiny / huge datasets. Still, somewhere along a line; not separately allocating nodes affects overall performance positively.

Once sorted maps were working properly, I had the crazy idea to put a hash on top just to see what happens. It turns out that dividing the dataset into a tuned number of ordered sets occasionally beats native maps in the synthetic performance game for millions of items. That's far from the end of the story though. I still haven't had enough time to ponder the consequences of having access to both hashed and ordered aspects of the data simultaneously, but I have a hunch it will bend the rules to my advantage in a number of tricky scenarios. Additionally; any kind of map can be hashed, which opens the door for multi level hashing where each chain is another hash that further divides the dataset along a potentially orthogonal axis.

## status
Basic functionality and testing in place, and evolving on a daily basis. The rest of godbase is currently being built on top, there are plenty of examples in the other sub packages.

## benchmarks
Several parameters are available for tuning the tests, they are defined in test.go

```
	go test -bench=.*
```

## license
NOP

## code
I trust you'll find godbase more RISC/Lispy than your everyday set/map api. Providing an optimal api is part of implementing an optimal algorithm, and there's more low hanging fruit in the garden of set/map apis than most places. It's obvious to me that academic dogmatics and software (or life in general, for that matter) isn't really the match made in heaven it's being sold as.

### interfaces

```go

// All map keys are requred to support Key
type Key interface {
	Less(Key) bool
}

// All hash maps require a hash fn
type HashFn func (Key) uint64

// Iters are circular and cheap, since they are nothing but a common 
// interface on top of actual nodes. 

type Iter interface {
	// Returns key for elem or nil if root
	Key() Key

	// Returns iter to next elem
	Next() Iter

	// Returns iter to prev elem
	Prev() Iter

	// Returns val for elem
	Val() interface{}

	// Returns true if not root
	Valid() bool
}

// Basic map ops supported by all implementations
type Any interface {
	// Clears all elems from map. Deallocates nodes for maps that use allocators.
	Clear()

	// Cuts elems from start to end for which fn returns non nil key into new set;
	// start, end & fn are all optional. When fn is specified, the returned key/val replaces
	// the original; except for maps with embedded nodes, where the returned val replaces the
	// entire node. No safety checks are provided; if you mess up the ordering, you're on your
	// own. Circular cuts, with start/end on opposite sides of root; are supported. 
	// Returns a cut from the start slot for hash maps.

	Cut(start, end Iter, fn MapFn) Any

	// Deletes elems from start to end, matching key/val;
	// start, end, key & val are all optional, nil means all elems. Specifying 
	// iters for hash maps only works within the same slot. Circular deletes,
	// with start/end on opposite sides of root; are supported. Returns an iter to next 
	// elem and number of deleted elems.

	Delete(start, end Iter, key Key, val interface{}) (Iter, int)

	// Returns iter for first elem after start matching key and ok;
	// start & val are optional, specifying a start iter for hash maps only works within the 
	// same slot.

	Find(start Iter, key Key, val interface{}) (Iter, bool)
	
	// Returns iter to first elem; not supported by hash maps
	First() Iter

	// Returns val for key and ok
	Get(key Key) (interface{}, bool)

	// Inserts key/val into map after start;
	// start & val are both optional, dup checks can be disabled by setting allowMulti to false. 
	// Returns iter to inserted val & true on success, or iter to existing val & false on dup. 
	// Specifying a start iter for hash maps only works within the same slot.

	Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool)

	// Returns a new, empty map of the same type as the receiver
	New() Any

	// Returns the number of elems in map
	Len() int64

	// Inserts/updates key to val and returns true on insert
	Set(key Key, val interface{}) bool

	// Returns string repr for printing
	String() string
}

type TestFn func (Key, interface{}) bool


```

### constructors

```go

type testKey int

func (k testKey) Less(other Key) bool {
	return k < other.(testKey)
}

func genHash(k Key) uint64 { return uint64(k.(testKey)) }
func genMapHash(k Key) interface{} { return k }

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

	// 1000 slots backed by a native array and generic slot allocator
	// could be used in any of the following examples,
	// but specializing the slot type allows allocating all slots at once and
	// accessing by value which makes a difference in some scenarios.
	// the allocator receives the key as param which enables choosing
	// differend kinds of slot chains for different keys.

	skipAlloc := func (_ Key) Any { return NewSkip(nil, 2) }
	as := NewSlots(1000, genHash, skipAlloc)
	NewHash(as)

	// 1000 slots backed by a native map and generic slot allocator
	// could also be used in any of the following examples, since it too
	// uses a generic allocator to allocate slots on demand.
	// what map slots bring to the table, is the ability to use any kind of
	// value except slices as hash keys; which is useful when
	// mapping your keys to an integer is problematic. On the other hand they
	// share the same limitations as native maps, no slice keys and relatively
	// expensive to create.

	ms := NewMapSlots(1000, genMapHash, skipAlloc)
	NewHash(ms)

	// 1000 skip slots backed by 2 level skip maps with slab allocated nodes
	ss := NewSkipSlots(1000, genHash, a, 2)
	NewHash(ss)

	// 1000 hash slots backed by embedded skip maps
	ess := NewESkipSlots(1000, genHash)
	NewHash(ess)

	// 1000 hash slots backed by hash maps with 100 embedded skip slots
	hs := NewHashSlots(1000, genHash, func (_ Key) Slots { return NewESkipSlots(100, genHash) })
	NewHash(hs)
}

```

### embedded nodes
I picked up the idea of embedding node infrastructure into elems from the Linux kernel, but I'm sure the idea is at least as old as the C language. It's a nice tool to reduce memory allocation which bends the rules enough for the previously undoable to become possible. If you don't mind keeping a reference per collection in your type, or sprinkling a pinch of unsafe magic on top; this might be for you. lists.EDouble contains a double-linked list implementation based on the same idea.

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

### wrapping it up
godbase provides scaffolding for trivial extension of the api in form of a Wrap struct, maps.Suffix serves well as an introduction:

```go

type Suffix struct {
	Wrap
}

func NewSuffix(m Any) *Suffix {
	res := new(Suffix)
	res.Init(m)
	return res
}

// override to delete all suffixes
func (m *Suffix) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	sk := key.(StringKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := m.wrapped.Delete(start, end, StringKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := m.wrapped.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

// override to insert all suffixes
func (m *Suffix) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	sk := key.(StringKey)

	for i := 1; i < len(sk) - 1; i++ {
		m.wrapped.Insert(start, StringKey(sk[i:]), val, allowMulti)
	}

	return m.wrapped.Insert(start, key, val, allowMulti)
}


```

maps.Suffix, like all wraps, can be used wherever maps.Any is expected; with the restriction that it only supports StringKeys, for obvious reasons. A suffix map is a nice tool to solve string completion problems, this one comes bundled with all the additional features of godbase map api.

```go

func TestSuffix(t *testing.T) {
	// NewSuffix wraps any map
	// iters only work within slots for hash maps; therefore, the obvious 
	// combination is with one of the ordered maps.

	m := NewSuffix(NewSkip(nil, 4))

	// keys must be of type StringKey
	// per key dup check control is inherited from the map api

	m.Insert(nil, StringKey("abc"), "abc", true)
	m.Insert(nil, StringKey("abcdef"), "abcdef", true)
	m.Insert(nil, StringKey("abcdefghi"), "abcdefghi", true)

	// find first suffix starting with "de" using wrapped Find()
	i, _ := m.Find(nil, StringKey("de"), nil)
	
	// since we're prefix searching, iter needs to be stepped once
	i = i.Next()

	// then we get all matching suffixes in order
	if i.Key().(StringKey) != "def" || i.Val().(string) != "abcdef" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	i = i.Next()

	if i.Key().(StringKey) != "defghi" || i.Val().(string) != "abcdefghi" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	// check that Delete removes all suffixes for specified val
	if res, cnt := m.Delete(nil, nil, StringKey("bcdef"), "abcdef"); 
	cnt != 4 || res.Next().Key().(StringKey) != "cdefghi" {
		t.Errorf("invalid delete res: %v", res.Next().Key())	
	}
}

```