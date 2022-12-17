package tablewriter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPadDefault(t *testing.T) {
	t.Parallel()

	t.Run("should pad right (left-aligned)", func(t *testing.T) {
		const (
			toPad    = "ABC"
			expected = "ABC  "
		)

		padded := padDefault(toPad, SPACE, 5)
		require.Equal(t, expected, padded)
	})

	t.Run("automatic alignment for numbers", func(t *testing.T) {
		t.Run("should pad number left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "123"
				expected = "  123"
			)

			require.True(t, isNumerical(toPad))
			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad decimal number left (point) (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "12.3"
				expected = " 12.3"
			)

			require.True(t, isNumerical(toPad))
			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad decimal number left (comma) (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "12,3"
				expected = " 12,3"
			)

			require.True(t, isNumerical(toPad))
			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad signed number left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "-123"
				expected = " -123"
			)

			require.True(t, isNumerical(toPad))
			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad signed number left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "+123"
				expected = " +123"
			)

			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad number with thousands separator (comma) (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "123,456,789"
				expected = " 123,456,789"
			)

			padded := padDefault(toPad, SPACE, 12)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad number with thousands separator (space) (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "123 456 789"
				expected = " 123 456 789"
			)

			padded := padDefault(toPad, SPACE, 12)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad % left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "2%"
				expected = "   2%"
			)

			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "94.2%"
				expected = " 94.2%"
			)

			padded := padDefault(toPad, SPACE, 6)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "-4.2%"
				expected = " -4.2%"
			)

			padded := padDefault(toPad, SPACE, 6)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad currency amount left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "$123"
				expected = " $123"
			)

			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad negative currency amount left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "$-123"
				expected = " $-123"
			)

			padded := padDefault(toPad, SPACE, 6)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad negative currency amount left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "-$123"
				expected = " -$123"
			)

			padded := padDefault(toPad, SPACE, 6)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad currency amount left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "€123"
				expected = " €123"
			)

			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad currency amount left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "123¥"
				expected = " 123¥"
			)

			padded := padDefault(toPad, SPACE, 5)
			require.Equal(t, expected, padded)
		})

		t.Run("should pad scientific number left (right-aligned)", func(t *testing.T) {
			const (
				toPad    = "1.2e-10"
				expected = " 1.2e-10"
			)

			padded := padDefault(toPad, SPACE, 8)
			require.Equal(t, expected, padded)
		})
	})
}

func TestPadCenter(t *testing.T) {
	t.Parallel()

	t.Run("should center string (odd)", func(t *testing.T) {
		const (
			toPad    = "abc"
			expected = " abc "
		)

		padded := padCenter(toPad, SPACE, 5)
		require.Equal(t, expected, padded)
	})

	t.Run("should center string (even)", func(t *testing.T) {
		const (
			toPad    = "abc"
			expected = " abc  "
		)

		padded := padCenter(toPad, SPACE, 6)
		require.Equal(t, expected, padded)
	})
}
