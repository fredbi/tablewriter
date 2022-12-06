package tablewriter

import (
	"github.com/logrusorgru/aurora/v4"
)

// Formatter is a formatting function from the github.com/logrusorgru/aurora/v4 package.
//
// It wraps some argument with an appropriate ANSI terminal escape sequence.
type Formatter = func(interface{}) aurora.Value

func format(in string, formatter Formatter) string {
	if formatter == nil {
		return in
	}

	return aurora.Sprintf(formatter(in))
}

// WithHeaderFormatters allows to specify ANSI terminal control sequences to format the header.
//
// In particular this may be used to colorize the header.
func WithHeaderFormatters(formatters map[int]Formatter) Option {
	return func(o *options) {
		o.headerParams = formatters
	}
}

// WithFooterFormatters allows to specify ANSI terminal control sequences to format the footer.
//
// In particular this may be used to colorize the footer.
func WithFooterFormatters(formatters map[int]Formatter) Option {
	return func(o *options) {
		o.footerParams = formatters
	}
}

// WithColFormatters allows to specify ANSI terminal control sequences to format cells by columns.
func WithColFormatters(formatters map[int]Formatter) Option {
	return func(o *options) {
		o.columnsParams = formatters
	}
}
