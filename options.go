package tablewriter

import (
	"io"
	"os"
)

const (
	// MaxColWidth is the default maximum width of a column.
	MaxColWidth = 30
)

// Default separator characters.
const (
	CENTER    = "+"
	ROW       = "-"
	COLUMN    = "|"
	SPACE     = " "
	NEWLINE   = "\n"
	NOPADDING = ""
)

type (
	// Border represent a borders specification for a table.
	Border struct {
		Left   bool
		Right  bool
		Top    bool
		Bottom bool
	}

	// Option to render a table
	Option func(*options)

	separatorOptions struct {
		pCenter string
		pRow    string
		pColumn string
		newLine string
	}

	formatOptions struct {
		headerParams  map[int]Formatter
		columnsParams map[int]Formatter
		footerParams  map[int]Formatter
	}

	alignOptions struct {
		headerAlign    HAlignment
		footerAlign    HAlignment
		cellAlign      HAlignment
		perColumnAlign map[int]HAlignment
	}

	options struct {
		rows        [][]string // input rows
		header      []string
		footer      []string
		captionText string

		// rendering target
		out io.Writer

		// width & height
		cs map[int]int // min width for a column
		ms map[int]int // max width for a column
		rs map[int]int // max lines per cell

		// header title-case
		autoFmt bool

		separatorOptions

		// borders
		borders Border

		// cell formatting
		// tRow           int
		// tColumn int
		autoWrap       bool
		reflowText     bool
		mW             int
		autoMergeCells bool
		noWhiteSpace   bool
		tablePadding   string

		// line separators
		separatorBetweenRows bool
		separatorAfterHeader bool
		separatorAfterFooter bool

		// formatting
		formatOptions

		// horizontal alignment
		alignOptions
	}
)

