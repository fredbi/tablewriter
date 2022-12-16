package tablewrappers

import (
	"log"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// buckets to determine p-values.
// For instance with buckets = 10, we'll have buckets as deciles: 10%; 20%, ..., 90%
const buckets = 10

type (
	columns []*column
	cells   []*cell
	rows    []*row

	column struct {
		j        int
		maxWidth int
		rows     rows
		cells    cells
		pvalues  []int
	}

	row struct {
		i          int
		totalWidth int
		cells      cells
	}

	cell struct {
		i       int
		j       int
		content []string
		words   []string
		width   int
		// passNo      int // TODO
		wordLengths   [][]int // triangular matrix of paragraphs made up with words
		maxWordLength int
	}
)

func newCell(i, j int, content []string, splitter Splitter) *cell {
	c := &cell{
		i:       i,
		j:       j,
		content: content,
		width:   cellWidth(content),
	}

	for _, line := range content {
		c.words = append(c.words, strings.FieldsFunc(line, splitter)...)
	}

	c.wordLengths, c.maxWordLength = buildLengthsMatrix(c.words, 1) // build a matrix of the lengths of all word arrangements into single-line paragraphs [this is done once]

	return c
}

func newColumn(j int, rows rows) *column {
	c := &column{
		j:    j,
		rows: rows,
	}
	c.maxWidth = cellsMaxWidth(c.Values())
	c.cells = make(cells, 0, rows.MaxLen())

	for _, r := range rows {
		var colFoundInRow bool
		for _, cell := range r.cells {
			if cell.j == j {
				c.cells = append(c.cells, cell)
				colFoundInRow = true

				break
			}
		}

		if !colFoundInRow {
			// pad cells with an empty cell
			padCell := newCell(r.i, j, []string{""}, BlankSplitter)
			c.cells = append(c.cells, padCell)
			r.cells = append(r.cells, padCell)
		}
	}

	return c
}

func newRow(i int, cells cells) *row {
	return &row{
		i:          i,
		cells:      cells,
		totalWidth: cells.TotalWidth(),
	}
}

func buildMatrix(matrix [][]string, splitter Splitter) (rows, columns) {
	lines := make(rows, 0, len(matrix))
	maxCols := 0

	for i, row := range matrix {
		numCols := 0
		c := make(cells, 0, len(row))

		for j, content := range row {
			numCols++
			c = append(c, newCell(i, j, []string{content}, splitter))
		}

		lines = append(lines, newRow(i, c))
		if numCols > maxCols {
			maxCols = numCols
		}
	}

	cols := make(columns, 0, maxCols)
	for j := 0; j < maxCols; j++ {
		cols = append(cols, newColumn(j, lines))
	}

	return lines, cols
}

func (c columns) Less(i, j int) bool {
	return c[i].maxWidth > c[j].maxWidth
}

func (c columns) Len() int {
	return len(c)
}

func (c columns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// SortRows in a set of columns sorts the inner rows by their (descending) total width,
func (c columns) SortRows() {
	for _, col := range c {
		sort.Stable(col.rows)
		sort.Stable(col.cells)
	}
}

func (c columns) Sort() {
	sort.Stable(c)
}

func (c columns) SortNatural() {
	for _, col := range c {
		col.rows.SortNatural()
		col.cells.SortNatural()
	}

	sort.Slice(c, func(i, j int) bool {
		return c[i].j < c[j].j
	})
}

func (c columns) TotalWidth() int {
	total := 0

	for _, col := range c {
		total += col.maxWidth
	}

	return total
}

func (c cells) Less(i, j int) bool {
	return c[i].width > c[j].width
}

func (c cells) Len() int {
	return len(c)
}

func (c cells) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c cells) TotalWidth() int {
	total := 0

	for _, cell := range c {
		total += cell.width
	}

	return total
}

func (c cells) Values() [][]string {
	values := make([][]string, len(c))

	for _, cell := range c {
		values[cell.j] = cell.content
	}

	return values
}

func (c cells) SortNatural() {
	sort.Slice(c, func(i, j int) bool {
		if c[i].i == c[j].i {
			return c[i].j < c[j].j
		}

		return c[i].i < c[j].i
	})
}

func (c rows) Less(i, j int) bool {
	return c[i].totalWidth > c[j].totalWidth
}

func (c rows) Len() int {
	return len(c)
}

func (c rows) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c rows) MaxLen() int {
	maxLen := 0
	for _, r := range c {
		if l := r.Len(); l > maxLen {
			maxLen = l
		}
	}

	return maxLen
}

