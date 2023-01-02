package tablewrappers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordWrapper(t *testing.T) {
	t.Parallel()

	matrix, maxWordLength, minWordLength := buildLengthsMatrix([]string{"characters", "too", "long"}, 1)
	require.Equal(t, 10, maxWordLength)
	require.Equal(t, 3, minWordLength)
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
			wrapMultiline(words, 10, 1),
		)

		require.EqualValues(t,
			[]string{
				"1 22", "333", "4444",
			},
			wrapMultiline(words, 6, 1),
		)
	})

	longwords := []string{"1111", "2222", "33333", "44444"}

	t.Run("should wrap words with larger words", func(t *testing.T) {
		require.EqualValues(t,
			[]string{
				"1111 2222", "33333", "44444",
			},
			wrapMultiline(longwords, 9, 1),
		)

		require.EqualValues(t,
			[]string{
				"1111", "2222", "33333", "44444",
			},
			wrapMultiline(longwords, 5, 1),
		)

		require.EqualValues(t,
			[]string{
				"1111", "2222", "33333", "44444",
			},
			wrapMultiline(longwords, 4, 1),
		)
	})

	t.Run("should wrap empty list", func(t *testing.T) {
		require.EqualValues(t,
			[]string{""},
			wrapMultiline([]string{}, 4, 1),
		)
	})

	t.Run("should wrap empty words", func(t *testing.T) {
		emptyWords := []string{"", "", "", ""}
		require.EqualValues(t,
			[]string{""},
			wrapMultiline(emptyWords, 4, 1),
		)
	})

	t.Run("should wrap words without padding", func(t *testing.T) {
		require.EqualValues(t,
			[]string{
				"11112222", "33333", "44444",
			},
			wrapMultiline(longwords, 8, 0),
		)
	})

	t.Run("should account for padding width", func(t *testing.T) {
		require.EqualValues(t,
			[]string{
				"1111  2222", "33333  44444",
			},
			wrapMultiline(longwords, 12, 2),
		)
	})

	t.Run("should reconstruct escape sequence over multiple lines", func(t *testing.T) {
		const (
			startInput     = "\033[43;30m"
			endInput       = "\033[00m"
			paragraphInput = startInput + "ABC XYZ 123" + endInput
		)
		escapedInput := strings.FieldsFunc(paragraphInput, BlankSplitter)
		t.Logf("%q", wrapMultiline(escapedInput, 7, 1))
		// word_wrapper_test.go:107: ["\x1b[43;30mABC" "XYZ" "123\x1b[00m"]
		// word_wrapper_test.go:107: ["\x1b[43;30mABC XYZ" "123\x1b[00m"]
	})
}
