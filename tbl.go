package godbase

import (
	"io"
	"time"
)

type Tbl interface {
	Def
	AddCol(Col) Col
	Clear()
	Col(string) Col
	Cols() Iter
	Dump(io.Writer) error
	Len() int64
	Load(Cx, Rec) (Rec, error)
	Reset(Rec) (Rec, error)
	Read(Rec, io.Reader) (Rec, error)
	Revision(Rec) int64
	OnDelete() *Evt
	OnLoad() *Evt
	OnUpsert() *Evt
	Slurp(Cx, io.Reader) error
	Upsert(Cx, Rec) (Rec, error)
	UpsertedAt(Rec) time.Time
	Write(Rec, io.Writer) error
}
