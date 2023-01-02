package linebreak

import (
	"github.com/fredbi/tablewriter/tablewrappers/internal/wordbreak"
)

type (
	Option func(*options)

	options struct {
		tolerance float64
		badness   float64
		demerits  demeritsT
		formatterOptions
	}

	formatterOptions struct {
		measurer    func(string) int // TODO: func([]rune) int
		scaleFactor float64          // the scale factor is a multiplier to adapt measures to better fit the algorithm's settings
		space       sums             // space widths

		wordBreak     bool                  // enable breaking words (hyphenations, ...)
		renderHyphens bool                  // enable the rendering of hyphens for hyphenated words
		hyphenPenalty float64               // penalty to give to hyphenated words
		hyphenator    wordbreaker.SplitFunc // word breaker for hyphenation TODO func([]rune) [][]rune ?
		minHyphenate  int                   // minimum length of a token for hyphenation to apply
		glueStretch   int
	}
)

func WithTolerance(tolerance float64) Option {
	return func(o *options) {
		o.tolerance = tolerance
	}
}

func WithScaleFactor(scale float64) Option {
	return func(o *options) {
		o.scaleFactor = scale
	}
}

func WithWordBreak(enabled bool) Option {
	return func(o *options) {
		o.wordBreak = enabled
	}
}

func WithMeasurer(measurer func(string) int) Option {
	return func(o *options) {
		o.measurer = measurer
	}
}

func WithHyphenator(hyphenator wordbreaker.SplitFunc) Option {
	return func(o *options) {
		o.wordBreak = true
		o.hyphenator = hyphenator
	}
}

func WithHyphenPenalty(penalty float64) Option {
	return func(o *options) {
		o.hyphenPenalty = penalty
		o.demerits.flagged = penalty
	}
}

func WithRenderHyphens(enabled bool) Option {
	return func(o *options) {
		o.renderHyphens = enabled
	}
}

func defaultOptions(opts []Option) *options {
	o := &options{
		tolerance:        8.6,
		badness:          100.00,
		demerits:         defaultDemerits(),
		formatterOptions: defaultFormatterOptions(),
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func defaultDemerits() demeritsT {
	return demeritsT{
		line:    10,
		flagged: 100,
		fitness: 3000,
	}
}

func defaultFormatterOptions() formatterOptions {
	return formatterOptions{
		renderHyphens: true,
		measurer:      func(in string) int { return len(in) }, // TODO: use rune width and strip ANSI escape seq
		scaleFactor:   3,                                      //3,
		space: sums{
			width:   3,
			stretch: 6,
			shrink:  9,
		},
		hyphenPenalty: 100,
		hyphenator:    func(in string) []string { return []string{in} },
		minHyphenate:  4, // minimum length of a word to be hyphenated // TODO: should be handled by the hyphenator
		glueStretch:   12,
	}
}
