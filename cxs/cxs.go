package cxs

import (
)

type Basic struct {
}

func New() *Basic {
	return new(Basic).Init()
}

func (self *Basic) Init() *Basic {
	return self
}
