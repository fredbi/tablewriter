package tablewrappers

type RowWrapper struct {
	*wrapOptions

	matrix       [][]string
	rowLimit     int
	wordSplitter Splitter
	noOp         bool
	columns      columns
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
	/*
		if w.noOp {
			return []string{w.matrix[row][col]}
		}
	*/

	return w.columns[col].cells[row].content
}

func (w *RowWrapper) prepare() {
	_, columns := buildMatrix(w.matrix, w.wordSplitter)
	if w.rowLimit < 0 || columns.TotalWidth() < w.rowLimit {
		w.noOp = true

		return // nothing to be wrapped.
	}
	// spew.Dump(columns)

	columns.SortRows() // each column gets its rows sorted by width, widest first
	columns.Sort()     // columns get sorted, so that the first element is the widest

	// TODO: run in passes
	// 1. First pass: try wrapping columns, no word breaking
	// 2. Word breaking on natural boundaries
	// 2. Word breaking anywhere
	// TODO: use lengths matrix in p-values
LOOP:
	for bucket := 0; bucket < buckets-1; bucket++ { // progressively more agressive: 90%-width, 80%-width, ...
		for _, col := range columns { // iterate over columns, widest first
			col.SetPValues(-9999) // computes the fixed-bucket histogram of widths (param to capture pass on words later on)

			limit := col.pvalues[bucket]
			if col.maxWidth <= limit {
				continue // the p-value for this bucket did not result in a signficant decrease. Skip to the next column.
			}

			// try with limiting the width to the max width of p% of values in this column
			col.WrapCells(limit)

			if col.TotalWidth() <= w.rowLimit {
				break LOOP
			}
		}
	}

	// reorder columns and rows by their natural order
	columns.SortNatural()

	w.columns = columns
}
