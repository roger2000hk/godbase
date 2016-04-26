package defs

import (
	"github.com/fncodr/godbase/maps"
	"strings"
)

type Any interface {
	maps.Key
	Name() string
}

type Basic struct {
	name string
}

func (d *Basic) Init(n string) *Basic {
	d.name = n
	return d
}

func (d *Basic) Less(other maps.Key) bool {
	return strings.Compare(d.name, other.(Any).Name()) < 0
}

func (d *Basic) Name() string {
	return d.name
}
