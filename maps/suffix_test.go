package maps

import (
       //"fmt"
       "testing"
)

func TestSuffix(t *testing.T) {
	// NewSuffix wraps any map
	// iters only work within slots for hash maps; therefore, the obvious 
	// combination is with one of the sorted maps.

	m := NewSuffix(NewSort(4))

	// keys must be of type StringKey
	// per key dup check control is inherited from the map api

	m.Insert(nil, StringKey("abc"), "abc", true)
	m.Insert(nil, StringKey("abcdef"), "abcdef", true)
	m.Insert(nil, StringKey("abcdefghi"), "abcdefghi", true)

	// find first suffix starting with "de" using wrapped Find()
	i, _ := m.Find(nil, StringKey("de"), nil)
	
	// since we're prefix searching, iter needs to be stepped once
	i = i.Next()

	// then we get all matching suffixes in order
	if i.Key().(StringKey) != "def" || i.Val().(string) != "abcdef" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	i = i.Next()

	if i.Key().(StringKey) != "defghi" || i.Val().(string) != "abcdefghi" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	// check that Delete removes all suffixes for specified val
	if res, cnt := m.Delete(nil, nil, StringKey("bcdef"), "abcdef"); 
	cnt != 4 || res.Next().Key().(StringKey) != "cdefghi" {
		t.Errorf("invalid delete res: %v", res.Next().Key())	
	}
}
