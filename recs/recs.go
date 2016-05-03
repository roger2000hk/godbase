package recs

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type TestFn func(godbase.Rec) bool
type NotFound godbase.UId
type Size uint32

func (e NotFound) Error() string {
	return fmt.Sprintf("rec not found: %v", e)
}
