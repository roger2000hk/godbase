package godbase

type Rec interface {
	Clear()
	Clone() Rec
	Delete(Col) bool
	Eq(Rec) bool
	Find(Col) (interface{}, bool)
	Get(Col) interface{}
	Id() UId
	Iter() Iter
	Len() int
	New() Rec
	Set(Col, interface{}) interface{}
}
