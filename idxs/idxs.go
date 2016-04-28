package idxs

import (
	"fmt"
	"hash"
	"hash/fnv"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
)

type Any interface {
	Insert(recs.Any) (recs.Any, error)
	Key(r recs.Any) interface{}
	Delete(recs.Any) error
}

type Basic struct {
	cols []cols.Any
	hash hash.Hash64
	recs maps.Any
	unique bool
}

type Key1 [1]maps.Key
type Key2 [2]maps.Key
type Key3 [3]maps.Key

type DupKey struct {
	key maps.Key
}

type KeyNotFound struct {
	key maps.Key
}

func (i *Basic) Delete(r recs.Any) error {
	k := i.Key(r)

	if _, cnt := i.recs.Delete(nil, nil, k, r.Id()); cnt == 0 {
		return &KeyNotFound{key: k}
	}

	return nil
}

func (e *DupKey) Error() string {
	return fmt.Sprintf("dup key: %v", e)
}

func (e *KeyNotFound) Error() string {
	return fmt.Sprintf("key not found: %v", e)
}

func (i *Basic) Init(rs maps.Any, cs []cols.Any, u bool) *Basic {
	i.cols = cs
	i.hash = fnv.New64()
	i.recs = rs
	i.unique = u
	return i
}

func (i *Basic) Key(r recs.Any) maps.Key {
	l := len(i.cols)
	var k1, k2, k3 maps.Key
	
	switch l {
	case 3: 
		if v, ok := r.Find(i.cols[2]); ok {
			k3 = i.cols[2].AsKey(v)
		}
		fallthrough
	case 2: 
		if v, ok := r.Find(i.cols[1]); ok {
			k2 = i.cols[1].AsKey(v)
		}
		fallthrough
	case 1: 
		if v, ok := r.Find(i.cols[0]); ok {
			k1 = i.cols[0].AsKey(v)
		}
	default:
		panic(fmt.Sprintf("invalid idx key len: %v", l))
	}

	switch l {
	case 1: return Key1{k1}
	case 2: return Key2{k1, k2}
	}
	
	return Key3{k1, k2, k3}
}

func (i *Basic) Insert(r recs.Any) (recs.Any, error) {
	k := i.Key(r)

	if res, ok := i.recs.Insert(nil, k, r.Id(), !i.unique); !ok && res.Val() != r.Id() {
		return nil, &DupKey{key: k}
	}

	return r, nil
}

func (k Key1) Less(other maps.Key) bool {
	return k[0].Less(other.(Key1)[0])
}

func (k Key2) Less(_other maps.Key) bool {
	other := _other.(Key2)
	return k[0].Less(other[0]) || k[1].Less(other[1])
}

func (k Key3) Less(_other maps.Key) bool {
	other := _other.(Key3)
	return k[0].Less(other[0]) || k[1].Less(other[1]) || k[2].Less(other[2])
}

func NewHash(cs []cols.Any, u bool, sc int, a *maps.SlabAlloc, ls int) *Basic {
	i := new(Basic)
	return i.Init(maps.NewHash(maps.NewSlabSlots(sc, genHashFn(i), a, ls)), cs, u)
}

func NewSorted(cs []cols.Any, u bool, a *maps.SlabAlloc, ls int) *Basic {
	return new(Basic).Init(maps.NewSlab(a, ls), cs, u)
}

func genHashFn(i *Basic) func(maps.Key) uint64 {
	return func(_key maps.Key) uint64 {
		i.hash.Reset()
		l := len(i.cols)

		switch l {
		case 1: 
			i.cols[0].Hash(_key.(Key1)[0], i.hash)
		case 2: 
			key := _key.(Key2)
			i.cols[0].Hash(key[0], i.hash)
			i.cols[1].Hash(key[1], i.hash)
		case 3: 
			key := _key.(Key3)
			i.cols[0].Hash(key[0], i.hash)
			i.cols[1].Hash(key[1], i.hash)
			i.cols[2].Hash(key[2], i.hash)
		default:
			panic(fmt.Sprintf("invalid idx key len: %v", l))
		}
	
		return i.hash.Sum64()
	}
}
