package idxs

import (
	"testing"
)

func TestKeyEq(t *testing.T) {
	if [2]interface{}{"abc", "def"} != [2]interface{}{"abc", "def"} {
		t.Errorf("not equal")
	}
}
