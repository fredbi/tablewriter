package tablewrappers

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMatrixInternals(t *testing.T) {
	rows, columns := buildMatrix(
		[][]string{
			{"r11xxx", "r12", "r13"},
			{"r21", "r22x", "r23"},
			{"r31", "r32", "r33"},
			{"r41", "r42", "r43xx"},
		},
		BlankSplitter,
	)

	// spew.Dump(rows)
	// spew.Dump(columns)

	t.Logf("%v", rows[0].Values())
	t.Logf("%v", columns[0].Values())

	columns.SortRows()
	columns.Sort()
	spew.Dump(columns)
}