func (c rows) SortNatural() {
	for _, r := range c {
		r.cells.SortNatural()
	}

	sort.Slice(c, func(i, j int) bool {
		return c[i].i < c[j].i
	})
}

func (r row) Len() int {
	return r.cells.Len()
}

func (r row) Values() [][]string {
	return r.cells.Values()
}

func (r row) TotalWidth() int {
	return r.cells.TotalWidth()
}

func (c column) SortRows() {
	sort.Stable(c.rows)
}

func (c column) Rows() rows {
	return c.rows
}

func (c *column) SetPValues(_ int) {
	// TODO: use word-level widths to enrich the histogram
	sort.Stable(c.cells)
	sum := 0.00
	cdf := make([]float64, 0, len(c.cells))
	pvalues := make([]int, 0, buckets-1)
	// spew.Dump(c.cells)

	for i := len(c.cells) - 1; i >= 0; i-- {
		sum += float64(c.cells[i].width)
		// log.Printf("DEBUG: cell[i=%d] width=%d, sum=%f", i, c.cells[i].width, sum)
		cdf = append(cdf, sum)
	}

	if sum > 0 {
		for i := range cdf {
			cdf[i] /= sum
		}

	LOOP:
		for bucket := 1; bucket < buckets; bucket++ { // at most buckets - 1 pvalues
			// log.Printf("search pvalue for bucket[%d]", bucket)
			threshold := 1 - float64(bucket)/float64(buckets) // e.g: 90% for bucket 1, 80% for bucket 2, etc.

			for i := len(cdf) - 1; i >= 0; i-- {
				// log.Printf("cdf[%d] for bucket[%d]", i, bucket)
				val := cdf[i]
				var next float64
				if i < len(cdf)-1 {
					next = cdf[i+1]
				}

				// log.Printf("DEBUG: bucket[%d], i=%d, val=%f, next=%f, threshold=%f", bucket, i, val, next, threshold)

				if val <= threshold && (next > threshold || i >= len(cdf)-1) {
					// the captured pvalue represents the width under which fall bucket/buckets of cells in this column
					pvalue := c.cells[i].width
					// log.Printf("DEBUG: bucket[%d], i+lastBucket=%d, pvalue=%d", bucket, i, pvalue)
					pvalues = append(pvalues, pvalue)

					break
				}

				if i == 0 {
					break LOOP // no point searching for more buckets
				}
			}
		}
	}

	// ensure that p-values vectors are all of the same normalized size
	needsPadding := buckets - 1 - len(pvalues)
	for i := 0; i < needsPadding; i++ {
		pvalues = append(pvalues, c.maxWidth)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(pvalues)))
	// log.Printf("DEBUG: pvalues: %v", pvalues)

	c.pvalues = pvalues
}

func (c column) Values() [][]string {
	values := make([][]string, len(c.rows))

	for _, row := range c.rows {
		var content []string
		for _, cell := range row.cells {
			if cell.j == c.j {
				content = cell.content

				break
			}
		}

		if content == nil {
			// col not defined for this row
			content = []string{""}
		}
		values[row.i] = content
	}

	return values
}

// WrapCells updates all cells in this column with a wrapped version to the new width limit.
//
// NOTE: we don't update the p-values, which remain in their initial state.
// No need to update the word lengths matrix.
func (c *column) WrapCells(limit int) {
	log.Printf("before")
	spew.Dump(c.cells)
	for _, cell := range c.cells {
		if limit >= cell.width {
			continue
		}
		log.Printf("wrapping cell [col=%d][row=%d][limit: %d] %q", c.j, cell.i, limit, cell.words)
		lines := wrapMultiline(cell.words, limit) // wrap whole words over multiple lines
		cell.content = lines
		cell.width = cellWidth(lines)
	}

	c.maxWidth = cellsMaxWidth(c.Values())
	log.Printf("after")
	spew.Dump(c.cells)
}

// TotalWidth is the total width of all the rows that this column contains.
func (c column) TotalWidth() int {
	maxTotal := 0

	for _, row := range c.rows {
		if w := row.TotalWidth(); w > maxTotal {
			maxTotal = w
		}
	}

	return maxTotal
}
