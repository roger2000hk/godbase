package idxs

import (
	"testing"
)

func TestKeyEq(t *testing.T) {
	gen := func(k1, k2 string) interface{} {
		return [2]interface{}{k1, k2}
	}
	
	if  l, r := gen("abc", "def"), gen("abc", "def"); l != r {
		t.Errorf("not equal")
	}

	if  l, r := gen("abc", "def"), gen("abc", "ghi"); l == r {
		t.Errorf("equal")
	}
}
