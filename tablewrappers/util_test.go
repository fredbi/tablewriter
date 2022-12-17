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
