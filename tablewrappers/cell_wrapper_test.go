package tablewrappers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRowWrapper(t *testing.T) {
	t.Run("should wrap table within global width constraint (40)", func(t *testing.T) {
		const (
			terminalWidth = 40
			numRows       = 4
			numCols       = 3
		)

		matrix := [][]string{
			{"Heading A", "Heading BCD", "Another heading, longer than the other ones"},
			{"elem 10 - short", "elem 11 - longer than the previous one", "elem 12"},
			{"elem 20 - not super long", "elem 21 - small", "elem 22 - intermediate"},
			{"Footer A", "Footer BCD", "Another footer, much longer than the other ones"},
		}

		expectedMaxWidths := []int{9, 9, 21} // 9 + 9 +21 = 39 < 40
		expectedWrapped := [][][]string{
			{{"Heading A"}, {"Heading", "BCD"}, {"Another heading,", "longer than the other", "ones"}},
			{{"elem 10 -", "short"}, {"elem 11", "- longer", "than the", "previous", "one"}, {"elem 12"}},
			{{"elem 20 -", "not super", "long"}, {"elem 21 -", "small"}, {"elem 22 -", "intermediate"}},
			{{"Footer A"}, {"Footer", "BCD"}, {"Another footer, much", "longer than the other", "ones"}},
		}

		w := NewRowWrapper(matrix, terminalWidth)
		originalCols := make([][]string, numCols)
		for col := 0; col < numCols; col++ {
			originalCols[col] = make([]string, 0, numRows)
			for row := 0; row < numRows; row++ {
				originalCols[col] = append(originalCols[col], matrix[row][col])
			}
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
			}
		}
	})

	t.Run("should wrap table within global width constraint (18)", func(t *testing.T) {
		const (
			terminalWidth = 18 // (30 - 12 overhead)
			numRows       = 8
			numCols       = 3
		)

		matrix := [][]string{
			{"Date", "Name", "Items", "Price"},
			{"", "", "Total", "$145.93"},
			{"1/1/2014", "Domain name", "2233", "$10.98"},
			{"1/1/2014", "January Hosting", "2233", "$54.95"},
			{"", "    (empty)\n    (empty)", "", ""},
			{"1/4/2014", "February Hosting", "2233", "$51.00"},
			{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
			{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
		}
		expected := [][][]string{
			{{"Date"}, {"Name"}, {"Items"}},
			{{""}, {""}, {"Total"}},
			{{"1/1/2014"}, {"Domain", "name"}, {"2233"}},
			{{"1/1/2014"}, {"January", "Hosting"}, {"2233"}},
			{{""}, {"(empty)", "(empty)"}, {""}},
			{{"1/4/2014"}, {"February", "Hosting"}, {"2233"}},
			{{"1/4/2014"}, {"February", "Extra", "Bandwidth"}, {"2233"}},
			{{"1/4/2014"}, {"(Discount)"}, {"2233"}},
		}

		w := NewRowWrapper(matrix, terminalWidth)
		for row := 0; row < numRows; row++ {
			for col := 0; col < numCols; col++ {
				wrapped := w.WrapCell(row, col)
				require.EqualValues(t, expected[row][col], wrapped)
			}
		}
	})
}

func TestRowWrapperWordBreak(t *testing.T) {
	const (
		terminalWidth = 10
		numRows       = 8
		numCols       = 3
	)

	matrix := [][]string{
		{"Date", "Name", "Items", "Price"},
		{"", "", "Total", "$145.93"},
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"", "    (empty)\n    (empty)", "", ""},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
	}
	expected := [][][]string{
		{{"Date"}, {"Name"}, {"Items"}},
		{{""}, {""}, {"Total"}},
		{{"1/1/2014"}, {"Domain", "name"}, {"2233"}},
		{{"1/1/2014"}, {"January", "Hosting"}, {"2233"}},
		{{""}, {"(empty)", "(empty)"}, {""}},
		{{"1/4/2014"}, {"February", "Hosting"}, {"2233"}},
		{{"1/4/2014"}, {"February", "Extra", "Bandwidth"}, {"2233"}},
		{{"1/4/2014"}, {"(Discount)"}, {"2233"}},
	}

	w := NewRowWrapper(matrix, terminalWidth)
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			wrapped := w.WrapCell(row, col)
			t.Logf("[%d][%d] %q => %v", row, col, matrix[row][col], wrapped)
			require.EqualValues(t, expected[row][col], wrapped)
		}
	}
}
