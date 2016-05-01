package godbase

import (
	"hash"
	"io"
)

type Col interface {
	Def
	AsKey(Rec, interface{}) Key
	CloneVal(interface{}) interface{}
	Decode(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(Rec, interface{}, hash.Hash64)
	Read(Rec, ValSize, io.Reader) (interface{}, error)
	Type() ColType
	Write(Rec, interface{}, io.Writer) error
}

type ColType interface {
	AsKey(Rec, interface{}) Key
	CloneVal(interface{}) interface{}
	Decode(interface{}) interface{}
	Encode(interface{}) interface{}
	Eq(interface{}, interface{}) bool
	Hash(Rec, interface{}, hash.Hash64)
	Name() string
	Read(Rec, ValSize, io.Reader) (interface{}, error)
	Write(Rec, interface{}, io.Writer) error
}
