package tablewrappers

import (
	"strings"
)

type breakLevel int

const (
	breakOnBoundaries breakLevel = iota
	// TODO: breakOnHyphenation breakLevel
	breakAnywhere
)

// breakWord breaks a word with some aggressiveness into parts smaller than the given limit.
//
// A word is assumed not to contain any blank space.
//
// With aggressiveness 0, words are broken along natural separators such as
// anything that is not a letter or a digit.
//
// With aggressiveness 1, words are broken anywhere.
func breakWord(word string, limit int, aggressiveness breakLevel) []string {
	switch aggressiveness {
	case breakOnBoundaries:
		return breakWordOnBoundaries(word, limit)
	default:
		return breakWordAnywhere(word, limit)
	}
}

// breakWordOnBoundarie will break a word along any non-letter or non-digit symbol.
func breakWordOnBoundaries(word string, limit int) []string {
	parts := breakAtFunc(word, WordBreaker)
	lines := make([]string, 0, len(parts))

	for _, part := range wrapWords(parts, 0, limit, defaultPenalty) {
		lines = append(lines, strings.Join(part, ""))
	}

	if len(lines) == 0 {
		// always ensure at least one line
		lines = []string{""}
	}

	return lines
}

func breakWordAnywhere(word string, limit int) []string {
	parts := breakAtFunc(word, func(r rune) bool { return true })
	lines := make([]string, 0, len(parts))

	for _, part := range wrapWords(parts, 0, limit, defaultPenalty) {
		lines = append(lines, strings.Join(part, ""))
	}

	if len(lines) == 0 {
		// always ensure at least one line
		lines = []string{""}
	}

	return lines
}

func breakAtFunc(word string, isBreak Splitter) []string {
	parts := make([]string, 0, len(word))
	previous := 0

	for i, r := range word {
		if isBreak(r) {
			parts = append(parts, word[previous:i+1])
			previous = i + 1
		}
	}

	if previous < len(word) {
		parts = append(parts, word[previous:])
	}

	return parts
}
