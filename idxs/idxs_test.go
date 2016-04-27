package idxs

import (
	"testing"
)

func TestKeyEq(t *testing.T) {
	// array equality is pretty cool
	// as long as the legth can be determined at compile time,
	// any comparable values can be packed up in a fixed generic array and compared in one op

	gen := func(k1, k2 interface{}) interface{} {
		return [2]interface{}{k1, k2}
	}
	
	if  l, r := gen("abc", "def"), gen("abc", "def"); l != r {
		t.Errorf("not equal")
	}

	if  l, r := gen("abc", "def"), gen("abc", "ghi"); l == r {
		t.Errorf("equal")
	}

	if  l, r := gen(1, 2), gen(1, 2); l != r {
		t.Errorf("not equal")
	}

	if  l, r := gen("abc", "def"), gen(1, 2); l == r {
		t.Errorf("equal")
	}

	if  l, r := gen("abc", 1), gen("abc", 1); l != r {
		t.Errorf("not equal")
	}
}
