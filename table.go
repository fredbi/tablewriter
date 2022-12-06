// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

// Create & Generate text based table
package tablewriter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	headerRowIdx = -1
	footerRowIdx = -2
)

type Table struct {
	*options

	// internal multi-line representations of rows
	lines                   [][][]string // multi-line cells
	headers                 [][]string
	footers                 [][]string
	numColumns              int
	columnsToAutoMergeCells map[int]bool
	columnsAlign            []HAlignment
}

// New builds a new empty table writer.
func New(opts ...Option) *Table {
	t := &Table{
		lines:      [][][]string{},
		headers:    [][]string{},
		footers:    [][]string{},
		numColumns: -1,

		options: defaultOptions(opts),
	}

	return t
}

// NewBuffered builds a new empty table writer that writes in a new bytes.Buffer.
func NewBuffered(opts ...Option) (*Table, *bytes.Buffer) {
	buf := &bytes.Buffer{}

	return New(append(opts, WithWriter(buf))...), buf
}

// Append a row to the table.
func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

// Rows of this table.
func (t *Table) Rows() [][]string {
	return t.rows
}

// Header of this table.
func (t *Table) Header() []string {
	return t.header
}

// Footer of this table.
func (t *Table) Footer() []string {
	return t.footer
}

// Render the table
func (t *Table) Render() {
	t.prepare()

	if t.borders.Top {
		t.printSepLine(true)
	}

	t.printHeader()

	if t.autoMergeCells {
		t.printRowsMergeCells()
	} else {
		t.printRows()
	}

	if !t.separatorBetweenRows && t.borders.Bottom {
		t.printSepLine(true)
	}

	t.printFooter()

	if len(t.captionText) > 0 {
		t.printCaption()
	}
}

func (t *Table) prepare() {
	t.setNumColumns()
	t.fillAlignments()

	for i, v := range t.header {
		lines := t.parseDimension(v, i, headerRowIdx)
		t.headers = append(t.headers, lines)
	}

	for i, cells := range t.rows {
		var rowLines [][]string
		for j, v := range cells {
			rowLines = append(rowLines, t.parseDimension(v, j, i))
		}
		t.lines = append(t.lines, rowLines)
	}

	for i, v := range t.footer {
		lines := t.parseDimension(v, i, footerRowIdx)
		t.footers = append(t.footers, lines)
	}
}

func (t *Table) fillAlignments() {
	num := t.numColumns
	t.columnsAlign = make([]HAlignment, 0, num)

	for i := 0; i < num; i++ {
		alignment, ok := t.perColumnAlign[i]
		if !ok {
			alignment = t.cellAlign
		}

		t.columnsAlign = append(t.columnsAlign, alignment)
	}
}

// setNumColumns determines the number of columns for this table, aligned to the row
// (or header, or footer) with the largest number of columns.
func (t *Table) setNumColumns() {
	t.numColumns = len(t.header)

	for _, row := range t.rows {
		if cols := len(row); cols > t.numColumns {
			t.numColumns = cols
		}
	}

	if cols := len(t.footer); cols > t.numColumns {
		t.numColumns = cols
	}

	// normalize all content to the same number of columns, adding trailing empty columns
	if len(t.header) > 0 {
		if missing := t.numColumns - len(t.header); missing > 0 {
			padded := make([]string, missing)
			t.header = append(t.header, padded...)
		}
	}

	for i, row := range t.rows {
		if missing := t.numColumns - len(row); missing > 0 {
			padded := make([]string, missing)
			t.rows[i] = append(t.rows[i], padded...)
		}
	}

	if len(t.footer) > 0 {
		if missing := t.numColumns - len(t.footer); missing > 0 {
			padded := make([]string, missing)
			t.footer = append(t.footer, padded...)
		}
	}
}

// center determines the character to print at an intersection line, based on the position and border.
func (t *Table) center(i int) string {
	switch {
	case i == -1 && !t.borders.Left:
		fallthrough // -
	case i == len(t.cs)-1 && !t.borders.Right:
		return t.pRow // -
	default:
		return t.pCenter // +
	}
}

// PrintSepLine prints a separation line based on the row width.
func (t *Table) printSepLine(withNewLine bool) {
	fmt.Fprint(t.out, t.center(-1)) // -

	for i := 0; i < len(t.cs); i++ {
		colWidth := t.cs[i]
		fmt.Fprintf(t.out, "%s%s%s%s",
			t.pRow,                           // -
			strings.Repeat(t.pRow, colWidth), // -...-
			t.pRow,                           // -
			t.center(i))                      // +|-
	}

	if withNewLine {
		fmt.Fprint(t.out, t.newLine)
	}
}

