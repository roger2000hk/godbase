package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"testing"
)

func TestSuffix(t *testing.T) {
	var s Suffix

	// keys must be of type godbase.StrKey
	// per key dup check control is inherited from the set api

	s.Insert(0, godbase.StrKey("abc"), false)
	s.Insert(0, godbase.StrKey("abcdef"), false)
	s.Insert(0, godbase.StrKey("abcdefghi"), false)

	// find first suffix starting with "de" using wrapped Find()
	i  := s.First(0, godbase.StrKey("de"))
	
	// then we get all matching suffixes in order
	// i+1 since we matched on a prefix instead of full key
	if k := s.Get(nil, i+1).(godbase.StrKey); k != "def" {
		t.Errorf("invalid find res: %v", k)
	}

	// then we get all matching suffixes in order
	if k := s.Get(nil, i+2).(godbase.StrKey); k != "defghi" {
		t.Errorf("invalid find res: %v", k)
	}

	// check that Delete removes all suffixes

	if _, cnt := s.DeleteAll(0, 0, godbase.StrKey("abcdefghi")); cnt != 8 {
		t.Errorf("invalid delete res: %v", cnt)
	}
}
