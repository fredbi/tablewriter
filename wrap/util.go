package wrap

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

func displayWidth(str string) int {
	return runewidth.StringWidth(ansi.ReplaceAllLiteralString(str, ""))
}
