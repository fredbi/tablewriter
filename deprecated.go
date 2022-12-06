package tablewriter

/*
import "io"

// @deprecated
func NewWriter(writer io.Writer, opts ...Option) *Table {
	t := &Table{
		lines:      [][][]string{},
		headers:    [][]string{},
		footers:    [][]string{},
		numColumns: -1,

		options: defaultOptions(opts),
	}
	t.out = writer // TODO(fred)

	return t
}

// Append row to table with color attributes
// TODO(fred): remove call to parseDimension
func (t *Table) Rich(row []string, colors []Colors) {
	rowSize := len(t.headers)
	if rowSize > t.numColumns {
		t.numColumns = rowSize
	}

	n := len(t.lines)
	line := [][]string{}
	for i, v := range row {

		// Detect string  width
		// Detect String height
		// Break strings into words
		out := t.parseDimension(v, i, n)

		if len(colors) > i {
			color := colors[i]
			out[0] = format(out[0], color)
		}

		// Append broken words
		line = append(line, out)
	}
	t.lines = append(t.lines, line)
}

// @deprecated
func (t *Table) AppendBulk(rows [][]string) {
	for _, row := range rows {
		t.Append(row)
	}
}
*/
