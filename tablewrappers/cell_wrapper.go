package tablewrappers

import (
// "log"
)

type (
	// RowWrapper wraps the content of a table with a single constraint on the table width.
	RowWrapper struct {
		*wrapOptions

		matrix       [][]string
		rowLimit     int
		wordSplitter Splitter
		noOp         bool
		columns      columns
	}

	// DefaultCellWrapper wraps the content of a table with predefined constraints on column widths.
	DefaultCellWrapper struct {
		*DefaultWrapper
		matrix      [][]string
		colMaxWidth map[int]int // max width for a column
	}
)

func NewDefaultCellWrapper(matrix [][]string, colMaxWidth map[int]int, opts ...Option) *DefaultCellWrapper {
	w := &DefaultCellWrapper{
		DefaultWrapper: NewDefault(opts...),
		matrix:         matrix,
		colMaxWidth:    colMaxWidth,
	}

	if w.colMaxWidth == nil {
		w.colMaxWidth = make(map[int]int)
	}

	return w
}

func (w *DefaultCellWrapper) WrapCell(row, col int) []string {
	limit := w.colMaxWidth[col]
	if limit == 0 {
		return []string{w.matrix[row][col]} // no op
	}

	return w.WrapString(w.matrix[row][col], limit)
}

// TODO: introduce colMinWidth, colMaxWidth local limits for backward-compatible layout
// TODO: bug wrapping control sequences bugs when words are broken
func NewRowWrapper(matrix [][]string, rowWidthLimit int, opts ...Option) *RowWrapper {
	w := &RowWrapper{
		wrapOptions: optionsWithDefaults(opts),
		matrix:      matrix,
		rowLimit:    rowWidthLimit,
	}
	w.wordSplitter = composeSplitters(w.splitters)
	w.prepare() // perform all the computations to derive column widths and wrap as few cells as possible

	return w
}

func (w *RowWrapper) WrapCell(row, col int) []string {
	if w.noOp {
		return []string{w.matrix[row][col]}
	}

	return w.columns[col].cells[row].content
}

func (w *RowWrapper) prepare() {
	if w.rowLimit < 0 || len(w.matrix) == 0 {
		// short circuit: wrapping can't achieve the limit
		w.noOp = true

		return
	}

	_, cols := buildMatrix(w.matrix, w.wordSplitter)

	currentWidth := cols.TotalWidth()
	if currentWidth < w.rowLimit {
		// short circuit: nothing to be wrapped
		w.noOp = true

		return
	}

	// TODO: preliminary check on single words: if some words are wider than the total width
	// cols.BreakLongestWords(wordBreakLevel, w.rowLimit-currentWidth, w.wordSplitter)
	// TODO: adaptable # buckets vs # rows in the matrix
	// TODO: if maxWordLength in a column > limit : break words at once in this column
	// TODO: if maxWordLength in a column + Sum(minWordLength) other columns > limit break words in this column

	cols.SortRows() // each column gets its rows sorted by width, widest first
	cols.Sort()     // columns get sorted, so that the first element is the widest

	for _, wordBreakLevel := range []breakLevel{
		breakNone,
		breakOnSeps,
		breakAnywhere,
	} {
		cols.BreakLongestWords(wordBreakLevel, currentWidth-w.rowLimit, w.wordSplitter)

		// shrink columns, widest-first
		currentWidth = w.shrinkColumns(cols)
		// log.Printf("DEBUG: current width after shrink=%d", currentWidth)

		if currentWidth <= w.rowLimit {
			break
		}
	}
	// log.Printf("DEBUG: last current width=%d", currentWidth)

	// reorder columns and rows by their natural order
	cols.SortNatural()

	w.columns = cols
}

// shrinkColumns rebalances words in the cells of columns.
// It returns the total width of the re-arranged table.
//
// This is a best-effort since this step doesn't break words.
//
// This function assesses the histogram of widths for columns, assuming columns come already sorted
// widest-first, then shrinks each candidate column to the next bucket.
func (w *RowWrapper) shrinkColumns(cols columns) int {
	var currentWidth int

LOOP:
	for bucket := 0; bucket < numBuckets-1; bucket++ { // progressively more agressive: 90%-width, 80%-width, ...
		for _, col := range cols { // iterate over columns, widest first
			col.SetPValues(numBuckets) // computes the fixed-bucket histogram of widths (param to capture pass on words later on)
			limit := max(col.pvalues[bucket], w.minColWidth[col.j])

			if maxw := w.maxColWidth[col.j]; maxw > 0 {
				limit = min(limit, maxw)
			}

			if limit > w.rowLimit { // maybe we should do this in a first pass
				continue LOOP // the column p-value cannot work. Skip to the next bucket
			}

			if col.maxWidth <= limit {
				continue // the p-value for this bucket did not result in a signficant decrease. Skip to the next column.
			}

			// try with limiting the width to the max width of p% of values in this column
			col.WrapCells(limit, w.wordSplitter)

			currentWidth = col.TotalRowWidth()
			if currentWidth <= w.rowLimit {
				return currentWidth
			}
		}
	}

	return currentWidth
}
