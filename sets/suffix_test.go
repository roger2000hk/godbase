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

	s.Insert(0, godbase.StrKey("abc"), true)
	s.Insert(0, godbase.StrKey("abcdef"), true)
	s.Insert(0, godbase.StrKey("abcdefghi"), true)

	// find first suffix starting with "de" using wrapped Find()
	i  := s.First(0, godbase.StrKey("def"))
	
	// then we get all matching suffixes in order
	if k := s.Get(nil, i).(godbase.StrKey); k != "def" {
		t.Errorf("invalid find res: %v", k)
	}

	// then we get all matching suffixes in order
	if k := s.Get(nil, i+1).(godbase.StrKey); k != "defghi" {
		t.Errorf("invalid find res: %v", k)
	}

	// check that Delete removes all suffixes

	if _, cnt := s.DeleteAll(0, -1, godbase.StrKey("abcdefghi")); cnt != 8 {
		t.Errorf("invalid delete res: %v", cnt)
	}
}
