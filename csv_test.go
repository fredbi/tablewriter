package tablewriter

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSVInfo(t *testing.T) {
	file, err := os.Open("testdata/test_info.csv")
	require.NoError(t, err)

	reader := csv.NewReader(file)
	var buf bytes.Buffer

	table, err := NewCSV(reader, true,
		WithCellAlignment(AlignLeft),
		WithAllBorders(false),
		WithWriter(&buf),
	)
	require.NoError(t, err)

	table.Render()

	const want = `   FIELD   |     TYPE     | NULL | KEY | DEFAULT |     EXTRA
-----------+--------------+------+-----+---------+-----------------
  user_id  | smallint(5)  | NO   | PRI | NULL    | auto_increment
  username | varchar(10)  | NO   |     | NULL    |
  password | varchar(100) | NO   |     | NULL    |
`

	t.Run("should not right-pad with blanks when no borded is rendered", func(t *testing.T) {
		checkEqual(t, buf.String(), want, "CSV info failed")
	})
}

func TestCSVSeparator(t *testing.T) {
	file, err := os.Open("testdata/test.csv")
	require.NoError(t, err)

	reader := csv.NewReader(file)
	var buf bytes.Buffer

	table, err := NewCSV(reader, true,
		WithRowLine(true),
		WithCenterSeparator("+"),
		WithColumnSeparator("|"),
		WithRowSeparator("-"),
		WithCellAlignment(AlignLeft),
		WithWriter(&buf),
	)
	require.NoError(t, err)
	table.Render()

	const want = `+------------+-----------+---------+
| FIRST NAME | LAST NAME |   SSN   |
+------------+-----------+---------+
| John       | Barry     | 123456  |
+------------+-----------+---------+
| Kathy      | Smith     | 687987  |
+------------+-----------+---------+
| Bob        | McCornick | 3979870 |
+------------+-----------+---------+
`

	checkEqual(t, buf.String(), want, "CSV info failed")
}
