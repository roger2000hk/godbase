package godbase

type Idx interface {
	TblDef
	Delete(Iter, Rec) error
	Drop(Iter, Rec) error
	Find(start Iter, key Key, val interface{}) (Iter, bool)
	Insert(Iter, Rec) (Iter, error)
	Load(Rec) (Rec, error)
	Key(...interface{}) Key
	RecKey(r Rec) Key
}
