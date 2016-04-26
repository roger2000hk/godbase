package cols

import (
	"bytes"
	"testing"
)

func TestReadWriteVal(t *testing.T) {
	var buf bytes.Buffer
	col := NewInt64("foo")

	col.WriteVal(int64(42), &buf)

	var err error
	var s ValSize

	if s, err = ReadSize(&buf); err != nil {
		t.Error(err)
	}

	var v interface{}

	if v, err = col.ReadVal(s, &buf); err != nil {
		t.Error(err)
	} else if v.(int64) != 42 {
		t.Errorf("invalid val read: %v", v)
	}
}
