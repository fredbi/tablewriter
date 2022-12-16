package tablewriter

import (
	"strings"

	wrap "github.com/fredbi/tablewriter/tablewrappers"
)

func (t *Table) prepare() {
	t.setNumColumns()
	t.fillAlignments()
	t.fillMaxWidths()

	// evaluate wrapped content
	t.setWrapper()

	for i := range t.header {
		lines := t.parseCell(i, headerRowIdx)
		t.headers = append(t.headers, lines)
	}

	for i, cells := range t.rows {
		var rowLines [][]string
		for j := range cells {
			rowLines = append(rowLines, t.parseCell(j, i))
		}
		t.lines = append(t.lines, rowLines)
	}

	for i := range t.footer {
		lines := t.parseCell(i, footerRowIdx)
		t.footers = append(t.footers, lines)
	}
}

func (t *Table) setWrapper() {
	if t.cellWrapperFactory != nil {
		// wrap is enabled with some wrapper
		wrapper := t.cellWrapperFactory(t)
		t.cellWrapper = func(row, col int) []string {
			rowOffset := row

			switch {
			case row == headerRowIdx:
				rowOffset = 0
			case row == footerRowIdx:
				rowOffset = 1
			default:
				if len(t.header) > 0 {
					rowOffset++
				}
				if len(t.footer) > 0 {
					rowOffset++
				}
			}

			return wrapper.WrapCell(rowOffset, col)
		}

		return
	}

	// wrap is disabled: set a noop wrapper. This preserves blank space and paragraphs.
	paragrapher := func(s string) []string { return strings.FieldsFunc(s, wrap.LineSplitter) }
	t.cellWrapper = func(row, col int) []string {
		switch {
		case row == headerRowIdx:
			return paragrapher(t.header[col])
		case row == footerRowIdx:
			return paragrapher(t.footer[col])
		default:
			return paragrapher(t.rows[row][col])
		}
	}
}

// parseCell analyzes a cell[col,row] of the table and computes its width and height.
// If wrapping is enabled, the content of the cell is wrapped.
//
// Works also for header and footer with special row indices.
func (t *Table) parseCell(col, row int) []string {
	paragraphs := t.cellWrapper(row, col)

	t.setColWidth(col, wrap.CellWidth(paragraphs))
	t.setRowHeight(row, len(paragraphs))

	return paragraphs
}

// defaultCellWrapperFactory provides a cell wrapper that abides by column-width constraints.
func defaultCellWrapperFactory() CellWrapperFactory {
	return func(t *Table) CellWrapper {
		return wrap.NewDefaultCellWrapper(
			makeMatrix(t),
			t.ColLimits(),
		)
	}
}

// rowCellWrapperFactory provides a cell wrapper that abide by a single table-width constraint.
func rowCellWrapperFactory(width int) CellWrapperFactory {
	return func(t *Table) CellWrapper {
		wrapper := wrap.NewRowWrapper(makeMatrix(t), width-t.Overhead())

		return wrapper
	}
}
