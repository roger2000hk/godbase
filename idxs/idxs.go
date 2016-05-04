package idxs

import (
	"fmt"
	"hash"
	"hash/fnv"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/maps"
)

type Map struct {
	Basic
	cols []godbase.Col
	hash hash.Hash64
	recs godbase.Map
	unique bool
}

type Key1 [1]godbase.Key
type Key2 [2]godbase.Key
type Key3 [3]godbase.Key

type DupKey struct {
	key godbase.Key
}

type KeyNotFound struct {
	key godbase.Key
}

func (self *Map) AddToTbl(tbl godbase.Tbl)  {
	tbl.AddIdx(self)
}

func (i *Map) Delete(start godbase.Iter, r godbase.Rec) error {
	return i.Drop(start, r)
}

func (i *Map) Drop(start godbase.Iter, r godbase.Rec) error {
	k := i.RecKey(r)

	if _, cnt := i.recs.Delete(start, nil, k, r.Id()); cnt == 0 {
		return &KeyNotFound{key: k}
	}

	return nil
}

func (e *DupKey) Error() string {
	return fmt.Sprintf("dup key: %v", e.key)
}

func (e *KeyNotFound) Error() string {
	return fmt.Sprintf("key not found: %v", e.key)
}

func (i *Map) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return i.recs.Find(start, key, val)
}

func (i *Map) Init(n string, cs []godbase.Col, u bool, rs godbase.Map) *Map {
	i.Basic.Init(n)
	i.cols = cs
	i.hash = fnv.New64()
	i.recs = rs
	i.unique = u
	return i
}

func (i *Map) Insert(start godbase.Iter, r godbase.Rec) (godbase.Iter, error) {
	k := i.RecKey(r)
	res, ok := i.recs.Insert(start, k, r.Id(), !i.unique)

	if !ok && res.Val() != r.Id() {
		return nil, &DupKey{key: k}
	} 

	return res, nil
}

func (i *Map) Load(rec godbase.Rec) (godbase.Rec, error) {
	i.recs.Set(i.RecKey(rec), rec.Id())
	return rec, nil
}

func (i *Map) Key(ks...interface{}) godbase.Key {
	il, kl := len(i.cols), len(ks)
	var k1, k2, k3 godbase.Key
	
	k1 = i.cols[0].AsKey(nil, ks[0])
	
	if kl > 1 {
		k2 = i.cols[1].AsKey(nil, ks[1])
	}

	if kl > 2 {
		k3 = i.cols[2].AsKey(nil, ks[2])
	}	

	if il == 1 {
		return Key1{k1}
	}	

	if il == 2 {
		return Key2{k1, k2}
	}

	return Key3{k1, k2, k3}
}

func (i *Map) RecKey(r godbase.Rec) godbase.Key {
	l := len(i.cols)
	var k1, k2, k3 godbase.Key
	
	if v, ok := r.Find(i.cols[0]); ok {
		k1 = i.cols[0].AsKey(r, v)
	}
	
	if l == 1 {
		return Key1{k1}
	}	

	if v, ok := r.Find(i.cols[1]); ok {
		k2 = i.cols[1].AsKey(r, v)
	}

	if l == 2 {
		return Key2{k1, k2}
	}

	if v, ok := r.Find(i.cols[2]); ok {
		k3 = i.cols[2].AsKey(r, v)
	}
	
	return Key3{k1, k2, k3}
}

func (k Key1) Less(_other godbase.Key) bool {
	other := _other.(Key1)
	nil0 := k[0] == nil || other[0] == nil
 
	return !nil0 && k[0].Less(other[0])
}

func (k Key2) Less(_other godbase.Key) bool {
	other := _other.(Key2)
	nil0, eq0 := k[0] == nil || other[0] == nil, k[0] == other[0]
	nil1 := k[1] == nil || other[1] == nil

	return (nil0 || k[0].Less(other[0])) ||
		((nil0 || eq0) && !nil1 && k[1].Less(other[1]))
}

func (k Key3) Less(_other godbase.Key) bool {
	other := _other.(Key3)
	nil0, eq0 := k[0] == nil || other[0] == nil, k[0] == other[0]
	nil1, eq1 := k[1] == nil || other[1] == nil, k[1] == other[1]
	nil2 := k[2] == nil || other[2] == nil

	return (nil0 || k[0].Less(other[0])) ||
		((nil0 || eq0) && !nil1 && k[1].Less(other[1])) ||
		((nil0 || eq0) && (nil1 || eq1) && !nil2 && k[2].Less(other[2]))
}

func New(n string, cs []godbase.Col, u bool, recs godbase.Map) *Map {
	return new(Map).Init(n, cs, u, recs)
}

func NewHash(n string, cs []godbase.Col, u bool, sc int, a *maps.SlabAlloc, ls int) *Map {
	i := new(Map)
	return i.Init(n, cs, u, maps.NewSlabHash(sc, genHashFn(i), a, ls))
}

func NewSort(n string, cs []godbase.Col, u bool, a *maps.SlabAlloc, ls int) *Map {
	return New(n, cs, u, maps.NewSlab(a, ls))
}

func genHashFn(i *Map) func(godbase.Key) uint64 {
	return func(_key godbase.Key) uint64 {
		i.hash.Reset()
		l := len(i.cols)

		switch l {
		case 1: 
			i.cols[0].Hash(nil, _key.(Key1)[0], i.hash)
		case 2: 
			key := _key.(Key2)
			i.cols[0].Hash(nil, key[0], i.hash)
			i.cols[1].Hash(nil, key[1], i.hash)
		case 3: 
			key := _key.(Key3)
			i.cols[0].Hash(nil, key[0], i.hash)
			i.cols[1].Hash(nil, key[1], i.hash)
			i.cols[2].Hash(nil, key[2], i.hash)
		default:
			panic(fmt.Sprintf("invalid idx key len: %v", l))
		}
	
		return i.hash.Sum64()
	}
}
