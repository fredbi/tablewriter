package tablewriter

import (
	"regexp"
	"strings"
)

type (
	// HAlignment describes how to horizontally align an element in a cell.
	HAlignment uint8

	padFunc func(in string, pad string, width int) string
)

// Horizontal alignment
const (
	AlignDefault HAlignment = iota
	AlignCenter
	AlignRight
	AlignLeft
)

var (
	decimal = regexp.MustCompile(`^-?(?:\d{1,3}(?:,\d{3})*|\d+)(?:\.\d+)?$`)
	percent = regexp.MustCompile(`^-?\d+\.?\d*$%$`)
)

// padder yields the appropriate padding function for the alignment type.
func (h HAlignment) padder() padFunc {
	switch h {
	case AlignLeft:
		return padRight
	case AlignRight:
		return padLeft
	case AlignCenter:
		return padCenter
	case AlignDefault:
		fallthrough
	default:
		return padDefault
	}
}

// padDefault pads numerical values to the left (right-aligned) and other values to the right (left-aligned).
func padDefault(s, pad string, width int) string {
	if isNumerical(s) {
		return padLeft(s, pad, width)
	}

	return padRight(s, pad, width)
}

// padCenter centers a string
func padCenter(s, pad string, width int) string {
	gap := width - displayWidth(s)
	if gap <= 0 {
		return s
	}

	gapLeft := gap / 2
	gapRight := gap - gapLeft

	return strings.Repeat(pad, gapLeft) + s + strings.Repeat(pad, gapRight)
}

func padRight(s, pad string, width int) string {
	gap := width - displayWidth(s)
	if gap <= 0 {
		return s
	}

	return s + strings.Repeat(pad, gap)
}

func padLeft(s, pad string, width int) string {
	gap := width - displayWidth(s)
	if gap <= 0 {
		return s
	}

	return strings.Repeat(pad, gap) + s
}

func isNumerical(str string) bool {
	return decimal.MatchString(strings.TrimSpace(str)) || percent.MatchString(strings.TrimSpace(str))
}
