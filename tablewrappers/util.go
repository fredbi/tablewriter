package tablewrappers

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

var (
	rexANSI = regexp.MustCompile(
		"\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]",
	)

	rexStripANSI = regexp.MustCompile(
		"^((?:\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K])*?)?([^\033]*)*?((?:\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K])*)?$",
	)
)

// DisplayWidth yields the size of a string when rendered on a terminal.
//
// ANSI escape sequences are discarded from this count.
func DisplayWidth(str string) int {
	return displayWidth(str)
}

// CellWidth determines the displayed width of a multi-lines cell.
func CellWidth(lines []string) int {
	return cellWidth(lines)
}

// displayWidth yields the display size of string on a terminal.
func displayWidth(str string) int {
	return runewidth.StringWidth(rexANSI.ReplaceAllLiteralString(str, ""))
}

// cellWidth returns the display width of a multi-line cell.
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

// cellsMaxWidth return the maxium width of multi-line column content in different rows.
func cellsMaxWidth(rows [][]string) int {
	maxWidth := 0
	for _, row := range rows {
		if w := cellWidth(row); w > maxWidth {
			maxWidth = w
		}
	}

	return maxWidth
}

func stripANSI(str string) (string, string, string) {
	matches := rexStripANSI.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return str, "", ""
	}
	groups := matches[0] // TODO: what if several ansi wrapped around a single word???
	if len(groups) < 2 {
		return str, "", ""
	}

	switch len(groups) {
	case 2:
		return "", groups[1], ""
	case 3:
		return groups[2], groups[1], ""
	default:
		return groups[2], groups[1], groups[3]
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
