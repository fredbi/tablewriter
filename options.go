package tablewriter

import (
	"io"
	"os"

	"github.com/fredbi/tablewriter/titlers"
	"github.com/fredbi/tablewriter/wrap"
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

	alignOptions struct {
		headerAlign    HAlignment
		footerAlign    HAlignment
		cellAlign      HAlignment
		perColumnAlign map[int]HAlignment
	}

	wrapOptions struct {
		wrapper Wrapper
	}

	options struct {
		rows        [][]string // input rows
		header      []string
		footer      []string
		captionText string

		// rendering target
		out io.Writer

		// width & height
		colWidth    map[int]int // min width for a column
		colMaxWidth map[int]int // max width for a column
		maxColWidth int

		// header title-case
		titler transformer

		separatorOptions

		// borders
		borders Border

		wrapOptions

		// cell formatting
		reflowText     bool
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
		colWidth:             make(map[int]int),
		colMaxWidth:          make(map[int]int),
		captionText:          "",
		reflowText:           true,
		maxColWidth:          MaxColWidth,
		wrapOptions:          defaultWrapOptions(),
		separatorOptions:     defaultSeparatorOptions(),
		alignOptions:         defaultAlignOptions(),
		formatOptions:        defaultFormatOptions(),
		separatorAfterHeader: true,
		separatorAfterFooter: true,
		borders:              Border{Left: true, Right: true, Bottom: true, Top: true},
		tablePadding:         SPACE,
		titler:               titlers.DefaultTitler,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func defaultWrapOptions() wrapOptions {
	return wrapOptions{
		wrapper: wrap.New(),
	}
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
		headerAlign:    AlignCenter,
		footerAlign:    AlignCenter,
		cellAlign:      AlignDefault,
		perColumnAlign: make(map[int]HAlignment),
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

// WithTitledHeader autoformats headers and footer as titles.
//
// By default, the title string is trimmed, uppercased. Underscores are replaced by blank spaces.
func WithTitledHeader(enabled bool) Option {
	return func(o *options) {
		if enabled {
			o.titler = titlers.DefaultTitler
		} else {
			o.titler = nil
		}
	}
}

// WithCustomTitler injects a transform function to apply to header and footer values.
func WithCustomTitler(titler func(string) string) Option {
	return func(o *options) {
		o.titler = titler
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
// Whenever enabled, the default wrapper is used. The default wrapper wraps cells into
// multiline content, based on their column maximum width, wrapping only on word boundaries.
func WithWrap(enabled bool) Option {
	return func(o *options) {
		if enabled {
			o.wrapper = wrap.New()
		} else {
			o.wrapper = nil
		}
	}
}

// WithWrapper allows to plug-in a customized cell content Wrapper.
func WithWrapper(wrapper Wrapper) Option {
	return func(o *options) {
		o.wrapper = wrapper
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
		o.maxColWidth = width
	}
}

// WithColMaxWidth defines the maximum width for a specific column.
//
// This overrides the setting defined by WithColWidth.
func WithColMaxWidth(column int, width int) Option {
	return func(o *options) {
		o.colMaxWidth[column] = width
	}
}

// WithColMaxWidths defines the maximum width for a set of columns.
func WithColMaxWidths(maxWidths map[int]int) Option {
	return func(o *options) {
		for k, v := range maxWidths {
			o.colMaxWidth[k] = v
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
		o.colWidth[column] = width
	}
}

// WithColAlignment defines the aligment for a set of columns.
func WithColAlignment(align map[int]HAlignment) Option {
	return func(o *options) {
		o.perColumnAlign = align
	}
}
