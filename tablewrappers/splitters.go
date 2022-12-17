package tablewrappers

import (
	"unicode"
)

type (
	Splitter func(rune) bool
)

var (
	BlankSplitter = unicode.IsSpace
	PunctSplitter = unicode.IsPunct
	LineSplitter  = func(r rune) bool { return r == '\n' || r == '\r' }
	WordBreaker   = func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }
)

func composeSplitters(splitters []Splitter) Splitter {
	return func(r rune) bool {
		for _, fn := range splitters {
			if fn(r) {
				return true
			}
		}

		return false
	}
}
