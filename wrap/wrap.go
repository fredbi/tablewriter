// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

//nolint:unused,unparam,staticcheck
package wrap

import (
	"math"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	nl             = "\n"
	sp             = " "
	tab            = "\t"
	defaultPenalty = 1e5
)

var ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")

type (
	DefaultWrapper struct {
		*wrapOptions
	}
)

// New builds a new default wrapper.
func New(opts ...Option) *DefaultWrapper {
	return &DefaultWrapper{
		wrapOptions: optionsWithDefaults(opts),
	}
}

// Wrap input string s into a paragraph of lines of limited length, with minimal raggedness.
func (w *DefaultWrapper) WrapString(s string, limit int) []string {
	options := w.wrapOptions

	// there are 2 levels of splitting:
	// * irrelevant boundaries such as blank space or new lines
	// * word boundaries that we want to keep, such as punctuation marks
	splitter := composeSplitters(options.splitters)
	words := strings.FieldsFunc(s, splitter)

	var lines []string
	max := 0
	for _, v := range words {
		if len(v) == 0 {
			continue
		}

		max = runewidth.StringWidth(v)
		if max > limit {
			limit = max
		}
	}

	for _, line := range wrapWords(words, 1, limit, defaultPenalty) {
		lines = append(lines, strings.Join(line, sp))
	}

	if options.strictWidth {
		// wrap harder -- TODO(fred)
	}

	if len(lines) == 0 {
		// always ensure at least one line
		lines = []string{""}
	}

	return lines
}

// wrapWords is the low-level line-breaking algorithm, useful if you need more
// control over the details of the text wrapping process. For most uses,
// WrapString will be sufficient and more convenient.
//
// WrapWords splits a list of words into lines with minimal "raggedness",
// treating each rune as one unit, accounting for spc units between adjacent
// words on each line, and attempting to limit lines to lim units. Raggedness
// is the total error over all lines, where error is the square of the
// difference of the length of the line and lim. Too-long lines (which only
// happen when a single word is longer than lim units) have pen penalty units
// added to the error.
func wrapWords(words []string, spc, limit, penalty int) [][]string {
	n := len(words)

	length := make([][]int, n)
	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = runewidth.StringWidth(words[i])
		for j := i + 1; j < n; j++ {
			length[i][j] = length[i][j-1] + spc + runewidth.StringWidth(words[j])
		}
	}

	nbrk := make([]int, n)
	cost := make([]int, n)
	for i := range cost {
		cost[i] = math.MaxInt32
	}

	for i := n - 1; i >= 0; i-- {
		if length[i][n-1] <= limit {
			cost[i] = 0
			nbrk[i] = n
		} else {
			for j := i + 1; j < n; j++ {
				d := limit - length[i][j-1]
				c := d*d + cost[j]
				if length[i][j-1] > limit {
					c += penalty // too-long lines get a worse penalty
				}
				if c < cost[i] {
					cost[i] = c
					nbrk[i] = j
				}
			}
		}
	}

	var lines [][]string
	i := 0
	for i < n {
		lines = append(lines, words[i:nbrk[i]])
		i = nbrk[i]
	}

	return lines
}

// DisplayWidth yields the size of a string when rendered on a terminal.
//
// ANSI escape sequences are discared.
func DisplayWidth(str string) int {
	return displayWidth(str)
}

func displayWidth(str string) int {
	return runewidth.StringWidth(ansi.ReplaceAllLiteralString(str, ""))
}
