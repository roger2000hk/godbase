package godbase

type Cx interface {
}

type BasicCx struct {
}

func NewCx() *BasicCx {
	return new(BasicCx).Init()
}

func (self *BasicCx) Init() *BasicCx {
	return self
}
