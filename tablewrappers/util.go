package tablewrappers

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

var ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")

// DisplayWidth yields the size of a string when rendered on a terminal.
//
// ANSI escape sequences are discared.
func DisplayWidth(str string) int {
	return displayWidth(str)
}

// how about unicode.IsControl()?
func displayWidth(str string) int {
	return runewidth.StringWidth(ansi.ReplaceAllLiteralString(str, ""))
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

/*
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
*/

// CellWidth determines the displayed width of a multi-lines cell.
func CellWidth(lines []string) int {
	return cellWidth(lines)
}

func cellWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if w := displayWidth(line); w > maxWidth {
			maxWidth = w
		}
	}

	return maxWidth
}

func cellsMaxWidth(rows [][]string) int {
	maxWidth := 0
	for _, row := range rows {
		if w := cellWidth(row); w > maxWidth {
			maxWidth = w
		}
	}

	return maxWidth
}
