package tablewrappers

type (
	Option func(*wrapOptions)

	wrapOptions struct {
		strictWidth bool
		splitters   []Splitter
	}
)

func optionsWithDefaults(opts []Option) *wrapOptions {
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
func WithWrapWordSplitters(splitters ...Splitter) Option {
	return func(o *wrapOptions) {
		o.splitters = splitters
	}
}

func WithWrapStrictMaxWidth(enabled bool) Option {
	return func(o *wrapOptions) {
		o.strictWidth = enabled
	}
}
