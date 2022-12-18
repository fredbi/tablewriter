package tablewriter

// Credits to Oleku Konko, whose work published at
// github.com/olekukonko/tablewriter has greatly inspired
// this piece of software.

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	wrap "github.com/fredbi/tablewriter/tablewrappers"
)

const (
	headerRowIdx = -1
	footerRowIdx = -2
)

type (

	// Table renders a text table.
	Table struct {
		*options

		// internal multi-line representations of rows
		lines                   [][][]string
		headers                 [][]string
		footers                 [][]string
		numColumns              int
		columnsToAutoMergeCells map[int]bool
		columnsAlign            []HAlignment
		rowMaxHeight            map[int]int // max lines per cell

		wrappers
	}

	wrappers struct {
		cellWrapper func(row, col int) []string
	}
)

// New builds a new empty table writer.
func New(opts ...Option) *Table {
	t := &Table{
		lines:        [][][]string{},
		headers:      [][]string{},
		footers:      [][]string{},
		numColumns:   -1,
		rowMaxHeight: make(map[int]int),

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

func (t *Table) fillMaxWidths() {
	for i := 0; i < t.numColumns; i++ {
		_, isDefined := t.colMaxWidth[i]
		if !isDefined {
			t.colMaxWidth[i] = t.maxColWidth
		}
	}
}

// setNumColumns determines the number of columns for this table, aligned to the row
// (or header, or footer) with the largest number of columns.
func (t *Table) setNumColumns() {
	nCols := len(t.header)

	for _, row := range t.rows {
		if cols := len(row); cols > nCols {
			nCols = cols
		}
	}

	if cols := len(t.footer); cols > nCols {
		nCols = cols
	}

	// normalize all content to the same number of columns, adding trailing empty columns
	if len(t.header) > 0 {
		if missing := nCols - len(t.header); missing > 0 {
			padded := make([]string, missing)
			t.header = append(t.header, padded...)
		}
	}

	for i, row := range t.rows {
		if missing := nCols - len(row); missing > 0 {
			padded := make([]string, missing)
			t.rows[i] = append(t.rows[i], padded...)
		}
	}

	if len(t.footer) > 0 {
		if missing := nCols - len(t.footer); missing > 0 {
			padded := make([]string, missing)
			t.footer = append(t.footer, padded...)
		}
	}

	t.numColumns = nCols
}

func (t *Table) lastCol() int {
	return t.numColumns - 1
}

// center determines the character to print at an intersection line, based on the position and border.
func (t *Table) center(i int) string {
	switch {
	case i == -1 && !t.borders.Left: // ICI - strange: for isLeftMost we have i == 0 && !t.borders.Left
		fallthrough // -
	// return i == t.lastCol() && !t.borders.Right
	case t.isRightMost(i):
		return t.pRow // -
	default:
		return t.pCenter // +
	}
}

// PrintSepLine prints a separation line based on the row width.
//
// BUG(fred): this doesn't work well with noWhiteSpace
func (t *Table) printSepLine(withNewLine bool) {
	if !t.noWhiteSpace {
		fmt.Fprint(t.out, t.center(-1)) // -
	}

	for i := 0; i < t.numColumns; i++ {
		if !t.noWhiteSpace {
			fmt.Fprint(t.out, t.pRow) // -
		}
		fmt.Fprint(t.out, strings.Repeat(t.pRow, t.colWidth[i])) // -...-
		if !t.noWhiteSpace {
			fmt.Fprint(t.out, t.pRow) // -
		}
		fmt.Fprint(t.out, t.center(i)) // +|-
	}

	if withNewLine {
		fmt.Fprint(t.out, t.newLine)
	}
}

// Prints a separation line based on row width with or without cell separator.
//
// This is used when rendering merged cells.
//
// TODO(fred): this should be factorized with printSepLine
func (t *Table) printLineOptionalCellSeparators(withNewLine bool, displayCellSeparator []bool) {
	fmt.Fprint(t.out, t.pCenter) // ?? + never -???

	for i := 0; i < t.numColumns; i++ {
		colWidth := t.colWidth[i]

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
				strings.Repeat(SPACE, colWidth+2),
				t.pCenter)
		}
	}

	if withNewLine {
		fmt.Fprint(t.out, t.newLine)
	}
}

