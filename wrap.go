// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

//nolint:unused,unparam,staticcheck
package tablewriter

import (
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

const (
	nl  = "\n"
	sp  = " "
	tab = "\t"
)

const defaultPenalty = 1e5

type (
	Splitter   func(rune) bool
	WrapOption func(*wrapOptions)

	Wrapper struct {
	}

	wrapOptions struct {
		strictWidth bool
		splitters   []Splitter
	}

	columns []column
	cells   []cell

	column struct {
		i        int
		maxWidth int
		cells    cells
	}

	cell struct {
		i       int
		j       int
		content *string
		pvalues []int
		width   int
		passNo  int
	}
)

func (c columns) Less(i, j int) bool {
	return c[i].maxWidth > c[j].maxWidth
}

func (c columns) Len() int {
	return len(c)
}

func (c columns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c cells) Less(i, j int) bool {
	return c[i].width > c[j].width
}

func (c cells) Len() int {
	return len(c)
}

func (c cells) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c column) SortRows() {
	sort.Sort(c.cells)
}

var (
	BlankSplitter = unicode.IsSpace
	PunctSplitter = unicode.IsPunct
	LineSplitter  = func(r rune) bool { return r == '\n' || r == '\r' }
)

func makeReplacer(separators []rune) *strings.Replacer {
	return strings.NewReplacer() // TODO
}

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

func wrapOptionsWithDefaults(opts []WrapOption) *wrapOptions {
	options := &wrapOptions{
		splitters: []Splitter{
			BlankSplitter,
			LineSplitter,
		},
	}

	for _, apply := range opts {
		apply(options)
	}

	return options
}

// WithWrapWordSplitters defines a wrapper's word boundaries split functions.
//
// The default is to break words on IsSpace runes and new-line/carriage return.
func WithWrapWordSplitters(splitters ...Splitter) WrapOption {
	return func(o *wrapOptions) {
		o.splitters = splitters
	}
}

func WithWrapStrictMaxWidth(enabled bool) WrapOption {
	return func(o *wrapOptions) {
		o.strictWidth = enabled
	}
}

// Wrap input string s into a paragraph of lines of length lim, with minimal
// raggedness.
// @deprecated
func WrapString(s string, lim int, opts ...WrapOption) ([]string, int) {
	return wrapString(s, lim, opts...)
}

func wrapString(s string, lim int, opts ...WrapOption) ([]string, int) {
	options := wrapOptionsWithDefaults(opts)

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
		if max > lim {
			lim = max
		}
	}

	for _, line := range WrapWords(words, 1, lim, defaultPenalty) {
		lines = append(lines, strings.Join(line, sp))
	}

	if options.strictWidth {
		// wrap harder -- TODO(fred)
	}

	if len(lines) == 0 {
		lines = []string{""}
	}

	return lines, lim
}

// WrapWords is the low-level line-breaking algorithm, useful if you need more
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
func WrapWords(words []string, spc, lim, pen int) [][]string {
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
		if length[i][n-1] <= lim {
			cost[i] = 0
			nbrk[i] = n
		} else {
			for j := i + 1; j < n; j++ {
				d := lim - length[i][j-1]
				c := d*d + cost[j]
				if length[i][j-1] > lim {
					c += pen // too-long lines get a worse penalty
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

// getLines decomposes a multiline string into a slice of strings.
func getLines(s string) []string {
	return strings.Split(s, nl)
}
