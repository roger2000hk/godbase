package cols

import (
	"bytes"
	"testing"
)

func TestReadWriteVal(t *testing.T) {
	var buf bytes.Buffer
	col := NewInt64("foo")

	col.WriteVal(int64(42), &buf)
	if v, err := col.ReadVal(&buf); err != nil {
		t.Error(err)
	} else if v.(int64) != 42 {
		t.Errorf("invalid val read: %v", v)
	}
}
