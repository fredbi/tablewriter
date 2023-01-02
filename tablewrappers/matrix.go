package tablewrappers

import (
	// "log"
	"math"
	"sort"
	"strings"
	// "github.com/davecgh/go-spew/spew"
)

// numBuckets to determine p-values.
//
// For instance with buckets = 10, we'll have buckets as deciles: 10%; 20%, ..., 90%
const numBuckets = 10

type (
	columns []*column
	cells   []*cell
	rows    []*row
	ratios  []ratio

	column struct {
		j          int
		maxWidth   int
		rows       rows
		cells      cells
		pvalues    []int
		minAllowed int // TODO
		maxAllowed int
	}

	row struct {
		i     int
		cells cells
	}

	cell struct {
		i       int
		j       int
		content []string
		// TODO: should index content that comes from break words

		width         int
		wordLengths   [][]int // triangular matrix of paragraphs made up with words
		maxWordLength int
		minWordLength int
	}

	ratio struct {
		j int
		r float64
	}
)

func newCell(i, j int, content []string, splitter Splitter) *cell {
	c := &cell{
		i:       i,
		j:       j,
		content: content,
		width:   cellWidth(content),
	}

	words := []string{}
	for _, line := range c.content {
		words = append(words, strings.FieldsFunc(line, splitter)...)
	}

	c.wordLengths, c.maxWordLength, c.minWordLength = buildLengthsMatrix(words, 1) // build a matrix of the lengths of all word arrangements into single-line paragraphs [this is done once]

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
		i:     i,
		cells: cells,
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

// TotalWidth yields the width of the table, adding up the max width of all columns.
func (c columns) TotalWidth() int {
	total := 0

	for _, col := range c {
		total += col.maxWidth
	}

	return total
}

// WordMaxWiths returns the maximum single word length in the set of colums.
// The result is provided in the order of the columns collection.
func (c columns) WordMaxWidths() ([]int, int) {
	wordWidths := make([]int, 0, len(c))
	var total int

	for _, col := range c {
		w := col.WordMaxWidth()
		wordWidths = append(wordWidths, w)
		total += w
	}

	return wordWidths, total
}

func newRatio(j int, r float64) ratio {
	return ratio{
		j: j,
		r: r,
	}
}

func (r ratios) Len() int {
	return len(r)
}

func (r ratios) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ratios) Less(i, j int) bool {
	return r[i].r > r[j].r
}

func (r ratios) Sort() {
	sort.Stable(r)
}

func (r ratios) SortNatural() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].j < r[j].j
	})
}

// WordLengthTargets computes the target lengths for colums,
// spreading word-breaking across columns according to their
// respective maximum // word lengths:
// longest words get broken more aggressively.
func (c columns) WordLengthTargets(gap int) []int {
	wordWidths, total := c.WordMaxWidths()

	spreads := make(ratios, len(wordWidths))
	for j := range wordWidths {
		spreads[j] = newRatio(j, float64(wordWidths[j])/float64(total))
	}
	spreads.Sort()

	var (
		sum  int
		full bool
	)

	for _, spread := range spreads {
		reduction := int(math.Ceil(float64(gap) * spread.r))
		sum += reduction
		if sum >= gap {
			if !full {
				reduction -= sum - gap
				full = true
			} else {
				reduction = 0
			}
		}

		wordWidths[spread.j] -= reduction
	}

	return wordWidths
}

