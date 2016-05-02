package maps

// package godbase/maps implements hashed and sorted maps based on skip lists with optionally 
// slab allocated or embedded nodes

import (
	"github.com/fncodr/godbase"
)

// map allocator interface
type Alloc func () godbase.Map
