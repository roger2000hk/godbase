package defs

import (
	"github.com/fncodr/godbase"
	"strings"
)

type Basic struct {
	name string
}

func (d *Basic) Init(n string) *Basic {
	d.name = n
	return d
}

func (d *Basic) Less(other godbase.Key) bool {
	return strings.Compare(d.name, other.(godbase.Def).Name()) < 0
}

func (d *Basic) Name() string {
	return d.name
}
