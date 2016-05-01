package godbase

type Idx interface {
	Def
	Delete(Iter, Rec) error
	Drop(Iter, Rec) error
	Find(start Iter, key Key, val interface{}) (Iter, bool)
	Insert(Iter, Rec) (Iter, error)
	Load(Rec) (Rec, error)
	Key(...interface{}) Key
	RecKey(r Rec) Key
}

type KVMapFn func (Key, interface{}) (Key, interface{})
type KVTestFn func (Key, interface{}) bool
