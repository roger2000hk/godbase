package recs

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"time"
)

type TestFn func(godbase.Rec) bool
type NotFound godbase.UId
type Size uint32

func (e NotFound) Error() string {
	return fmt.Sprintf("rec not found: %v", e)
}

func Init(rec *Basic, id godbase.UId) *Basic {
	rec.SetTime(cols.CreatedAt(), time.Now())
	rec.SetUId(cols.RecId(), id)
	return rec
}
