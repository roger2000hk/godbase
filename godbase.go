package godbase

import (
	"encoding/binary"
)

var ByteOrder = binary.BigEndian

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
