package maps

import (
       //"fmt"
       "testing"
)

func runSuffixTests(t *testing.T, alloc Alloc) {
	m := NewSuffix(alloc())
	m.Insert(nil, StringKey("abc"), 1, true)
	m.Insert(nil, StringKey("abcdef"), 2, true)
	m.Insert(nil, StringKey("abcdefghi"), 3, true)

	i, _ := m.Find(nil, StringKey("de"), nil)
	i = i.Next()

	if i.Key().(StringKey) != "def" || i.Val().(int) != 2 {
		t.Errorf("invalid find res: %v", i.Val())
	}

	i = i.Next()
	if i.Key().(StringKey) != "defghi" || i.Val().(int) != 3 {
		t.Errorf("invalid find res: %v", i.Key())
	}

	if res, cnt := m.Delete(nil, nil, StringKey("bcdef"), 2); 
	cnt != 4 || res.Next().Key().(StringKey) != "bcdefghi" {
		t.Errorf("invalid delete res: %v", res.Next().Key())	
	}


}

func TestSuffix(t *testing.T) {
	allocSkip := func() Any {
		return NewSkip(testSkipAlloc, testLevels)
	}

	runSuffixTests(t, allocSkip)
}