// headerPrepadder yields a transform applied BEFORE padding to header or footer cells.
func (t *Table) headerPrepadder() transformer {
	if t.titler != nil {
		return t.titler.Title
	}

	return identity
}

// isRightMost position on this row?
func (t *Table) isRightMost(i int) bool {
	return i == t.lastCol() && !t.borders.Right
}

// isLefttMost position on this row?
func (t *Table) isLeftMost(i int) bool {
	return i == 0 && !t.borders.Left
}

// hasEscSeq indicates if we have ANSI escape sequence parameters
// set for at least one colum.
func (t *Table) hasEscSeq(params map[int]Formatter) bool {
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
			if t.hasEscSeq(params) {
				// apply formatting escape sequence, if any
				in = format(in, params[i])
			}

			if t.isRightMost(i) {
				// remove extraneous trailing blanks on rows without a right border
				in = strings.TrimRightFunc(in, wrap.BlankSplitter)
			}

			return in
		}
	}
}

func (t *Table) printHeader() {
	if len(t.headers) == 0 {
		return
	}

	padder := t.headerAlign.padder()
	aligner := func(_ int) padFunc { return padder }
	maxHeight := t.rowMaxHeight[headerRowIdx]
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

	colRightPad := func(_ string, i, _ int) string {
		if t.isRightMost(i) {
			return NOPADDING
		}

		var middlePad string
		if t.hasEscSeq(t.headerParams) || !t.noWhiteSpace { // why when escape seq???
			middlePad = SPACE
		}

		if !t.noWhiteSpace {
			return middlePad + stringIf(
				t.isRightMost(i),
				SPACE, t.pColumn,
			)
		}

		return middlePad + stringIf(
			t.isRightMost(i),
			SPACE, t.tablePadding,
		)
	}

	prepadding := t.headerPrepadder()
	transform := t.transformer(t.headerParams)

	t.renderRowWithPadding(
		headerLines,
		maxHeight,
		aligner,
		colLeftPad, colRightPad,
		prepadding,
		transform,
		t.startOfLinePad,
	)

	if t.separatorAfterHeader {
		t.printSepLine(true)
	}
}

// renderRowWithPadding captures the rendering logic to display a single row of any type.
//
// NOTE: further refactoring is desirable to eventually use similar padding approaches for
// header, footer and table row.
func (t *Table) renderRowWithPadding(
	cells [][]string,
	maxHeight int,
	cellAligner colAligner,
	leftPadder, rightPadder colPadder,
	prepadder transformer,
	transform colTransformer,
	lineStarter func() string,
) {
	if prepadder == nil {
		prepadder = identity
	}

	for line := 0; line < maxHeight; line++ {
		if lineStarter != nil {
			fmt.Fprint(t.out, lineStarter())
		}

		for col := 0; col < t.numColumns; col++ {
			cellPadder := cellAligner(col) // each column may use a different alignment
			colWidth := t.colWidth[col]
			value := cells[col][line]

			fmt.Fprint(t.out, leftPadder(value, col, line))
			fmt.Fprint(t.out, transform(col)(
				cellPadder(
					prepadder(value),
					SPACE, colWidth,
				),
			))
			fmt.Fprint(t.out, rightPadder(value, col, line))
		}

		fmt.Fprint(t.out, t.newLine)
	}
}

