package tablewrappers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRowWrapper(t *testing.T) {
	matrix := [][]string{
		{"Heading A", "Heading BCD", "Another heading, longer than the other ones"},
		{"elem 10 - short", "elem 11 - longer than the previous one", "elem 12"},
		{"elem 20 - not super long", "elem 21 - small", "elem 22 - intermediate"},
		{"Footer A", "Footer BCD", "Another footer, much longer than the other ones"},
	}
	const (
		terminalWidth = 40
		numRows       = 4
		numCols       = 3
	)
	w := NewRowWrapper(matrix, terminalWidth)

	originalCols := make([][]string, numCols)
	for col := 0; col < numCols; col++ {
		originalCols[col] = make([]string, 0, numRows)
		for row := 0; row < numRows; row++ {
			originalCols[col] = append(originalCols[col], matrix[row][col])
		}
	}

	expectedMaxWidths := []int{9, 11, 12} // 9 + 11 +12 = 32 < 40
	expectedWrapped := [][][]string{
		{{"Heading A"}, {"Heading BCD"}, {"Another", "heading,", "longer than", "the other", "ones"}},
		{{"elem 10 -", "short"}, {"elem 11", "- longer", "than the", "previous", "one"}, {"elem 12"}},
		{{"elem 20 -", "not super", "long"}, {"elem 21 -", "small"}, {"elem 22 -", "intermediate"}},
		{{"Footer A"}, {"Footer BCD"}, {"Another", "footer, much", "longer than", "the other", "ones"}},
	}

	for row := 0; row < numRows; row++ {
		require.LessOrEqual(t, w.columns[0].rows[row].TotalWidth(), terminalWidth)

		for col := 0; col < numCols; col++ {
			// sanity check
			require.LessOrEqualf(t,
				w.columns[col].maxWidth,
				cellWidth(originalCols[col]),
				"column max width %d should be less or equal than original width %d for column (%#v)",
				w.columns[col].maxWidth,
				cellWidth(originalCols[col]),
				originalCols[col],
			)

			require.Equal(t, expectedMaxWidths[col], w.columns[col].maxWidth)

			wrapped := w.WrapCell(row, col)
			require.Equal(t, expectedWrapped[row][col], wrapped)

			/*
				t.Logf("[%d,%d] -> %d\n%q -> \n%s",
					row, col,
					w.columns[col].maxWidth,
					w.matrix[row][col], strings.Join(wrapped, "\n"),
				)
			*/
		}
	}
}
