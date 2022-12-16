package tablewrappers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordWrapper(t *testing.T) {
	t.Parallel()

	matrix, maxWordLength := buildLengthsMatrix([]string{"characters", "too", "long"}, 1)
	require.Equal(t, 10, maxWordLength)
	require.EqualValues(t,
		[][]int{{10, 14, 19}, {0, 3, 8}, {0, 0, 4}},
		matrix,
	)
}

func TestWrapMultiline(t *testing.T) {
	t.Parallel()

	t.Run("should wrap words with small words", func(t *testing.T) {
		words := []string{"1", "22", "333", "4444"}
		require.EqualValues(t,
			[]string{
				"1 22 333", "4444",
			},
			wrapMultiline(words, 10),
		)

		require.EqualValues(t,
			[]string{
				"1 22", "333", "4444",
			},
			wrapMultiline(words, 6),
		)
	})

	t.Run("should wrap words with larger words", func(t *testing.T) {
		longwords := []string{"1111", "2222", "33333", "44444"}
		require.EqualValues(t,
			[]string{
				"1111 2222", "33333", "44444",
			},
			wrapMultiline(longwords, 9),
		)

		require.EqualValues(t,
			[]string{
				"1111", "2222", "33333", "44444",
			},
			wrapMultiline(longwords, 5),
		)

		require.EqualValues(t,
			[]string{
				"1111", "2222", "33333", "44444",
			},
			wrapMultiline(longwords, 4),
		)
	})

	t.Run("should wrap empty list", func(t *testing.T) {
		require.EqualValues(t,
			[]string{""},
			wrapMultiline([]string{}, 4),
		)
	})

	t.Run("should wrap empty words", func(t *testing.T) {
		emptyWords := []string{"", "", "", ""}
		require.EqualValues(t,
			[]string{""},
			wrapMultiline(emptyWords, 4),
		)
	})
}