// BreakLongestWords breaks down the longest words in a set of colums in order to tentatively
// reduce the provided gap.
//
// "gap" is the width reduction target.
func (c columns) BreakLongestWords(wordBreakLevel breakLevel, gap int, splitter Splitter) {
	if wordBreakLevel == breakNone || gap <= 0 {
		return
	}

	for j, col := range c {
		targets := c.WordLengthTargets(gap) // targets are refreshed after every break.
		col.BreakLongestWords(wordBreakLevel, targets[j], splitter)

		// update the max width of the column
		col.maxWidth = cellsMaxWidth(col.Values())
	}
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

func (c cells) Sort() {
	sort.Stable(c)
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
	return c[i].cells.TotalWidth() > c[j].cells.TotalWidth()
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

// SetPValues computes the quantile values on a column widths.
//
// This ways we know the width thresholds that affect x% of the cells in a column.
func (c *column) SetPValues(buckets int) {
	c.cells.Sort()
	sum := 0.00
	cdf := make([]float64, 0, len(c.cells))
	pvalues := make([]int, 0, buckets-1)

	for i := len(c.cells) - 1; i >= 0; i-- {
		sum += float64(c.cells[i].width)
		cdf = append(cdf, sum)
	}

	if sum > 0 {
		for i := range cdf {
			cdf[i] /= sum
		}

	LOOP:
		for bucket := 1; bucket < buckets; bucket++ { // at most buckets - 1 pvalues
			threshold := 1 - float64(bucket)/float64(buckets) // e.g: 90% for bucket 1, 80% for bucket 2, etc.

			for i := len(cdf) - 1; i >= 0; i-- {
				val := cdf[i]
				var next float64
				if i < len(cdf)-1 {
					next = cdf[i+1]
				}

				if val <= threshold && (next > threshold || i >= len(cdf)-1) {
					// the captured pvalue represents the width under which fall bucket/buckets of cells in this column
					pvalue := c.cells[i].width
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

	c.pvalues = pvalues
}

// Values returns the multiline content of all cells in a colum.
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
func (c *column) WrapCells(limit int, splitter Splitter) {
	for _, cell := range c.cells {
		if limit >= cell.width {
			continue
		}

		lines := make([]string, 0, len(cell.content))
		for _, line := range cell.content {
			words := strings.FieldsFunc(line, splitter)
			lines = append(lines, wrapMultiline(words, limit, 1)...) // wrap whole words over multiple lines
		}
		cell.content = lines
		cell.width = cellWidth(lines)
	}

	c.maxWidth = cellsMaxWidth(c.Values())
}

// TotalWidth is the total width of all the rows that this column contains.
func (c column) TotalRowWidth() int {
	maxPerColumn := make(map[int]int)
	maxTotal := 0

	for _, row := range c.rows {
		// log.Printf("DEBUG: row with=%d", row.TotalWidth())
		for _, cell := range row.cells {
			if w := cell.width; w > maxPerColumn[cell.j] {
				maxPerColumn[cell.j] = w
			}
		}
	}

	for _, w := range maxPerColumn {
		maxTotal += w
	}

	return maxTotal
}

// WordMaxWidth yields the width of the widest single word in the column.
func (c column) WordMaxWidth() int {
	var maxWidth int

	for _, cell := range c.cells {
		for i := range cell.wordLengths {
			maxWidth = max(maxWidth, cell.wordLengths[i][i])
		}
	}

	return maxWidth
}

// BreakLongestWords breaks words larger than limit in all this column's cells.
func (c *column) BreakLongestWords(wordBreakLevel breakLevel, limit int, splitter Splitter) {
	if wordBreakLevel == breakNone {
		return
	}

	for _, cell := range c.cells {
		if cellWidth(cell.content) <= limit {
			continue // unchanged cell
		}

		newLines := make([]string, 0, len(cell.content))

		for _, line := range cell.content { // iterate over the lines in this cell
			if displayWidth(line) <= limit {
				newLines = append(newLines, line) // unchanged line

				continue
			}

			wordsOnTheLine := newWords(strings.FieldsFunc(line, splitter))
			wordsOnTheLine.Sort() // widest word on the line comes first

			for _, word := range wordsOnTheLine {
				word.Break(limit, wordBreakLevel)
				if wordsOnTheLine.Width() <= limit { // enough word-breaking to abide by this limit
					break
				}
			}

			wordsOnTheLine.SortNatural() // return to the original ordering of words

			for _, word := range wordsOnTheLine {
				// parts := wrapMultiline(word.parts, limit, 0)
				// log.Printf("DEBUG: after rewrap (%d): => %v", limit, parts)
				newLines = append(newLines, word.parts...)
			}
		}
		// TODO: figure out how to break/not break when there are control chars
		// TODO: implement backtracking logic on word breaking - ugh!

		cell.content = newLines

		words := []string{}
		for _, line := range cell.content {
			words = append(words, strings.FieldsFunc(line, splitter)...)
		}

		cell.wordLengths, cell.maxWordLength, cell.minWordLength = buildLengthsMatrix(words, 1) // updates the matrix of the lengths of all word arrangements into single-line paragraphs [this is done once]
		cell.width = cellWidth(cell.content)
	}
}
