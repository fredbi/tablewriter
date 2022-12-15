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

	for i, content := range t.header {
		lines := t.parseCell(content, i, headerRowIdx)
		t.headers = append(t.headers, lines)
	}

	for i, cells := range t.rows {
		var rowLines [][]string
		for j, v := range cells {
			rowLines = append(rowLines, t.parseCell(v, j, i))
		}
		t.lines = append(t.lines, rowLines)
	}

	for i, content := range t.footer {
		lines := t.parseCell(content, i, footerRowIdx)
		t.footers = append(t.footers, lines)
	}
}

func (t *Table) setWrapper() {
	if t.cellWrapperFactory != nil {
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

	// TODO: deprecate this - cellWrapper should be the unique interface
	if t.stringWrapperFactory != nil {
		wrapper := t.stringWrapperFactory(t)
		t.stringWrapper = wrapper.WrapString
	}
}

// parseCell analyzes a cell[col,row] of the table and computes its width and height.
// If wrapping is enabled, the content of the cell is wrapped.
//
// Works also for header and footer with special row indices.
func (t *Table) parseCell(str string, col, row int) []string { // TODO replace str by content cell address
	paragraphs := t.cellWrapper(row, col)

	maxWidth := min(t.colMaxWidth[col], wrap.CellWidth(paragraphs))
	t.setColWidth(col, maxWidth)
	t.setRowHeight(row, len(paragraphs))

	return paragraphs

	/*
		// previous implem with paragraph wrapping

		paragraphs := strings.FieldsFunc(str, wrap.LineSplitter)
		maxWidth := wrap.CellWidth(paragraphs)

		if t.stringWrapper != nil {
			maxWidth = min(t.colMaxWidth[col], maxWidth)
			maxWidth, paragraphs = t.wrapParagraphs(maxWidth, paragraphs)
		}

		t.setColWidth(col, maxWidth)
		t.setRowHeight(row, len(paragraphs))

		return paragraphs
	*/
}

// wrapParagraphs wraps the text inside a multi-lines cell and returns the new width and set of lines.
func (t *Table) wrapParagraphs(maxWidth int, paragraphs []string) (int, []string) {
	if t.reflowText {
		// make a single paragraph of everything.
		paragraphs = []string{strings.Join(paragraphs, SPACE)}
	}

	newMaxWidth := maxWidth
	wrappedParagraphs := make([]string, 0, len(paragraphs))

	for i, paragraph := range paragraphs {
		wrappedParagraph := t.stringWrapper(paragraph, maxWidth)
		newMaxWidth = max(wrap.CellWidth(wrappedParagraph), maxWidth)

		if i > 0 && !t.reflowText {
			// separate paragraphs with an empty line (there is no point if reflow is enabled)
			wrappedParagraphs = append(wrappedParagraphs, SPACE)
		}

		wrappedParagraphs = append(wrappedParagraphs, wrappedParagraph...)
	}

	return newMaxWidth, wrappedParagraphs
}
