package tablewrappers

import (
	"strings"
)

type (
	DefaultWrapper struct {
		*wrapOptions
	}
)

// New builds a new default wrapper.
func NewDefault(opts ...Option) *DefaultWrapper {
	return &DefaultWrapper{
		wrapOptions: optionsWithDefaults(opts),
	}
}

// Wrap input string s into a paragraph of lines of limited length, with minimal raggedness.
func (w *DefaultWrapper) WrapString(s string, limit int) []string {
	options := w.wrapOptions

	splitter := composeSplitters(options.splitters)
	words := strings.FieldsFunc(s, splitter) // default: splits according to blanks & lines
	limit = max(limit, cellWidth(words))     // readjust limit to maximum width of a single word
	lines := wrapMultiline(words, limit)

	if options.strictWidth {
		// wrap harder -- TODO(fred)
	}

	return lines
}
