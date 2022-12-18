package tablewrappers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBreakWord(t *testing.T) {
	t.Parallel()

	t.Run("should break word on non-letter|digit boundaries", func(t *testing.T) {
		t.Parallel()

		lvl := breakOnBoundaries
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