func (t *Table) printFooter() {
	if len(t.footers) == 0 {
		return
	}

	padder := t.footerAlign.padder()
	aligner := func(_ int) padFunc { return padder }
	maxHeight := t.rowMaxHeight[footerRowIdx]
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

	erasePad := make([]bool, len(t.footers))
	colRightPad := func(in string, i, j int) string {
		if j == 0 {
			// right padding on first line of footer
			if len(strings.TrimRightFunc(in, wrap.BlankSplitter)) == 0 {
				erasePad[i] = true

				return stringIf(t.isRightMost(i), NOPADDING, SPACE+SPACE)
			}

			return stringIf(t.isRightMost(i), NOPADDING, SPACE+t.pColumn)
		}

		if erasePad[i] {
			return stringIf(t.isRightMost(i), NOPADDING, SPACE+SPACE)
		}

		return stringIf(t.isRightMost(i), NOPADDING, SPACE+t.pColumn)
	}

	prepadding := t.headerPrepadder()
	transform := t.transformer(t.footerParams)

	t.renderRowWithPadding(
		footerLines,
		maxHeight,
		aligner,
		colLeftPad, colRightPad,
		prepadding,
		transform,
		t.startOfLinePad,
	)

	if t.separatorAfterFooter {
		t.printFooterSeparator()
	}
}

// print special separator line below the footer
func (t *Table) printFooterSeparator() {
	hasPrinted := false

	for col := 0; col < t.numColumns; col++ {
		colWidth := t.colWidth[col]
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
			if col == t.lastCol() && !t.borders.Right {
				center = t.pRow
			}
		}

		if center == SPACE {
			if col < t.lastCol() && len(t.footers[col+1][0]) != 0 {
				if !t.borders.Left {
					center = t.pRow
				} else {
					center = t.pCenter
				}
			}
		}

		fmt.Fprint(t.out, pad)
		fmt.Fprint(t.out, strings.Repeat(pad, colWidth))
		fmt.Fprint(t.out, pad)
		fmt.Fprint(t.out, center)
	}

	fmt.Fprint(t.out, t.newLine)
}

// Print caption text
func (t Table) printCaption() {
	width := t.getTableWidth()
	captionWrapper := wrap.NewDefault()
	paragraph := captionWrapper.WrapString(t.captionText, width)

	for linecount := 0; linecount < len(paragraph); linecount++ {
		fmt.Fprintln(t.out, format(paragraph[linecount], t.captionParams))
	}
}

// getTableWidth yields the total number of characters in a row
func (t Table) getTableWidth() int {
	chars := t.overhead()

	for _, width := range t.colWidth {
		chars += width
	}

	return chars
}

// Overhead yields the amount extra padding and separators needed to display the table.
func (t Table) Overhead() int {
	return t.overhead()
}

func (t Table) overhead() int {
	var chars int

	colSepWidth := wrap.DisplayWidth(t.pColumn)
	paddingWidth := wrap.DisplayWidth(t.tablePadding)

	if !t.noWhiteSpace {
		chars += paddingWidth * t.numColumns * 2
		chars += colSepWidth * (t.numColumns + 1)

		return chars
	}

	if t.numColumns > 0 {
		return paddingWidth * (t.numColumns - 1)
	}

	return 0
}

// printRows renders all multi-lines rows
func (t Table) printRows() {
	for i, rowLines := range t.lines {
		t.printRow(rowLines, i)
	}
}

func (t *Table) cellAligner(col int) padFunc {
	return t.columnsAlign[col].padder()
}

