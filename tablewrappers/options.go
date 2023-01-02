package tablewrappers

type (
	Option func(*wrapOptions)

	wrapOptions struct {
		strictWidth bool // TODO for StringWrapper
		splitters   []Splitter
		minColWidth map[int]int
		maxColWidth map[int]int
	}
)

func optionsWithDefaults(opts []Option) *wrapOptions {
	options := &wrapOptions{
		splitters: []Splitter{
			BlankSplitter,
			LineSplitter,
		},
		minColWidth: make(map[int]int),
		maxColWidth: make(map[int]int),
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

func WithColMinWidth(column int, width int) Option {
	return func(o *wrapOptions) {
		o.minColWidth[column] = width
	}
}

func WithColMaxWidth(column int, width int) Option {
	return func(o *wrapOptions) {
		o.maxColWidth[column] = width
	}
}

// TODO: NoColWordBreak(cols ...int) Option
