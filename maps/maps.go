package maps

import (
	"github.com/fncodr/godbase"
)

// map allocator interface
type Alloc func () godbase.Map