// Prints a separation line based on row width with or without cell separator.
//
// This is used when rendering merged cells.
//
// TODO(fred): can't this be factorized with printSepLine??
func (t *Table) printLineOptionalCellSeparators(withNewLine bool, displayCellSeparator []bool) {
	fmt.Fprint(t.out, t.pCenter) // ?? + never -???

	for i := 0; i < len(t.cs); i++ {
		colWidth := t.cs[i]
		if i > len(displayCellSeparator) || displayCellSeparator[i] {
			// display the cell separator
			fmt.Fprintf(t.out, "%s%s%s%s",
				t.pRow,
				strings.Repeat(t.pRow, colWidth),
				t.pRow,
				t.pCenter, // t.center(i),
			)
		} else {
			// don't display the cell separator for this cell
			fmt.Fprintf(t.out, "%s%s",
				strings.Repeat(" ", colWidth+2),
				t.pCenter)
		}
	}

	if withNewLine {
		fmt.Fprint(t.out, t.newLine)
	}
}

// headerPrepadder yields a transform applied to header and footer cells,
// before padding.
func (t *Table) headerPrepadder() transformer {
	if t.autoFmt {
		return title
	}

	return identity
}

func (t *Table) isRightMost(i int) bool {
	end := len(t.cs) - 1
	return i == end && !t.borders.Right
}

func (t *Table) isLeftMost(i int) bool {
	return i == 0 && !t.borders.Left
}

func (t *Table) isEscSeq(params map[int]Formatter) bool {
	return len(params) > 0
}

func (t *Table) startOfLinePad() string {
	if !t.noWhiteSpace {
		return stringIf(t.borders.Left,
			t.pColumn, // |
			t.tablePadding,
		)
	}

	return NOPADDING
}

// transformer yields a functor to apply transforms on cell values.
func (t *Table) transformer(params map[int]Formatter) colTransformer {
	return func(i int) transformer {
		return func(in string) string {
			if t.isEscSeq(params) {
				in = format(in, params[i])
			}

			if t.isRightMost(i) {
				in = strings.TrimRightFunc(in, BlankSplitter)
			}

			return in
		}
	}
}

// Print heading information
func (t *Table) printHeader() {
	if len(t.headers) == 0 {
		return
	}

	padder := t.headerAlign.padder()
	maxHeight := t.rs[headerRowIdx]
	headerLines := normalizeRowHeight(t.headers, maxHeight)

	colLeftPad := func(in string, i, _ int) string {
		if len(in) == 0 && t.isRightMost(i) {
			return NOPADDING
		}

		if !t.noWhiteSpace {
			return SPACE
		}

		return NOPADDING
	}
	/*
		footer version: should we take into account noWhiteSpace??
		colLeftPad := func(in string, i, _ int) string {
			if len(in) == 0 && t.isRightMost(i) {
				return NOPADDING
			}
			return SPACE
		}
	*/

	colRightPad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		if !t.noWhiteSpace {
			return stringIf(
				t.isRightMost(i),
				SPACE, t.pColumn,
			)
		}

		return stringIf(
			t.isRightMost(i),
			SPACE, t.tablePadding,
		)
	}
	/* footer version
	colRightPad := func(in string, i, j int) string {
		if j == 0 {
			// right padding on first line of footer
			if len(in) == 0 {
				erasePad[i] = true

				return stringIf(t.isRightMost(i), NOPADDING, SPACE)
			}

			return stringIf(t.isRightMost(i), NOPADDING, t.pColumn)
		}

		if erasePad[i] {
			return stringIf(t.isRightMost(i), NOPADDING, SPACE)
		}

		return stringIf(t.isRightMost(i), NOPADDING, t.pColumn)
	}
	*/

	colMiddlePad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		if t.isEscSeq(t.headerParams) || !t.noWhiteSpace {
			return SPACE
		}

		return NOPADDING
	}
	/* footer version: why do have a space on header with ansi escape seq???
	colMiddlePad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		return SPACE
	}
	*/

	prepadding := t.headerPrepadder()
	transform := t.transformer(t.headerParams)

	for line := 0; line < maxHeight; line++ {
		fmt.Fprint(t.out, t.startOfLinePad())

		for col := 0; col < len(t.cs); col++ {
			value := headerLines[col][line]
			fmt.Fprint(t.out, colLeftPad(value, col, line))
			fmt.Fprint(t.out, transform(col)(
				padder(
					prepadding(value),
					SPACE, t.cs[col],
				),
			),
			)
			fmt.Fprint(t.out, colMiddlePad(value, col, line)) // TODO: simplify this (may be mergd with RightPad)
			fmt.Fprint(t.out, colRightPad(value, col, line))
		}

		fmt.Fprint(t.out, t.newLine)
	}

	if t.separatorAfterHeader {
		t.printSepLine(true)
	}
}