// printRow renders a single multi-lines row
func (t *Table) printRow(columns [][]string, rowIdx int) {
	maxHeight := t.rowMaxHeight[rowIdx]
	columns = normalizeRowHeight(columns, maxHeight)

	aligner := t.cellAligner
	transform := t.transformer(t.columnsParams)

	colLeftPad := func(in string, i, _ int) string {
		if t.isRightMost(i) {
			if !t.noWhiteSpace {
				if len(strings.TrimRightFunc(in, wrap.BlankSplitter)) > 0 {
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

		if !t.noWhiteSpace && t.borders.Right && i == t.lastCol() {
			return SPACE + stringIf(t.borders.Right, t.pColumn, SPACE)
		}

		return t.tablePadding
	}

	t.renderRowWithPadding(
		columns,
		maxHeight,
		aligner,
		colLeftPad, colRightPad,
		nil, // at this moment, we don't have cell transforms configurable for rows
		transform,
		nil, // at this moment, the padding logic for rendering row is different: in that case, no start-of-line padding
	)

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
		// we store the display of the current line in a tmp writer, as we need to know which border needs to be print above
		previousLine, displayCellBorder = t.printRowMergeCells(&tmpWriter, lines, i, previousLine)
		if i > 0 { // we don't need to print borders above first line
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

// print row Information to a writer and merge identical cells.
// Adjust column alignment based on type
//
// TODO(fred): this should be refactored along the same lines as printHeader and printRow.
func (t *Table) printRowMergeCells(writer io.Writer, columns [][]string, rowIdx int, previousLine []string) ([]string, []bool) {
	max := t.rowMaxHeight[rowIdx]
	numColumns := len(columns)
	columns = normalizeRowHeight(columns, max)

	var displayCellBorder []bool
	for x := 0; x < max; x++ {
		for y := 0; y < numColumns; y++ {

			// Check if border is set
			fmt.Fprint(writer, stringIf((!t.borders.Left && y == 0), SPACE, t.pColumn))
			fmt.Fprint(writer, SPACE)

			str := columns[y][x]

			// Embedding escape sequence with column value
			if t.hasEscSeq(t.columnsParams) {
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
				fullLine := strings.TrimRight(strings.Join(columns[y], SPACE), SPACE)
				if len(previousLine) > y && fullLine == previousLine[y] && fullLine != NOPADDING && mergeCell {
					// If this cell is identical to the one above but not empty, we don't display the border and keep the cell empty.
					displayCellBorder = append(displayCellBorder, false)
					str = NOPADDING
				} else {
					// First line or different content, keep the content and print the cell border
					displayCellBorder = append(displayCellBorder, true)
				}
			}

			// This would print alignment
			// Default alignment  would use multiple configuration
			switch t.columnsAlign[y] {
			case AlignCenter: //
				fmt.Fprintf(writer, "%s", padCenter(str, SPACE, t.colWidth[y]))
			case AlignRight:
				fmt.Fprintf(writer, "%s", padLeft(str, SPACE, t.colWidth[y]))
			case AlignLeft:
				fmt.Fprintf(writer, "%s", padRight(str, SPACE, t.colWidth[y]))
			default:
				if isNumerical(str) {
					fmt.Fprintf(writer, "%s", padLeft(str, SPACE, t.colWidth[y]))
				} else {
					fmt.Fprintf(writer, "%s", padRight(str, SPACE, t.colWidth[y]))
				}
			}
			fmt.Fprint(writer, SPACE)
		}

		// Check if border is set
		// Replace with space if not set
		fmt.Fprint(writer, stringIf(t.borders.Left, t.pColumn, SPACE))
		fmt.Fprint(writer, t.newLine)
	}

	// The new previous line is the current one
	previousLine = make([]string, numColumns)
	for y := 0; y < numColumns; y++ {
		previousLine[y] = strings.TrimRight(strings.Join(columns[y], SPACE), SPACE) // store the full line for multi-lines cells
	}

	// Returns the newly added line and wether or not a border should be displayed above.
	return previousLine, displayCellBorder
}

func (t *Table) setColWidth(col, width int) {
	previous := t.colWidth[col]
	if previous == 0 || previous < width {
		t.colWidth[col] = width
	}
}

func (t *Table) setRowHeight(row, height int) {
	previous := t.rowMaxHeight[row]

	if previous == 0 || previous < height {
		t.rowMaxHeight[row] = height
	}
}
