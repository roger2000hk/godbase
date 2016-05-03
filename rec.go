package godbase

type ColValTestFn func(Col, interface{}) bool

type Rec interface {
	Clear()
	Clone() Rec
	Delete(Col) bool
	Eq(Rec) bool
	Find(Col) (interface{}, bool)
	Get(Col) interface{}
	Id() UId
	Len() int
	Set(Col, interface{}) interface{}
	While(ColValTestFn) bool
}
