package tablewriter

type (
	// StringWrapper knows how to wrap a string into multiple lines,
	// under the constraint of a maximum display width.
	StringWrapper interface {
		WrapString(input string, maxWidth int) []string
	}

	// CellWrapper knows how to wrap the content of a cell into multiple lines,
	// under the constraint of a maximum display width.
	CellWrapper interface {
		WrapCell(row, col int) []string
	}

	// Titler knows how to format an input string, suitable to display headings.
	Titler interface {
		Title(string) string
	}

	CellWrapperFactory func(*Table) CellWrapper
)
