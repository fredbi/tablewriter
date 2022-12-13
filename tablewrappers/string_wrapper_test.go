package tablewrappers

import (
	"fmt"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/require"
)

func TestDefaultWrapper(t *testing.T) {
	t.Parallel()

	const text = "The quick brown fox jumps over the lazy dog."

	t.Run("should wrap text in lines", func(t *testing.T) {
		t.Parallel()

		const maxWidth = 6

		expected := []string{
			"The", "quick", "brown", "fox",
			"jumps", "over", "the", "lazy", "dog."}

		w := NewDefault()
		got := w.WrapString(text, maxWidth)
		require.Len(t, got, len(expected))
	})

	t.Run("should wrap into a single line", func(t *testing.T) {
		t.Parallel()

		const maxWidth = 500

		expected := "The quick brown fox jumps over the lazy dog."

		w := NewDefault()
		actual := w.WrapString(text, maxWidth)
		require.Len(t, actual, 1)
		require.EqualValues(t, []string{expected}, actual)
	})

	t.Run("should wrap unicode according to the displayed width, not the number of runes", func(t *testing.T) {
		t.Parallel()

		w := NewDefault()
		for _, toPin := range []struct {
			Input         string
			MaxWidth      int
			ExpectedLines int
		}{
			{
				Input:         "Česká řeřicha",
				MaxWidth:      13,
				ExpectedLines: 1,
			},
			{
				Input:         "〒177-0034 日本, 東京都練馬区富士見台1丁目9番2号",
				MaxWidth:      13,
				ExpectedLines: 2,
			},
		} {
			testCase := toPin

			// TODO(fred): apparently there is a catch with CJK locale.
			// See runewidth package to understand better.
			t.Run(fmt.Sprintf("%q should fit in %d lines of width %d - CJK locale:%t",
				testCase.Input, testCase.ExpectedLines, testCase.MaxWidth, runewidth.IsEastAsian(),
			), func(t *testing.T) {
				t.Parallel()

				actual := w.WrapString(testCase.Input, testCase.MaxWidth)
				require.Len(t, actual, testCase.ExpectedLines)
			})
		}
	})

	t.Run("should wrap string into paragraphs", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"ああああああああああああああああああああああああ",
			"あああああああ",
		}

		w := NewDefault()
		actual := w.WrapString(
			"ああああああああああああああああああああああああ あああああああ",
			55,
		)

		require.EqualValues(t, expected, actual)
	})
}
