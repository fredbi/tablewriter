package tablewrappers

import (
	"strings"
)

type (
	DefaultWrapper struct {
		*wrapOptions
		splitter Splitter
	}
)

// New builds a new default wrapper.
func NewDefault(opts ...Option) *DefaultWrapper {
	w := &DefaultWrapper{
		wrapOptions: optionsWithDefaults(opts),
	}
	w.splitter = composeSplitters(w.splitters)

	return w
}

// Wrap input string s into a paragraph of lines of limited length, with minimal raggedness.
func (w *DefaultWrapper) WrapString(s string, limit int) []string {
	words := strings.FieldsFunc(s, w.splitter) // default: splits according to blanks & lines
	limit = max(limit, cellWidth(words))       // readjust limit to maximum width of a single word
	lines := wrapMultiline(words, limit, 1)

	/*
		if w.strictWidth {
			// wrap harder -- TODO(fred)
		}
	*/

	return lines
}
