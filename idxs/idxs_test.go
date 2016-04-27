package idxs

import (
	"testing"
)

func TestKeyEq(t *testing.T) {
	gen := func() interface{} {
		return [2]interface{}{"abc", "def"}
	}
	
	if  l, r := gen(), gen(); l != r {
		t.Errorf("not equal")
	}
}