/* factorize with this?
func (t *Table) renderRow(
	lines [][]string,
	maxHeight int,
	cellPadder padFunc,
	leftPadder, rightPadder colPadder,
	prepadder transformer,
	transform colTransformer,
) {
	for line := 0; line < maxHeight; line++ {
		fmt.Fprint(t.out, t.startOfLinePad())

		for col := 0; col <= len(t.cs); col++ {
			colWidth := t.cs[col]
			value := lines[col][line]

			fmt.Fprint(t.out, leftPadder(value, col, line))
			fmt.Fprint(t.out, transform(col)(cellPadder(prepadder(value), SPACE, colWidth)))
			fmt.Fprint(t.out, rightPadder(value, col, line))
		}

		fmt.Fprint(t.out, t.newLine)
	}
}
*/

// Print column footer
func (t *Table) printFooter() {
	if len(t.footers) == 0 {
		return
	}

	end := len(t.cs) - 1
	padder := t.footerAlign.padder()
	maxHeight := t.rs[footerRowIdx]
	footerLines := normalizeRowHeight(t.footers, maxHeight)

	if !t.borders.Bottom {
		t.printSepLine(true)
	}

	colLeftPad := func(in string, i, _ int) string {
		if len(in) == 0 && t.isRightMost(i) {
			return NOPADDING
		}
		return SPACE
	}
	colMiddlePad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		return SPACE
	}
	erasePad := make([]bool, len(t.footers))
	colRightPad := func(in string, i, j int) string {
		if j == 0 {
			// right padding on first line of footer
			if len(strings.TrimRightFunc(in, BlankSplitter)) == 0 {
				erasePad[i] = true

				return stringIf(t.isRightMost(i), NOPADDING, SPACE)
			}

			return stringIf(t.isRightMost(i), NOPADDING, t.pColumn)
		}

		if erasePad[i] {
			return stringIf(t.isRightMost(i), NOPADDING, SPACE)
		}

		return stringIf(t.isRightMost(i), NOPADDING, t.pColumn)
	}

	prepadding := t.headerPrepadder()
	transform := t.transformer(t.footerParams)

	for line := 0; line < maxHeight; line++ {
		fmt.Fprint(t.out, t.startOfLinePad())

		for col := 0; col <= end; col++ {
			colWidth := t.cs[col]
			value := footerLines[col][line]

			fmt.Fprint(t.out, colLeftPad(value, col, line))
			fmt.Fprint(t.out, transform(col)(padder(prepadding(value), SPACE, colWidth)))
			fmt.Fprint(t.out, colMiddlePad(value, col, line))
			fmt.Fprint(t.out, colRightPad(value, col, line))
		}

		fmt.Fprint(t.out, t.newLine)
	}

	if t.separatorAfterFooter {
		t.printFooterSeparator()
	}
}

// print special separator line below the footer
func (t *Table) printFooterSeparator() {
	end := len(t.cs) - 1
	hasPrinted := false

	for col := 0; col <= end; col++ {
		colWidth := t.cs[col]
		pad := t.pRow
		center := t.pCenter
		length := len(t.footers[col][0])

		if length > 0 {
			hasPrinted = true
		}

		if length == 0 && !t.borders.Right {
			center = SPACE
		}

		if col == 0 {
			if length > 0 && !t.borders.Left {
				center = t.pRow
			}
			fmt.Fprint(t.out, center)
		}

		if length == 0 {
			pad = SPACE
		}

		if hasPrinted || t.borders.Left {
			pad = t.pRow
			center = t.pCenter
		}

		if center != SPACE {
			if col == end && !t.borders.Right {
				center = t.pRow
			}
		}

		if center == SPACE {
			if col < end && len(t.footers[col+1][0]) != 0 {
				if !t.borders.Left {
					center = t.pRow
				} else {
					center = t.pCenter
				}
			}
		}

		fmt.Fprintf(t.out, "%s%s%s%s",
			pad,
			strings.Repeat(pad, colWidth),
			pad,
			center)

	}

	fmt.Fprint(t.out, t.newLine)
}

