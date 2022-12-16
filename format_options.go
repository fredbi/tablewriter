package tablewriter

import (
	"github.com/logrusorgru/aurora/v4"
)

type formatOptions struct {
	headerParams  map[int]Formatter
	columnsParams map[int]Formatter
	footerParams  map[int]Formatter
	captionParams Formatter
}

func defaultFormatOptions() formatOptions {
	return formatOptions{
		headerParams:  make(map[int]Formatter),
		columnsParams: make(map[int]Formatter),
		footerParams:  make(map[int]Formatter),
	}
}

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

// WithCaptionFormatter allows to specify ANSI terminal control sequences to format the table caption.
func WithCaptionFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.captionParams = formatter
	}
}
