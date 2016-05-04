# sets
#### sorted & hashed sets based on slices

### design
This package contains a sorted set implementation based on sorted slices, a corresponding hashed implmentation is also provided. They are very close to godbase's maps in spirit, not so much in body. Where the maps are optimized for huge datasets and frequent modifications, the sets are lightweight and snappy as long as they don't get to big and are read more than written. Even though a lot of care has been taken to minimize the number of separate memory allocations in the maps package, nothing beats allocating your entire dataset as a single slab of memory and accessing by index.

### status
Basic functionality and testing in place, and evolving on a daily basis. The rest of godbase is currently being built on top, see recs/basic.go for a real world example.

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
	// specify end=0 for rest of set
	// returns next idx and nr of deleted elems

	DeleteAll(start, end int, key Key) (int, int64)

	// returns first index of key, from start; regardless of actually finding key
	// second result is true if key was found
	First(start int, key Key) (int, bool)

	// returns elem at index i, within slot for hash sets; key is ignored for sorted sets
	Get(key Key, i int) Key

	// returns last index of key, from start to end (exclusive); 
	// regardless of actually finding key
	// specify end=0 for rest of set
	// second result is true if key was found

	Last(start, end int, key Key) (int, bool)

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
Extending the sets api is as simple as embedding one of the implementations in your struct and optionally overriding parts of the api. sets.Suffix implements a suffix set on top of sets.Sort:

```go

type Suffix struct {
	Sort
}

// override to delete all suffixes
func (self *Suffix) Delete(start int, key godbase.Key) int {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Delete(start, godbase.StrKey(sk[i:]))
	}

	return self.Sort.Delete(start, sk)
}

// override to delete all suffixes
func (self *Suffix) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	sk := key.(godbase.StrKey)
	res := int64(0)

	for i := 1; i < len(sk) - 1; i++ {
		_, cnt := self.Sort.DeleteAll(start, end, godbase.StrKey(sk[i:]))
		res += cnt
	}

	i, cnt := self.Sort.DeleteAll(start, end, sk)
	return i, res + cnt 
}

// override to insert all suffixes
func (self *Suffix) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Insert(start, godbase.StrKey(sk[i:]), multi)
	}

	return self.Sort.Insert(start, key, multi)
}

func TestSuffix(t *testing.T) {
	var s Suffix

	// keys must be of type godbase.StrKey
	// per key dup check control is inherited from the set api

	s.Insert(0, godbase.StrKey("abc"), false)
	s.Insert(0, godbase.StrKey("abcdef"), false)
	s.Insert(0, godbase.StrKey("abcdefghi"), false)

	// find first suffix starting with "de" using wrapped Find()
	i, ok := s.First(0, godbase.StrKey("de"))

	// we shouldn't get a clean find on "de"
	if ok {
		t.Errorf("found: %v", i)		
	}

	// then we get all matching suffixes in order
	// i+1 since we matched on a prefix instead of full key
	if k := s.Get(nil, i+1).(godbase.StrKey); k != "def" {
		t.Errorf("invalid find res: %v", k)
	}

	// then we get all matching suffixes in order
	if k := s.Get(nil, i+2).(godbase.StrKey); k != "defghi" {
		t.Errorf("invalid find res: %v", k)
	}

	// check that Delete removes all suffixes

	if _, cnt := s.DeleteAll(0, 0, godbase.StrKey("abcdefghi")); cnt != 8 {
		t.Errorf("invalid delete res: %v", cnt)
	}
}

```

### wrapping it up
godbase provides scaffolding for trivial ad-hoc extension of the sets api in form of a Wrap struct. sets.Trace serves well as an introduction:

```go

// embedding Wrap gives you default delegation to wrapped set

type Trace struct {
	Wrap

	// we're adding an id for logging
	id string
}

func NewTrace(s godbase.Set, id string) *Trace {
	res := new(Trace)
	res.Init(s)
	res.id = id
	return res
}

// override to log actions before updating wrapped map

func (self *Trace) Delete(start int, key godbase.Key) int {
	log.Printf("%v.Delete %v: '%v'", self.id, start, key)
	return self.wrapped.Delete(start, key)
}

func (self *Trace) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	log.Printf("%v.DeleteAll %v/%v: '%v'", self.id, start, end, key)
	return self.wrapped.DeleteAll(start, end, key)
}

func (self *Trace) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	log.Printf("%v.Insert/%v %v: '%v'", self.id, multi, start, key)
	return self.wrapped.Insert(start, key, multi)
}

```