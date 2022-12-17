package tablewrappers

import (
	"log"
	// "github.com/davecgh/go-spew/spew"
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
	log.Printf("RowWrapper limit: %d", w.rowLimit)
	if w.rowLimit < 0 {
		// short circuit: wrapping can't achieve the limit
		w.noOp = true

		return
	}

	_, cols := buildMatrix(w.matrix, w.wordSplitter)
	log.Printf("RowWrapper TotalWidth: %d", cols.TotalWidth())
	if cols.TotalWidth() < w.rowLimit {
		// short circuit: nothing to be wrapped
		w.noOp = true

		return
	}
	// spew.Dump(columns)

	cols.SortRows() // each column gets its rows sorted by width, widest first
	cols.Sort()     // columns get sorted, so that the first element is the widest

	// TODO: run in passes
	// 1. First pass: try wrapping columns, no word breaking
	// 2. Word breaking on natural boundaries
	// 2. Word breaking anywhere
	// TODO: use lengths matrix in p-values
LOOP:
	for bucket := 0; bucket < buckets-1; bucket++ { // progressively more agressive: 90%-width, 80%-width, ...
		for _, col := range cols { // iterate over columns, widest first
			log.Printf("assessing bucket[%d] col[%d]", bucket, col.j)
			col.SetPValues(-9999) // computes the fixed-bucket histogram of widths (param to capture pass on words later on)

			limit := col.pvalues[bucket]
			if limit > w.rowLimit { // maybe we should do this in a first pass
				log.Printf("skip to next bucket: p-value[%d]=%d (rowLimit=%d)", bucket, limit, w.rowLimit)

				continue LOOP // the column p-value cannot work. Skip to the next bucket
			}

			if col.maxWidth <= limit {
				log.Printf("skip to next column: col.maxWidth=%d (limit=%d)", col.maxWidth, limit)

				continue // the p-value for this bucket did not result in a signficant decrease. Skip to the next column.
			}

			// try with limiting the width to the max width of p% of values in this column
			log.Printf("RowWrapper WrapCells: %d [%T]", limit, col)
			col.WrapCells(limit)

			if col.TotalWidth() <= w.rowLimit {
				break LOOP
			}
		}
	}

	// reorder columns and rows by their natural order
	cols.SortNatural()

	w.columns = cols
}
