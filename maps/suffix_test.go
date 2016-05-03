package maps

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"testing"
)

func TestSuffix(t *testing.T) {
	m := NewSuffix(nil, 3)

	// keys must be of type godbase.StrKey
	// per key dup check control is inherited from the map api

	m.Insert(nil, godbase.StrKey("abc"), "abc", true)
	m.Insert(nil, godbase.StrKey("abcdef"), "abcdef", true)
	m.Insert(nil, godbase.StrKey("abcdefghi"), "abcdefghi", true)

	// find first suffix starting with "de" using wrapped Find()
	i, _ := m.Find(nil, godbase.StrKey("de"), nil)
	
	// since we're prefix searching, iter needs to be stepped once
	i = i.Next()

	// then we get all matching suffixes in order
	if i.Key().(godbase.StrKey) != "def" || i.Val().(string) != "abcdef" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	i = i.Next()

	if i.Key().(godbase.StrKey) != "defghi" || i.Val().(string) != "abcdefghi" {
		t.Errorf("invalid find res: %v", i.Key())
	}

	// check that Delete removes all suffixes for specified val
	if res, cnt := m.Delete(nil, nil, godbase.StrKey("bcdef"), "abcdef"); 
	cnt != 4 || res.Next().Key().(godbase.StrKey) != "cdefghi" {
		t.Errorf("invalid delete res: %v", res.Next().Key())	
	}
}
