package tablewrappers

import (
	"sort"
	"strings"
)

type (
	breakLevel int

	// words represent a sentence made of words.
	// This collection knows how to be sorted with the widest word first and
	// by its natural order in the original sentence.
	words []*word

	// word represent the n-th word in a sentence, possibly split in parts.
	word struct {
		n     int
		parts []string
	}
)

const (
	breakNone breakLevel = iota
	breakOnSeps
	// TODO: breakOnHyphenation breakLevel
	breakAnywhere
)

// newWords builds a new collection of words for a split sentence.
func newWords(sentence []string) words {
	out := make(words, 0, len(sentence))
	for i, w := range sentence {
		out = append(out, &word{
			n:     i,
			parts: []string{w},
		})
	}

	return out
}

func (w words) Len() int {
	return len(w)
}

func (w words) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w words) Less(i, j int) bool {
	return cellWidth(w[i].parts) > cellWidth(w[j].parts)
}

func (w words) Sort() {
	sort.Stable(w)
}

func (w words) SortNatural() {
	sort.Slice(w, func(i, j int) bool {
		return w[i].n < w[j].n
	})
}

// Width yields the maximum width on display of a collection of words,
// reassembled as space-separated sentence, but with broken parts on separate lines.
func (w words) Width() int {
	if len(w) == 0 {
		return 0
	}

	total := w[0].Width()
	for _, word := range w[1:] {
		total += word.Width() + 1 // +1 space between words
	}

	return total
}

// Width yields the width on display of a word, possibly broken over multiple lines.
func (w *word) Width() int {
	return cellWidth(w.parts)
}

// Break a word given the width limit and the aggressiveness of the word-breaker.
func (w *word) Break(limit int, aggressiveness breakLevel) {
	if w.Width() <= limit {
		return
	}

	newParts := make([]string, 0, len(w.parts))

	for _, part := range w.parts {
		if displayWidth(part) <= limit {
			newParts = append(newParts, part)

			continue
		}

		newParts = append(newParts, breakWord(part, limit, aggressiveness)...)
	}

	w.parts = newParts
}

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
	case breakNone:
		return []string{word}
	case breakOnSeps:
		return wordBreaker(WordBreaker)(word, limit)
	default:
		return wordBreaker(func(r rune) bool { return true })(word, limit)
	}
}

func wordBreaker(splitter Splitter) func(string, int) []string {
	return func(word string, limit int) []string {
		parts := breakAtFunc(word, splitter)
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
}

// breakAtFunc works like strings.FieldsFunc, but retain separators.
//
// Break always happen _after_ the separator.
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