// Print caption text
func (t Table) printCaption() {
	width := t.getTableWidth()
	paragraph, _ := WrapString(t.captionText, width) // TODO: use Wrapper

	for linecount := 0; linecount < len(paragraph); linecount++ {
		fmt.Fprintln(t.out, format(paragraph[linecount], t.captionParams))
	}
}

// Calculate the total number of characters in a row
//
// TODO: what happens when noBlankSpace is true?
// TODO: test when with or without borders
func (t Table) getTableWidth() int {
	var chars int

	col := displayWidth(t.pColumn)
	padding := displayWidth(t.tablePadding)

	//if t.borders.Left {
	chars += col
	chars += padding
	//}

	for _, width := range t.cs {
		chars += width
		chars += padding
		chars += col
	}

	//if t.borders.Right {
	chars += displayWidth(t.tablePadding)
	chars += col
	//}

	return chars

	// OLD return (chars + (3 * t.numColumns) + 2)
}

// printRows renders all row lines.
func (t Table) printRows() {
	for i, rowLines := range t.lines {
		t.printRow(rowLines, i)
	}
}

func (t *Table) cellAligner(col int) padFunc {
	return t.columnsAlign[col].padder()
}

// printRow renders a single table row.
//
// Adjust column alignment based on type
func (t *Table) printRow(columns [][]string, rowIdx int) {
	maxHeight := t.rs[rowIdx]
	numColumns := len(columns)
	columns = normalizeRowHeight(columns, maxHeight)

	transform := t.transformer(t.columnsParams)
	colLeftPad := func(in string, i, _ int) string {
		if t.isRightMost(i) {
			if !t.noWhiteSpace {
				if len(strings.TrimRightFunc(in, BlankSplitter)) > 0 {
					return stringIf(t.isLeftMost(i), SPACE, t.pColumn) + SPACE
				}
				return stringIf(t.isLeftMost(i), SPACE, t.pColumn)
			}

			return NOPADDING
		}

		if !t.noWhiteSpace {
			return stringIf(t.isLeftMost(i), SPACE, t.pColumn) + SPACE
		}

		return NOPADDING
	}

	colRightPad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		end := len(t.cs) - 1
		if !t.noWhiteSpace && t.borders.Right && i == end {
			return SPACE + stringIf(t.borders.Right, t.pColumn, SPACE)
		}

		return t.tablePadding
	}

	for line := 0; line < maxHeight; line++ {
		for col := 0; col < numColumns; col++ {
			padder := t.cellAligner(col)
			colWidth := t.cs[col]
			value := columns[col][line]

			fmt.Fprint(t.out, colLeftPad(value, col, line))
			fmt.Fprint(t.out, transform(col)(padder(value, SPACE, colWidth)))
			fmt.Fprint(t.out, colRightPad(value, col, line))
		}

		fmt.Fprint(t.out, t.newLine)
	}

	if t.separatorBetweenRows {
		t.printSepLine(true)
	}
}

// Print the rows of the table and merge the cells that are identical
func (t *Table) printRowsMergeCells() {
	var (
		previousLine      []string
		displayCellBorder []bool
		tmpWriter         bytes.Buffer
	)

	for i, lines := range t.lines {
		// We store the display of the current line in a tmp writer, as we need to know which border needs to be print above
		previousLine, displayCellBorder = t.printRowMergeCells(&tmpWriter, lines, i, previousLine)
		if i > 0 { // We don't need to print borders above first line
			if t.separatorBetweenRows {
				t.printLineOptionalCellSeparators(true, displayCellBorder)
			}
		}
		_, _ = tmpWriter.WriteTo(t.out)
	}

	if t.separatorBetweenRows {
		t.printSepLine(true)
	}
}

