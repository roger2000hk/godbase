package idxs

import (
	"testing"
)

func TestKey(t *testing.T) {
	if [2]string{"abc", "def"} != [2]string{"abc", "def"} {
		t.Errorf("not equal")
	}
}
