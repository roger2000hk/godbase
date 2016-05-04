# sets
#### sorted & hashed sets based on slices

### design
This package contains a sorted set implementation based on sorted slices, a corresponding hashed implmentation is also provided. They are very close to godbase's maps in spirit, not so much in body. Where the maps are optimized for huge datasets and frequent modifications, the sets are lightweight and snappy as long as they don't get to big and are read more than written. Even though a lot of care has been taken to minimize the number of separate memory allocations in the maps package, nothing beats allocating your entire dataset as a single slab of memory and accessing by index.

### status
Basic functionality and testing in place, and evolving on a daily basis. The rest of godbase is currently being built on top, there are plenty of examples in the other packages.

### performance
The short story is that the sorted implementation is around 50 times slower than a native hashed map for 10k elems; a properly tuned hashed set is comparable to a native map, while still providing the same features as the sorted implementation. One spot where sorted sets blow native maps out of the water is cloning. Basic performance tests are in sets_test.go

```
	go test -bench=.*
```

### interfaces

```go

// all set keys are requred to support the key interface
type Key interface {
	Less(Key) bool
}

// all hash sets require a hash fn
type HashFn func (Key) uint64

// the set interface is supported by all implementations

type Set interface {
	// returns clone of set
	Clone() Set

	// deletes key from start,
	// returns next idx, or -1 if not found

	Delete(start int, key Key) int

	// deletes keys from start to end (exclusive),
	// returns next idx, or -1 if not found; and nr of deleted elems

	DeleteAll(start, end int, key Key) (int, int64)

	// returns first index of key, from start; or -1 if not found
	First(start int, key Key) int

	// returns elem at index i, within slot for hash sets; key is ignored for sorted sets
	Get(key Key, i int) Key

	// returns last index of key, from start to end (exclusive); or -1 if not found
	Last(start, end int, key Key) int

	// inserts key into set, from start; rejects dup keys if multi=false
	// returns updated set and final index, or org set and -1 if dup

	Insert(start int, key Key, multi bool) (int, bool)

	// returns number of elems in set
	Len() int64

	// calls fn with successive elems until it returns false; returns false on early exit
	While(IKTestFn) bool
}

// callbacks

type IKTestFn func (int, Key) bool


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
	
	NewMap(100)

	// a sorted set
	NewSort()

	// note that the zero value works fine as well
	// var sortSet Sort

	// hashed set with 1000 slots
	NewSortHash(1000, genHash)
}

```

### extending
Extending the map api is as simple as embedding one of the implementations in your struct and optionally overriding parts of the api. maps.Suffix implements a suffix map on top of maps.Sorted:

```

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

	// keys must be of type godbase.StringKey
	// per key dup check control is inherited from the map api

	m.Insert(nil, godbase.StringKey("abc"), "abc", true)
	m.Insert(nil, godbase.StringKey("abcdef"), "abcdef", true)
	m.Insert(nil, godbase.StringKey("abcdefghi"), "abcdefghi", true)

	// find first suffix starting with "de" using wrapped Find()
	i, _ := m.Find(nil, godbase.StringKey("de"), nil)
	
	// since we're prefix searching, iter needs to be stepped once
	i = i.Next()

	// then we get all matching suffixes in order
	if i.Key().(godbase.StringKey) != "def" || i.Val().(string) != "abcdef" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	i = i.Next()

	if i.Key().(godbase.StringKey) != "defghi" || i.Val().(string) != "abcdefghi" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	// check that Delete removes all suffixes for specified val
	if res, cnt := m.Delete(nil, nil, godbase.StringKey("bcdef"), "abcdef"); 
	cnt != 4 || res.Next().Key().(godbase.StringKey) != "cdefghi" {
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
	log.Printf("%v.delete '%v': '%v'", self.id, key, val)
	return self.wrapped.Delete(start, end, key, val)
}

func (self *Trace) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	multi bool) (godbase.Iter, bool) {
	log.Printf("%v.insert/%v '%v': '%v'", self.id, multi, key, val)
	return self.wrapped.Insert(start, key, val, multi)
}

```