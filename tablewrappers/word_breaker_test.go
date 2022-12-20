package tablewrappers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBreakWord(t *testing.T) {
	t.Parallel()

	t.Run("should break word on non-letter|digit boundaries", func(t *testing.T) {
		t.Parallel()

		lvl := breakOnSeps
		require.Equal(t,
			[]string{"abcdefg"},
			breakWord("abcdefg", 4, lvl),
		)
		require.Equal(t,
			[]string{"abcd|", "efg"},
			breakWord("abcd|efg", 4, lvl),
		)

		require.Equal(t,
			[]string{"1234.", "34"},
			breakWord("1234.34", 4, lvl),
		)

		require.Equal(t,
			[]string{"1234.", "345."},
			breakWord("1234.345.", 4, lvl),
		)

		require.Equal(t,
			[]string{"ABC|", "1234.", "345."},
			breakWord("ABC|1234.345.", 7, lvl),
		)

		require.Equal(t,
			[]string{"ABC|1234.", "345"},
			breakWord("ABC|1234.345", 9, lvl),
		)
	})

	t.Run("should break word anywhere", func(t *testing.T) {
		t.Parallel()

		lvl := breakAnywhere

		require.Equal(t,
			[]string{"abcd", "efg"},
			breakWord("abcdefg", 4, lvl),
		)

		require.Equal(t,
			[]string{"abcd", "|efg"},
			breakWord("abcd|efg", 4, lvl),
		)

		require.Equal(t,
			[]string{"1234", ".34"},
			breakWord("1234.34", 4, lvl),
		)

		require.Equal(t,
			[]string{"1234", ".345", "."},
			breakWord("1234.345.", 4, lvl),
		)

		require.Equal(t,
			[]string{"ABC|123", "4.345."},
			breakWord("ABC|1234.345.", 7, lvl),
		)

		require.Equal(t,
			[]string{"ABC|1234.", "345"},
			breakWord("ABC|1234.345", 9, lvl),
		)
	})
}

func TestWords(t *testing.T) {
	const (
		sentence = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
		limit    = 6
	)

	tokens := strings.FieldsFunc(sentence, BlankSplitter)
	words := newWords(tokens)

	words.Sort() // sort by decreasing size
	for _, w := range words {
		t.Logf("[%d]: %#v", w.n, w.parts)
	}

	for _, word := range words {
		word.Break(limit, breakAnywhere)
	}

	words.SortNatural()
	for _, w := range words {
		t.Logf("[%d]: %#v", w.n, w.parts)
	}
}
