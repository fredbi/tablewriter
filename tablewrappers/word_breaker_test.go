package tablewrappers

import (
	"testing"
)

func TestBreakWord(t *testing.T) {
	t.Run("should break word on non-letter|digit boundaries", func(t *testing.T) {
		lvl := breakOnBoundaries
		t.Logf("%#v", breakWord("abcdefg", 4, lvl))
		t.Logf("%#v", breakWord("abcd|efg", 4, lvl))
		t.Logf("%#v", breakWord("1234.34", 4, lvl))
		t.Logf("%#v", breakWord("1234.345.", 4, lvl))
		t.Logf("%#v", breakWord("ABC|1234.345.", 7, lvl))
		t.Logf("%#v", breakWord("ABC|1234.345", 9, lvl))
	})

	t.Run("should break word anywhere", func(t *testing.T) {
		lvl := breakAnywhere
		t.Logf("%#v", breakWord("abcdefg", 4, lvl))
		t.Logf("%#v", breakWord("abcd|efg", 4, lvl))
		t.Logf("%#v", breakWord("1234.34", 4, lvl))
		t.Logf("%#v", breakWord("1234.345.", 4, lvl))
		t.Logf("%#v", breakWord("ABC|1234.345.", 7, lvl))
		t.Logf("%#v", breakWord("ABC|1234.345", 9, lvl))
	})
}