func defaultOptions(opts []Option) *options {
	o := &options{
		out:                  os.Stdout,
		rows:                 [][]string{},
		cs:                   make(map[int]int),
		ms:                   make(map[int]int),
		rs:                   make(map[int]int),
		captionText:          "",
		autoFmt:              true,
		autoWrap:             true,
		reflowText:           true,
		mW:                   MaxColWidth,
		separatorOptions:     defaultSeparatorOptions(),
		alignOptions:         defaultAlignOptions(),
		formatOptions:        defaultFormatOptions(),
		separatorAfterHeader: true,
		separatorAfterFooter: true,
		borders:              Border{Left: true, Right: true, Bottom: true, Top: true},
		tablePadding:         SPACE,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func defaultSeparatorOptions() separatorOptions {
	return separatorOptions{
		pCenter: CENTER,
		pRow:    ROW,
		pColumn: COLUMN,
		newLine: NEWLINE,
	}
}

func defaultAlignOptions() alignOptions {
	return alignOptions{
		headerAlign:    AlignDefault,
		footerAlign:    AlignDefault,
		cellAlign:      AlignDefault,
		perColumnAlign: make(map[int]HAlignment),
	}
}

func defaultFormatOptions() formatOptions {
	return formatOptions{
		headerParams:  make(map[int]Formatter),
		columnsParams: make(map[int]Formatter),
		footerParams:  make(map[int]Formatter),
	}
}

// WithWriter specifies the output writer to render this table.
//
// The default is os.Stdout.
func WithWriter(writer io.Writer) Option {
	return func(o *options) {
		o.out = writer
	}
}

// WithRows specifies the rows of the table, each being a record of columns.
//
// The input is not required to contain the same number of columns for each row.
func WithRows(rows [][]string) Option {
	return func(o *options) {
		o.rows = rows
	}
}

// WithHeader specifies the header fields for this table.
func WithHeader(header []string) Option {
	return func(o *options) {
		o.header = header
	}
}

// WithFooter specifies the footer fields for this table.
func WithFooter(footer []string) Option {
	return func(o *options) {
		o.footer = footer
	}
}

// WithTitledHeader autoformats headers with Title-case.
func WithTitledHeader(enabled bool) Option {
	return func(o *options) {
		o.autoFmt = enabled
	}
}

// WithCaption displays a caption under the table.
func WithCaption(caption string) Option {
	return func(o *options) {
		o.captionText = caption
	}
}

// WithWrap enables content wrapping inside columns to abide
// by column width constraints.
//
// Wrapping is enabled by default (the default maximum column width is 30 characters).
//
// Some WrapOptions may be passed to further tune the behavior of the wrapper.
func WithWrap(enabled bool, opts ...WrapOption) Option {
	// TODO(fred): wrap options
	return func(o *options) {
		o.autoWrap = enabled
	}
}

// WithCenterSeparator defines the string used to represent intersections of
// the table grid.
//
// The default is '+'.
func WithCenterSeparator(sep string) Option {
	return func(o *options) {
		o.pCenter = sep
	}
}

// WithRowSeparator defines the string used to separate rows.
//
// The default is '-'.
func WithRowSeparator(sep string) Option {
	return func(o *options) {
		o.pRow = sep
	}
}

// WithAllBorders enables (resp. disables) all table borders.
//
// Borders are enabled by default.
func WithAllBorders(enabled bool) Option {
	return func(o *options) {
		o.borders = Border{enabled, enabled, enabled, enabled}
	}
}

// WithBorders allows for a detailed specification of which borders are rendered.
func WithBorders(border Border) Option {
	return func(o *options) {
		o.borders = border
	}
}

// WithRowLine indicates that each row is followed by a separation line.
//
// By default, rows are packed without line separator.
func WithRowLine(enabled bool) Option {
	return func(o *options) {
		o.separatorBetweenRows = enabled
	}
}

// WithNewLine defines the end of line character.
//
// The default is '\n'.
func WithNewLine(nl string) Option {
	return func(o *options) {
		o.newLine = nl
	}
}

// WithHeaderLine prints a separation line under the header.
//
// This is enabled by default.
func WithHeaderLine(enabled bool) Option {
	return func(o *options) {
		o.separatorAfterHeader = enabled
	}
}

// WithFooterLine prints a separation line under the footer.
//
// This is enabled by default.
func WithFooterLine(enabled bool) Option {
	return func(o *options) {
		o.separatorAfterFooter = enabled
	}
}

// No White Space. TODO: more tests and bug fixes.
func WithNoWhiteSpace(enabled bool) Option {
	return func(o *options) {
		o.noWhiteSpace = enabled
	}
}

// WithFooterAlignment defines the alignment for all footer fields.
//
// The default is CENTER.
func WithFooterAlignment(footerAlign HAlignment) Option {
	return func(o *options) {
		o.footerAlign = footerAlign
	}
}

// WithHeaderAlignment defines the alignment for all headings.
//
// The default is CENTER.
func WithHeaderAlignment(align HAlignment) Option {
	return func(o *options) {
		o.headerAlign = align
	}
}

// WithPadding defines the padding character inside the table.
//
// The default is a blank space.
func WithPadding(padding string) Option {
	return func(o *options) {
		o.tablePadding = padding
	}
}

// WithColumnSeparator defines the character to separate columns.
//
// The default is '|'.
func WithColumnSeparator(sep string) Option {
	return func(o *options) {
		o.pColumn = sep
	}
}

// WithColWidth defines the maximum width for all columns (in characters).
//
// The default is 30.
func WithColWidth(width int) Option {
	return func(o *options) {
		o.mW = width
	}
}

// WithColMaxWidth defines the maximum width for a specific column.
//
// This overrides the setting defined by WithColWidth.
func WithColMaxWidth(column int, width int) Option {
	return func(o *options) {
		o.ms[column] = width
	}
}

// WithColMaxWidths defines the maximum width for a set of columns.
func WithColMaxWidths(maxWidths map[int]int) Option {
	return func(o *options) {
		for k, v := range maxWidths {
			o.ms[k] = v
		}
	}
}

// WithMarkdown reproduces classifical markdown tables.
//
// This option is a shortcut to:
//
//	WithCenterSeparator("|")
//	WithBorders(Border{Left: true, Top: false, Right: true, Bottom: false})
func WithMarkdown(enabled bool) Option {
	return func(o *options) {
		o.borders = Border{Left: true, Top: false, Right: true, Bottom: false}
		o.pCenter = "|"
	}
}

func WithWrapReflow(enabled bool) Option {
	return func(o *options) {
		o.reflowText = enabled
	}
}

// WithCellAlignment defines the default alignment for row cells.
//
// The default is CENTER for strings, RIGHT for numbers (and %).
func WithCellAlignment(align HAlignment) Option {
	return func(o *options) {
		o.cellAlign = align
	}
}

// WithMergeCells enables the merging of adjacent cells with the same value.
func WithMergeCells(enabled bool) Option {
	return func(o *options) {
		o.autoMergeCells = enabled
	}
}

// WithColMinWidth specifies the minimum width of columns. TODO: testing. Does this work?
func WithColMinWidth(column int, width int) Option {
	return func(o *options) {
		o.cs[column] = width
	}
}

// WithColAlignment defines the aligment for a set of columns.
func WithColAlignment(align map[int]HAlignment) Option {
	return func(o *options) {
		o.perColumnAlign = align
	}
}
