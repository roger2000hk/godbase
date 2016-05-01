package cxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
)

type Basic struct {
	mapAlloc maps.SlabAlloc
}

func New(mss int) godbase.Cx {
	return new(Basic).Init(mss)
}

func (self *Basic) InitRecId(_rec godbase.Rec, id godbase.UId) godbase.Rec {
	rec := _rec.(*recs.Basic)
	rec.Init(&self.mapAlloc)
	rec.InitId(id)
	return rec
}

func (self *Basic) InitRec(rec godbase.Rec) godbase.Rec {
	rec.(*recs.Basic).Init(&self.mapAlloc)
	return rec
}

func (self *Basic) Init(mss int) *Basic {
	self.mapAlloc.Init(mss)
	return self
}
