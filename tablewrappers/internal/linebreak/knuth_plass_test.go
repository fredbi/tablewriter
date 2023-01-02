package linebreak

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredbi/tablewriter/tablewrappers/internal/wordbreak/hyphenator"
	"github.com/stretchr/testify/require"
)

func TestLineBreaker(t *testing.T) {
	/*
		t.Run("should left-align a paragraph (width 50)", func(t *testing.T) {
			const (
				paragraph = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, ` +
					`sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. ` +
					`Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ` +
					`Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ` +
					`Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
				display = 50
			)
			t.Run("render", testLeftAlign(paragraph, display))

			t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display-10))
		})

		t.Run("should left-align Knuth's classical example (width 25)", func(t *testing.T) {
			const (
				paragraph = `In olden times when wishing still helped one, there lived a king ` +
					`whose daughters were all beautiful, but the youngest was so beautiful ` +
					`that the sun itself, which has seen so much, was astonished whenever it ` +
					`shone in her face. Close by the king's castle lay a great dark forest, ` +
					`and under an old lime-tree in the forest was a well, and when the day ` +
					`was very warm, the king's child went out into the forest and sat down by ` +
					`the side of the cool fountain, and when she was bored she took a golden ball, ` +
					`and threw it up on high and caught it, and this ball was her favorite plaything.`
				display = 30
			)

			t.Run("render", testLeftAlign(paragraph, display))

			t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display-10))
		})
	*/

	t.Run("should hyphenate", func(t *testing.T) {
		const (
			paragraph = `Honorificabilitudinitatibus is super long`
			display   = 20
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
	})

	/*
		t.Run("should hyphenate", func(t *testing.T) {
			const (
				paragraph = `In olden times when wishing still helped one, there lived a king ` +
					`whose daughters were all beautiful, but the youngest was so beautiful ` +
					`that the sun itself, which has seen so much, was astonished whenever it ` +
					`shone in her face. Close by the king's castle lay a great dark forest, ` +
					`and under an old lime-tree in the forest was a well, and when the day ` +
					`was very warm, the king's child went out into the forest and sat down by ` +
					`the side of the cool fountain, and when she was bored she took a golden ball, ` +
					`and threw it up on high and caught it, and this ball was her favorite plaything.`
				display = 8
			)

			t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
		})
	*/
}

func testLeftAlign(paragraph string, display int) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)

		lb := New() // with all defaults
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}

func testRenderLines(lines []string, display int) {
	for _, line := range lines {
		pad := strings.Repeat(".", max(0, display-len(line)))
		fmt.Printf("|%s%s|\n", line, pad)
	}
}

func testLeftAlignHyphenize(paragraph string, display int) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)

		h := hyphenator.New()
		lb := New(
			WithRenderHyphens(true),
			WithHyphenator(h.BreakWord),
			WithHyphenPenalty(5.00),
		)
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}
