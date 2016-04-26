package cols

import (
	"github.com/fncodr/godbase"
)

type Any interface {
	godbase.Def
	SizeOfVal(interface{}) uint64
}
