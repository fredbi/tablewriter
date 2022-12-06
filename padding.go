package tablewriter

import (
	"strings"
)

type (
	HAlignment uint8
	padFunc    func(in string, pad string, width int) string
)

// Horizontal alignment
const (
	AlignDefault HAlignment = iota
	AlignCenter
	AlignRight
	AlignLeft
)

// Return the appropriate padding function for the alignment type selected.
func pad(align HAlignment) padFunc {
	switch align {
	case AlignLeft:
		return padRight
	case AlignRight:
		return padLeft
	case AlignDefault:
		fallthrough
	default:
		return padCenter
	}
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
