package cxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/maps"
)

type Basic struct {
	mapAlloc maps.SlabAlloc
}

func New(mss int) godbase.Cx {
	return new(Basic).Init(mss)
}

func (self *Basic) Init(mss int) *Basic {
	self.mapAlloc.Init(mss)
	return self
}

func (self *Basic) MapAlloc() *maps.SlabAlloc {
	return &self.mapAlloc
}