// Print Row Information to a writer and merge identical cells.
// Adjust column alignment based on type
func (t *Table) printRowMergeCells(writer io.Writer, columns [][]string, rowIdx int, previousLine []string) ([]string, []bool) {
	// Get Maximum Height
	max := t.rs[rowIdx]
	numColumns := len(columns)
	columns = normalizeRowHeight(columns, max)

	var displayCellBorder []bool
	for x := 0; x < max; x++ {
		for y := 0; y < numColumns; y++ {

			// Check if border is set
			fmt.Fprint(writer, conditionString((!t.borders.Left && y == 0), SPACE, t.pColumn))
			fmt.Fprint(writer, SPACE)

			str := columns[y][x]

			// Embedding escape sequence with column value
			if t.isEscSeq(t.columnsParams) {
				str = format(str, t.columnsParams[y])
			}

			if t.autoMergeCells {
				var mergeCell bool
				if t.columnsToAutoMergeCells != nil {
					// Check to see if the column index is in columnsToAutoMergeCells.
					if t.columnsToAutoMergeCells[y] {
						mergeCell = true
					}
				} else {
					// columnsToAutoMergeCells was not set.
					mergeCell = true
				}

				// Store the full line to merge mutli-lines cells
				fullLine := strings.TrimRight(strings.Join(columns[y], " "), " ")
				if len(previousLine) > y && fullLine == previousLine[y] && fullLine != "" && mergeCell {
					// If this cell is identical to the one above but not empty, we don't display the border and keep the cell empty.
					displayCellBorder = append(displayCellBorder, false)
					str = ""
				} else {
					// First line or different content, keep the content and print the cell border
					displayCellBorder = append(displayCellBorder, true)
				}
			}

			// This would print alignment
			// Default alignment  would use multiple configuration
			switch t.columnsAlign[y] {
			case AlignCenter: //
				fmt.Fprintf(writer, "%s", padCenter(str, SPACE, t.cs[y]))
			case AlignRight:
				fmt.Fprintf(writer, "%s", padLeft(str, SPACE, t.cs[y]))
			case AlignLeft:
				fmt.Fprintf(writer, "%s", padRight(str, SPACE, t.cs[y]))
			default:
				if isNumerical(str) {
					fmt.Fprintf(writer, "%s", padLeft(str, SPACE, t.cs[y]))
				} else {
					fmt.Fprintf(writer, "%s", padRight(str, SPACE, t.cs[y]))
				}
			}
			fmt.Fprint(writer, SPACE)
		}
		// Check if border is set
		// Replace with space if not set
		fmt.Fprint(writer, conditionString(t.borders.Left, t.pColumn, SPACE))
		fmt.Fprint(writer, t.newLine)
	}

	// The new previous line is the current one
	previousLine = make([]string, numColumns)
	for y := 0; y < numColumns; y++ {
		previousLine[y] = strings.TrimRight(strings.Join(columns[y], " "), " ") // Store the full line for multi-lines cells
	}

	// Returns the newly added line and wether or not a border should be displayed above.
	return previousLine, displayCellBorder
}

func (t *Table) parseDimension(str string, colKey, rowKey int) []string {
	var (
		raw      []string
		maxWidth int
	)

	raw = getLines(str)
	maxWidth = 0
	for _, line := range raw {
		if w := DisplayWidth(line); w > maxWidth {
			maxWidth = w
		}
	}

	// If wrapping, ensure that all paragraphs in the cell fit in the
	// specified width.
	if t.autoWrap {
		maxAllowedWidth, isDefined := t.ms[colKey]
		if !isDefined {
			maxAllowedWidth = t.mW
		}

		// If there's a maximum allowed width for wrapping, use that.
		if maxWidth > maxAllowedWidth {
			maxWidth = maxAllowedWidth
		}

		// In the process of doing so, we need to recompute maxWidth. This
		// is because perhaps a word in the cell is longer than the
		// allowed maximum width in maxAllowedWidth.
		newMaxWidth := maxWidth
		newRaw := make([]string, 0, len(raw))

		if t.reflowText {
			// Make a single paragraph of everything.
			raw = []string{strings.Join(raw, " ")}
		}

		for i, para := range raw {
			paraLines, _ := wrapString(para, maxWidth)
			for _, line := range paraLines {
				if w := displayWidth(line); w > newMaxWidth {
					newMaxWidth = w
				}
			}
			if i > 0 {
				newRaw = append(newRaw, SPACE)
			}
			newRaw = append(newRaw, paraLines...)
		}
		raw = newRaw
		maxWidth = newMaxWidth
	}

	// Store the new known maximum width.
	v, ok := t.cs[colKey]
	if !ok || v < maxWidth || v == 0 {
		t.cs[colKey] = maxWidth
	}

	// Remember the number of lines for the row printer.
	h := len(raw)
	v, ok = t.rs[rowKey]

	if !ok || v < h || v == 0 {
		t.rs[rowKey] = h
	}

	return raw
}
