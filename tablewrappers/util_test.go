package tablewrappers

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/require"
)

func TestDisplayWidth(t *testing.T) {
	t.Parallel()

	expectedWidth := 13
	if runewidth.IsEastAsian() {
		// NOTE(fred): this check is relative to the current locale. For CJK, display widths are altered.
		expectedWidth = 14
	}

	t.Run("should get expected display width", func(t *testing.T) {
		const input = "Česká řeřicha"
		require.Equal(t, expectedWidth, DisplayWidth(input))
	})

	t.Run("should ignore ANSI escape sequence", func(t *testing.T) {
		const input = "\033[43;30mČeská řeřicha\033[00m"
		require.Equal(t, expectedWidth, DisplayWidth(input))
	})
}

func TestStripANSI(t *testing.T) {
	t.Parallel()

	t.Run("strip ANSI should leave non-escaped string unchanged", func(t *testing.T) {
		const input = "ABC"
		stripped, start, end := stripANSI(input)

		require.Equal(t, input, stripped)
		require.Empty(t, start)
		require.Empty(t, end)
	})

	const (
		startInput = "\033[43;30m"
		endInput   = "\033[00m"
		wordInput  = "Česká řeřicha"
	)

	t.Run("strip ANSI should isolate string between start and end escape sequences", func(t *testing.T) {
		const input = startInput + wordInput + endInput

		stripped, start, end := stripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput, start)
		require.Equal(t, endInput, end)
	})

	t.Run("strip ANSI should isolate string when missing end escape sequence", func(t *testing.T) {
		const input = startInput + wordInput

		stripped, start, end := stripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput, start)
		require.Empty(t, end)
	})

	t.Run("strip ANSI should isolate string when missing start escape sequence", func(t *testing.T) {
		const input = wordInput + endInput

		stripped, start, end := stripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Empty(t, start)
		require.Equal(t, endInput, end)
	})

	t.Run("strip ANSI should isolate empty string when pure start/end escape sequence", func(t *testing.T) {
		const input = startInput + endInput

		stripped, start, end := stripANSI(input)
		require.Empty(t, stripped)
		require.Empty(t, start)
		require.Equal(t, startInput+endInput, end)
	})

	t.Run("strip ANSI should isolate string when multiple start/end escape sequences", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput

		stripped, start, end := stripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput+startInput, start)
		require.Equal(t, endInput+endInput, end)
	})

	t.Run("strip ANSI should NOT isolate string with nested start and end escape sequences (unsupported)", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput + startInput + wordInput + endInput

		stripped, start, end := stripANSI(input)
		require.Equal(t, input, stripped)
		require.Empty(t, start)
		require.Empty(t, end)
	})
}
