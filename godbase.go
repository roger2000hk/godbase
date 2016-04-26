package godbase

type Def interface {
	Name() string
}

type BasicDef struct {
	name string
}

func (d *BasicDef) Init(n string) *BasicDef {
	d.name = n
	return d
}
