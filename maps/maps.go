package maps

// package maps implements sorted & hashed maps based on deterministic skip lists

import (
	"github.com/fncodr/godbase"
)

// map allocator interface
type Alloc func () godbase.Map
