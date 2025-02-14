package tablewriter

import (
	"regexp"
	"strings"

	wrap "github.com/fredbi/tablewriter/tablewrappers"
)

type (
	// HAlignment describes how to horizontally align an element in a cell.
	HAlignment uint8

	padFunc    func(in string, pad string, width int) string
	colAligner func(col int) padFunc
)

// Horizontal alignment
const (
	AlignDefault HAlignment = iota
	AlignCenter
	AlignRight
	AlignLeft
)

var (
	rexNumerical = regexp.MustCompile(`^\s*((\+|-)?\pS)?(\+|-)?((\pN+?)|(\pN{3}[\s,]))+([\.,]\pN*)?(%|\pS|([eE][\+-]{0,1}\pN+))?\s*$`)
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
	gap := width - wrap.DisplayWidth(s)
	if gap <= 0 {
		return s
	}

	gapLeft := gap / 2
	gapRight := gap - gapLeft

	return strings.Repeat(pad, gapLeft) + s + strings.Repeat(pad, gapRight)
}

func padRight(s, pad string, width int) string {
	gap := width - wrap.DisplayWidth(s)
	if gap <= 0 {
		return s
	}

	return s + strings.Repeat(pad, gap)
}

func padLeft(s, pad string, width int) string {
	gap := width - wrap.DisplayWidth(s)
	if gap <= 0 {
		return s
	}

	return strings.Repeat(pad, gap) + s
}

// isNumerical detects numbers, percentages and currency amounts.
func isNumerical(str string) bool {
	return rexNumerical.MatchString(str)
}
