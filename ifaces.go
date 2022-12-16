package tablewriter

import (
	"github.com/logrusorgru/aurora/v4"
)

type (
	// CellWrapper knows how to wrap the content of a table cell into multiple lines.
	//
	// The wrapper knows about the display constraints.
	//
	// A few useful wrappers are provided by the package tablewrappers.
	CellWrapper interface {
		WrapCell(row, col int) []string
	}

	// Titler knows how to format an input string, suitable to display headings.
	Titler interface {
		Title(string) string
	}

	// CellWrapperFactory produces a cell wrapper with the knowledge of the table to be rendered.
	CellWrapperFactory func(*Table) CellWrapper

	// Formatter is a formatting function from the github.com/logrusorgru/aurora/v4 package,
	// to be used for nice terminal formatting such as colors, bold face, etc.
	//
	// It wraps some argument with an appropriate ANSI terminal escape sequence.
	Formatter = func(interface{}) aurora.Value
)
