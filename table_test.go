// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestInternals(t *testing.T) {
	t.Run("number of lines should match rows", func(t *testing.T) {
		data := [][]string{
			{"A", "The Good", "500"},
			{"B", "The Very very Bad Man", "288"},
			{"C", "The Ugly", "120"},
			{"D", "The Gopher", "800"},
		}

		buf := &bytes.Buffer{}
		table := New(
			WithWriter(buf),
			WithHeader([]string{"Name", "Sign", "Rating"}),
			WithRows(data),
		)
		table.prepare()

		checkEqual(t, len(table.lines), len(data), "Number of lines failed")
	})
}

func TestNoBorder(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"", "    (empty)\n    (empty)", "", ""},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
	}

	t.Run("should render with footer", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
			WithFooter([]string{"", "", "Total", "$145.93"}),
			WithRows(data),
			WithWrap(false),
			WithAllBorders(false),
		)
		table.Render()
		const want = `    DATE   |       DESCRIPTION        |  CV2  | AMOUNT
-----------+--------------------------+-------+----------
  1/1/2014 | Domain name              |  2233 | $10.98
  1/1/2014 | January Hosting          |  2233 | $54.95
           |     (empty)              |       |
           |     (empty)              |       |
  1/4/2014 | February Hosting         |  2233 | $51.00
  1/4/2014 | February Extra Bandwidth |  2233 | $30.00
  1/4/2014 |     (Discount)           |  2233 | -$1.00
-----------+--------------------------+-------+----------
                                        TOTAL | $145.93
                                      --------+----------
`

		checkEqual(t, buf.String(), want, "border table rendering failed")
	})

	t.Run("should render without footer", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
			WithRows(data),
			WithWrap(false),
			WithAllBorders(false),
		)
		table.Render()

		const want = `    DATE   |       DESCRIPTION        | CV2  | AMOUNT
-----------+--------------------------+------+---------
  1/1/2014 | Domain name              | 2233 | $10.98
  1/1/2014 | January Hosting          | 2233 | $54.95
           |     (empty)              |      |
           |     (empty)              |      |
  1/4/2014 | February Hosting         | 2233 | $51.00
  1/4/2014 | February Extra Bandwidth | 2233 | $30.00
  1/4/2014 |     (Discount)           | 2233 | -$1.00
`

		checkEqual(t, buf.String(), want, "border table rendering failed")
	})
}

func TestWithBorder(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"", "    (empty)\n    (empty)", "", ""},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		{"1/4/2014", "    (Discount)", "2233", "-$1.00"},
	}

	table, buf := NewBuffered(
		WithWrap(false),
		WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
		WithFooter([]string{"", "", "Total", "$145.93"}),
		WithRows(data),
	)
	table.Render()

	want := `+----------+--------------------------+-------+---------+
|   DATE   |       DESCRIPTION        |  CV2  | AMOUNT  |
+----------+--------------------------+-------+---------+
| 1/1/2014 | Domain name              |  2233 | $10.98  |
| 1/1/2014 | January Hosting          |  2233 | $54.95  |
|          |     (empty)              |       |         |
|          |     (empty)              |       |         |
| 1/4/2014 | February Hosting         |  2233 | $51.00  |
| 1/4/2014 | February Extra Bandwidth |  2233 | $30.00  |
| 1/4/2014 |     (Discount)           |  2233 | -$1.00  |
+----------+--------------------------+-------+---------+
|                                       TOTAL | $145.93 |
+----------+--------------------------+-------+---------+
`

	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestPrintingInMarkdown(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
	}

	table, buf := NewBuffered(
		WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
		WithMarkdown(true),
		// WithBorders(Border{Left: true, Top: false, Right: true, Bottom: false}),
		// WithCenterSeparator("|"),
		WithRows(data),
	)
	table.Render()

	want := `|   DATE   |       DESCRIPTION        | CV2  | AMOUNT |
|----------|--------------------------|------|--------|
| 1/1/2014 | Domain name              | 2233 | $10.98 |
| 1/1/2014 | January Hosting          | 2233 | $54.95 |
| 1/4/2014 | February Hosting         | 2233 | $51.00 |
| 1/4/2014 | February Extra Bandwidth | 2233 | $30.00 |
`
	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestTitleCase(t *testing.T) {
	line := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"}
	const (
		titledCols = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C |\n"
		cols       = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |\n"
		seps       = "+---+---+---+---+---+---+---+---+---+---+---+---+\n"
	)

	t.Run("should print header with Title-case", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader(line),
		)
		table.prepare()
		table.printHeader()

		checkEqual(t, buf.String(), titledCols+seps, "header rendering failed")
	})

	t.Run("should print header without alteration", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader(line),
			WithTitledHeader(false),
		)
		table.prepare()
		table.printHeader()

		checkEqual(t, buf.String(), cols+seps, "header rendering failed")
	})
}

func TestPrintFooter(t *testing.T) {
	line := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"}
	const (
		data = "  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c\n"
	)

	t.Run("with borders", func(t *testing.T) {
		const (
			cols            = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C |\n"
			seps            = "+---+---+---+---+---+---+---+---+---+---+---+---+\n"
			dataWithBorders = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |\n"
		)
		options := []Option{
			WithHeader(line),
			WithFooter(line),
		}

		t.Run("shoud render footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printFooter()
			checkEqual(t, buf.String(), cols+seps, "footer rendering failed")
		})

		t.Run("shoud render header the same as footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printHeader()
			checkEqual(t, buf.String(), cols+seps, "header rendering failed")
		})

		t.Run("shoud render data the same", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), seps+cols+seps+dataWithBorders+seps+cols+seps, "rendering failed")
		})
	})

	t.Run("without borders", func(t *testing.T) {
		const (
			cols = "  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C\n"
			seps = "----+---+---+---+---+---+---+---+---+---+---+----\n"
		)

		options := []Option{
			WithHeader(line),
			WithFooter(line),
			WithAllBorders(false),
		}

		t.Run("shoud render footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printFooter()
			checkEqual(t, buf.String(), seps+cols+seps, "footer rendering failed")
		})

		t.Run("shoud render header the same as footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printHeader()
			checkEqual(t, buf.String(), cols+seps, "header rendering failed")
		})
	})

	t.Run("with bottom border", func(t *testing.T) {
		const (
			cols = "  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C\n"
			seps = "----+---+---+---+---+---+---+---+---+---+---+----\n"
		)
		options := []Option{
			WithHeader(line),
			WithFooter(line),
			WithBorders(Border{Left: false, Top: false, Right: false, Bottom: true}),
		}

		t.Run("shoud render footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printFooter()
			checkEqual(t, buf.String(), cols+seps, "footer rendering failed")
		})

		t.Run("shoud render header the same as footer", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.prepare()
			table.printHeader()
			checkEqual(t, buf.String(), cols+seps, "header rendering failed")
		})

		t.Run("shoud render data line the same as header", func(t *testing.T) {
			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), cols+seps+data+seps+cols+seps, "rendering failed")
		})
	})

	t.Run("with empty headings", func(t *testing.T) {
		incompleteLine := []string{"", "2", "", "4", "", "6", "7", "8", "9", "a", "b", ""}

		t.Run("with borders, shoud exhibit differences between header and footer", func(t *testing.T) {
			const (
				cols  = "|   | 2 |   | 4 |   | 6 | 7 | 8 | 9 | A | B |   |\n"
				data  = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |\n"
				seps  = "+---+---+---+---+---+---+---+---+---+---+---+---+\n"
				fcols = "|     2 |     4 |     6 | 7 | 8 | 9 | A | B |    \n"
			)

			options := []Option{
				WithHeader(incompleteLine),
				WithFooter(incompleteLine),
				WithAllBorders(true),
			}

			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), seps+cols+seps+data+seps+fcols+seps, "rendering failed")
		})

		t.Run("without borders, shoud exhibit differences between header and footer", func(t *testing.T) {
			const (
				cols  = "    | 2 |   | 4 |   | 6 | 7 | 8 | 9 | A | B |\n"
				data  = "  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c\n"
				seps  = "----+---+---+---+---+---+---+---+---+---+---+----\n"
				fcols = "      2 |     4 |     6 | 7 | 8 | 9 | A | B |\n"
				fseps = "    ----+---+---+---+---+---+---+---+---+---+----\n"
			)

			options := []Option{
				WithHeader(incompleteLine),
				WithFooter(incompleteLine),
				WithAllBorders(false),
			}

			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), cols+seps+data+seps+fcols+fseps, "rendering failed")
		})

		t.Run("with left border, shoud exhibit differences between header and footer", func(t *testing.T) {
			const (
				cols  = "|   | 2 |   | 4 |   | 6 | 7 | 8 | 9 | A | B |\n"
				data  = "| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c\n"
				seps  = "+---+---+---+---+---+---+---+---+---+---+---+----\n"
				fcols = "|     2 |     4 |     6 | 7 | 8 | 9 | A | B |\n"
				fseps = " ---+---+---+---+---+---+---+---+---+---+---+----\n"
			)

			options := []Option{
				WithHeader(incompleteLine),
				WithFooter(incompleteLine),
				WithBorders(Border{Left: true, Top: false, Right: false, Bottom: false}),
			}

			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), cols+seps+data+seps+fcols+fseps, "rendering failed")
		})

		t.Run("with right border, shoud exhibit differences between header and footer", func(t *testing.T) {
			const (
				cols  = "    | 2 |   | 4 |   | 6 | 7 | 8 | 9 | A | B |   |\n"
				data  = "  1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |\n"
				seps  = "----+---+---+---+---+---+---+---+---+---+---+---+\n"
				fcols = "      2 |     4 |     6 | 7 | 8 | 9 | A | B |    \n"
				fseps = "+   +---+---+---+---+---+---+---+---+---+---+---+\n" // ICI
			)

			options := []Option{
				WithHeader(incompleteLine),
				WithFooter(incompleteLine),
				WithBorders(Border{Left: false, Top: false, Right: true, Bottom: false}),
			}

			table, buf := NewBuffered(options...)
			table.Append(line)
			table.Render()
			checkEqual(t, buf.String(), cols+seps+data+seps+fcols+fseps, "rendering failed")
		})
	})
	t.Run("with no whitespace", func(t *testing.T) {
		// TODO
	})
}

func TestPrintFooterUntitled(t *testing.T) {
	table, buf := NewBuffered(
		WithTitledHeader(false),
		WithHeader([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"}),
		WithFooter([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c"}),
	)
	table.prepare()
	table.printFooter()

	const want = `| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | a | b | c |
+---+---+---+---+---+---+---+---+---+---+---+---+
`
	checkEqual(t, buf.String(), want, "footer rendering failed")
}

func TestPrintShortCaption(t *testing.T) {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table, buf := NewBuffered(
		WithHeader([]string{"Name", "Sign", "Rating"}),
		WithCaption("Short caption."),
		WithRows(data),
	)
	table.Render()

	const want = `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
| B    | The Very very Bad Man |    288 |
| C    | The Ugly              |    120 |
| D    | The Gopher            |    800 |
+------+-----------------------+--------+
Short caption.
`
	checkEqual(t, buf.String(), want, "long caption for short example rendering failed")
}

func TestPrintLongCaptionWithShortExample(t *testing.T) {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table, buf := NewBuffered(
		WithHeader([]string{"Name", "Sign", "Rating"}),
		WithCaption("This is a very long caption. The text should wrap. If not, we have a problem that needs to be solved."),
		WithRows(data),
	)
	table.Render()

	const want = `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
| B    | The Very very Bad Man |    288 |
| C    | The Ugly              |    120 |
| D    | The Gopher            |    800 |
+------+-----------------------+--------+
This is a very long caption. The text
should wrap. If not, we have a problem
that needs to be solved.
`
	checkEqual(t, buf.String(), want, "long caption for short example rendering failed")
}

func TestPrintCaptionWithFooter(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
		{"1/1/2014", "January Hosting", "2233", "$54.95"},
		{"1/4/2014", "February Hosting", "2233", "$51.00"},
		{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
	}

	table, buf := NewBuffered(
		WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
		WithFooter([]string{"", "", "Total", "$146.93"}),
		WithCaption("This is a very long caption. The text should wrap to the width of the table."),
		WithAllBorders(false),
		WithRows(data),
	)
	table.Render()

	const want = `    DATE   |       DESCRIPTION        |  CV2  | AMOUNT
-----------+--------------------------+-------+----------
  1/1/2014 | Domain name              |  2233 | $10.98
  1/1/2014 | January Hosting          |  2233 | $54.95
  1/4/2014 | February Hosting         |  2233 | $51.00
  1/4/2014 | February Extra Bandwidth |  2233 | $30.00
-----------+--------------------------+-------+----------
                                        TOTAL | $146.93
                                      --------+----------
This is a very long caption. The text should wrap to the
width of the table.
`
	checkEqual(t, buf.String(), want, "border table rendering failed")
}

func TestPrintLongCaptionWithLongExample(t *testing.T) {
	header := []string{"Name", "Sign", "Rating"}
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "Some Data", "Another Data"},
		{"Instead of lining up the letters all", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards"},
	}
	const (
		expectedTable = `+--------------------------------+--------------------------------+-------------------------------+
|              NAME              |              SIGN              |            RATING             |
+--------------------------------+--------------------------------+-------------------------------+
| Learn East has computers       | Some Data                      | Another Data                  |
| with adapted keyboards with    |                                |                               |
| enlarged print etc             |                                |                               |
| Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards |
| letters all                    | keyboard in two                |                               |
+--------------------------------+--------------------------------+-------------------------------+
`
		expectedCaption = `This is a very long caption. The text should wrap. If not, we have a problem that needs to be
solved.
`
	)

	t.Run("should render caption", func(t *testing.T) {
		table, buf := NewBuffered(
			WithCaption("This is a very long caption. The text should wrap. If not, we have a problem that needs to be solved."),
			WithHeader(header),
			WithRows(data),
		)
		table.Render()

		const expected = expectedTable + expectedCaption
		checkEqual(t, buf.String(), expected, "long caption for long example rendering failed")
	})

	t.Run("should wrap first col only", func(t *testing.T) {
		table, buf := NewBuffered(
			WithWrap(true),
			WithHeader(header),
			WithColWidth(50),
			WithColMaxWidth(0, 10),
			WithRows(data),
		)
		table.Render()

		const (
			expected = `+------------+-----------------------------------------------+-------------------------------+
|    NAME    |                     SIGN                      |            RATING             |
+------------+-----------------------------------------------+-------------------------------+
| Learn      | Some Data                                     | Another Data                  |
| East has   |                                               |                               |
| computers  |                                               |                               |
| with       |                                               |                               |
| adapted    |                                               |                               |
| keyboards  |                                               |                               |
| with       |                                               |                               |
| enlarged   |                                               |                               |
| print etc  |                                               |                               |
| Instead    | the way across, he splits the keyboard in two | Like most ergonomic keyboards |
| of lining  |                                               |                               |
| up the     |                                               |                               |
| letters    |                                               |                               |
| all        |                                               |                               |
+------------+-----------------------------------------------+-------------------------------+
`
		)

		checkEqual(t, buf.String(), expected)
	})

	t.Run("should wrap first and third cols only", func(t *testing.T) {
		table, buf := NewBuffered(
			WithWrap(true),
			WithHeader(header),
			WithColWidth(50),
			WithColMaxWidths(map[int]int{0: 10, 2: 12}),
		)

		for _, v := range data {
			table.Append(v)
		}
		table.Render()

		const (
			expected = `+------------+-----------------------------------------------+--------------+
|    NAME    |                     SIGN                      |    RATING    |
+------------+-----------------------------------------------+--------------+
| Learn      | Some Data                                     | Another Data |
| East has   |                                               |              |
| computers  |                                               |              |
| with       |                                               |              |
| adapted    |                                               |              |
| keyboards  |                                               |              |
| with       |                                               |              |
| enlarged   |                                               |              |
| print etc  |                                               |              |
| Instead    | the way across, he splits the keyboard in two | Like most    |
| of lining  |                                               | ergonomic    |
| up the     |                                               | keyboards    |
| letters    |                                               |              |
| all        |                                               |              |
+------------+-----------------------------------------------+--------------+
`
		)

		checkEqual(t, buf.String(), expected)
	})
}

func TestPrintSepLine(t *testing.T) {
	header := make([]string, 12)
	val := " "
	want := ""
	for i := range header {
		header[i] = val
		want = fmt.Sprintf("%s+-%s-", want, strings.ReplaceAll(val, " ", "-"))
		val += " "
	}

	want += "+"
	var buf bytes.Buffer
	table := New(
		WithWriter(&buf),
		WithHeader(header),
	)
	table.prepare()
	table.printSepLine(false)
	checkEqual(t, buf.String(), want, "line rendering failed")
}

func TestAnsiStrip(t *testing.T) {
	header := make([]string, 12)
	val := " "
	want := ""
	for i := range header {
		header[i] = "\033[43;30m" + val + "\033[00m"
		want = fmt.Sprintf("%s+-%s-", want, strings.ReplaceAll(val, " ", "-"))
		val += " "
	}
	want += "+"
	var buf bytes.Buffer
	table := New(
		WithWriter(&buf),
		WithHeader(header),
	)
	table.prepare()
	table.printSepLine(false)
	checkEqual(t, buf.String(), want, "line rendering failed")
}

func TestSubclass(t *testing.T) {
	buf := new(bytes.Buffer)
	table := newCustomizedTable(buf)

	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	want := `  A  The Good               500
  B  The Very very Bad Man  288
  C  The Ugly               120
  D  The Gopher             800
`
	checkEqual(t, buf.String(), want, "test subclass failed")
}

func TestAutoMergeRows(t *testing.T) {
	data := [][]string{
		{"A", "The Good", "500"},
		{"A", "The Very very Bad Man", "288"},
		{"B", "The Very very Bad Man", "120"},
		{"B", "The Very very Bad Man", "200"},
	}

	t.Run("should render merged cells", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader([]string{"Name", "Sign", "Rating"}),
			WithRows(data),
			WithMergeCells(true),
		)

		table.Render()

		const want = `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
|      | The Very very Bad Man |    288 |
| B    |                       |    120 |
|      |                       |    200 |
+------+-----------------------+--------+
`
		checkEqual(t, buf.String(), want, "failed to render merged cells")
	})

	t.Run("should render merged cells with row separators", func(t *testing.T) {
		table, buf := NewBuffered(
			WithHeader([]string{"Name", "Sign", "Rating"}),
			WithRows(data),
			WithMergeCells(true),
			WithRowLine(true),
		)

		table.Render()

		const want = `+------+-----------------------+--------+
| NAME |         SIGN          | RATING |
+------+-----------------------+--------+
| A    | The Good              |    500 |
+      +-----------------------+--------+
|      | The Very very Bad Man |    288 |
+------+                       +--------+
| B    |                       |    120 |
+      +                       +--------+
|      |                       |    200 |
+------+-----------------------+--------+
`
		checkEqual(t, buf.String(), want)
	})

	t.Run("should render merged cells with more row separators", func(t *testing.T) {
		dataWithlongText := [][]string{
			{"A", "The Good", "500"},
			{"A", "The Very very very very very Bad Man", "288"},
			{"B", "The Very very very very very Bad Man", "120"},
			{"C", "The Very very Bad Man", "200"},
		}

		table, buf := NewBuffered(
			WithHeader([]string{"Name", "Sign", "Rating"}),
			WithRows(dataWithlongText),
			WithMergeCells(true),
			WithRowLine(true),
		)

		table.Render()

		const want = `+------+--------------------------------+--------+
| NAME |              SIGN              | RATING |
+------+--------------------------------+--------+
| A    | The Good                       |    500 |
+      +--------------------------------+--------+
|      | The Very very very very very   |    288 |
|      | Bad Man                        |        |
+------+                                +--------+
| B    |                                |    120 |
|      |                                |        |
+------+--------------------------------+--------+
| C    | The Very very Bad Man          |    200 |
+------+--------------------------------+--------+
`
		checkEqual(t, buf.String(), want)
	})

	t.Run("should render merged cells with more row separators", func(t *testing.T) {
		dataWithlongText2 := [][]string{
			{"A", "The Good", "500"},
			{"A", "The Very very very very very Bad Man", "288"},
			{"B", "The Very very Bad Man", "120"},
		}

		table, buf := NewBuffered(
			WithHeader([]string{"Name", "Sign", "Rating"}),
			WithRows(dataWithlongText2),
			WithMergeCells(true),
			WithRowLine(true),
		)

		table.Render()

		const want = `+------+--------------------------------+--------+
| NAME |              SIGN              | RATING |
+------+--------------------------------+--------+
| A    | The Good                       |    500 |
+      +--------------------------------+--------+
|      | The Very very very very very   |    288 |
|      | Bad Man                        |        |
+------+--------------------------------+--------+
| B    | The Very very Bad Man          |    120 |
+------+--------------------------------+--------+
`
		checkEqual(t, buf.String(), want)
	})
}

func TestClearRows(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
	}

	t.Run("should render without wrapping", func(t *testing.T) {
		var buf bytes.Buffer
		table := New(
			WithWriter(&buf),
			WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
			WithFooter([]string{"", "", "Total", "$145.93"}),
			WithWrap(false),
			WithRows(data),
		)
		table.Render()

		const originalWant = `+----------+-------------+-------+---------+
|   DATE   | DESCRIPTION |  CV2  | AMOUNT  |
+----------+-------------+-------+---------+
| 1/1/2014 | Domain name |  2233 | $10.98  |
+----------+-------------+-------+---------+
|                          TOTAL | $145.93 |
+----------+-------------+-------+---------+
`
		want := originalWant

		checkEqual(t, buf.String(), want, "table clear rows failed")
	})

	t.Run("should render without inner rows", func(t *testing.T) {
		var buf bytes.Buffer
		table := New(
			WithWriter(&buf),
			WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
			WithFooter([]string{"", "", "Total", "$145.93"}),
			WithWrap(false),
		)
		table.Render()

		const want = `+------+-------------+-------+---------+
| DATE | DESCRIPTION |  CV2  | AMOUNT  |
+------+-------------+-------+---------+
+------+-------------+-------+---------+
|                      TOTAL | $145.93 |
+------+-------------+-------+---------+
`

		checkEqual(t, buf.String(), want, "table clear rows failed")
	})

	t.Run("should render without footer", func(t *testing.T) {
		var buf bytes.Buffer
		table := New(
			WithWriter(&buf),
			WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
			WithWrap(false),
			WithRows(data),
		)
		table.Render()

		const want = `+----------+-------------+------+--------+
|   DATE   | DESCRIPTION | CV2  | AMOUNT |
+----------+-------------+------+--------+
| 1/1/2014 | Domain name | 2233 | $10.98 |
+----------+-------------+------+--------+
`

		checkEqual(t, buf.String(), want, "table clear rows failed")
	})
}

func TestNoWrap(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "Domain name", "2233", "$10.98"},
	}

	table, buf := NewBuffered(
		WithWrap(false),
		WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
		WithRows(data),
	)
	table.Render()

	want := `+----------+-------------+------+--------+
|   DATE   | DESCRIPTION | CV2  | AMOUNT |
+----------+-------------+------+--------+
| 1/1/2014 | Domain name | 2233 | $10.98 |
+----------+-------------+------+--------+
`

	checkEqual(t, buf.String(), want)
}

func TestMoreDataColumnsThanHeaders(t *testing.T) {
	var (
		header = []string{"A", "B", "C"}
		data   = [][]string{
			{"a", "b", "c", "d"},
			{"1", "2", "3", "4"},
		}
	)

	const want = `+---+---+---+---+
| A | B | C |   |
+---+---+---+---+
| a | b | c | d |
| 1 | 2 | 3 | 4 |
+---+---+---+---+
`

	table, buf := NewBuffered(
		WithHeader(header),
		WithRows(data),
	)

	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestMoreFooterColumnsThanHeaders(t *testing.T) {
	var (
		header = []string{"A", "B", "C"}
		data   = [][]string{
			{"a", "b", "c", "d"},
			{"1", "2", "3", "4"},
		}
		footer = []string{"a", "b", "c", "d", "e"}
	)

	const want = `+---+---+---+---+---+
| A | B | C |   |   |
+---+---+---+---+---+
| a | b | c | d |   |
| 1 | 2 | 3 | 4 |   |
+---+---+---+---+---+
| A | B | C | D | E |
+---+---+---+---+---+
`

	table, buf := NewBuffered(
		WithHeader(header),
		WithFooter(footer),
		WithRows(data),
	)
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestSetColMinWidth(t *testing.T) {
	var (
		header = []string{"AAA", "BBB", "CCC"}
		data   = [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}
		footer = []string{"a", "b", "cccc"}
	)

	const want = `+-----+-----+-------+
| AAA | BBB |  CCC  |
+-----+-----+-------+
| a   | b   | c     |
|   1 |   2 |     3 |
+-----+-----+-------+
|  A  |  B  | CCCC  |
+-----+-----+-------+
`
	table, buf := NewBuffered(
		WithHeader(header),
		WithFooter(footer),
		WithRows(data),
		WithColMinWidth(2, 5),
	)

	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestNumberAlign(t *testing.T) {
	data := [][]string{
		{"AAAAAAAAAAAAA", "BBBBBBBBBBBBB", "CCCCCCCCCCCCCC"},
		{"A", "B", "C"},
		{"123456789", "2", "3"},
		{"1", "2", "123,456,789"},
		{"1", "123,456.789", "3"},
		{"-123,456", "-2", "-3"},
	}

	const want = `+---------------+---------------+----------------+
| AAAAAAAAAAAAA | BBBBBBBBBBBBB | CCCCCCCCCCCCCC |
| A             | B             | C              |
|     123456789 |             2 |              3 |
|             1 |             2 |    123,456,789 |
|             1 |   123,456.789 |              3 |
|      -123,456 |            -2 |             -3 |
+---------------+---------------+----------------+
`
	table, buf := NewBuffered(
		WithRows(data),
	)
	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestCustomAlign(t *testing.T) {
	var (
		header = []string{"AAA", "BBB", "CCC"}
		data   = [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}
		footer = []string{"a", "b", "cccc"}
	)

	const want = `+-----+-----+-------+
| AAA | BBB |  CCC  |
+-----+-----+-------+
| a   |  b  |     c |
| 1   |  2  |     3 |
+-----+-----+-------+
|  A  |  B  | CCCC  |
+-----+-----+-------+
`

	table, buf := NewBuffered(
		WithHeader(header),
		WithFooter(footer),
		WithRows(data),
		WithColMinWidth(2, 5),
		WithColAlignment(map[int]HAlignment{
			0: AlignLeft,
			1: AlignCenter,
			2: AlignRight,
		}),
	)

	table.Render()

	checkEqual(t, buf.String(), want)
}

func TestKubeFormat(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "jan_hosting", "2233", "$10.98"},
		{"1/1/2014", "feb_hosting", "2233", "$54.95"},
		{"1/4/2014", "feb_extra_bandwidth", "2233", "$51.00"},
		{"1/4/2014", "mar_hosting", "2233", "$30.00"},
	}

	table, buf := NewBuffered(
		WithHeader([]string{"Date", "Description", "CV2", "Amount"}),
		WithWrap(false),
		WithTitledHeader(true),
		WithHeaderAlignment(AlignLeft),
		WithCellAlignment(AlignLeft),
		WithCenterSeparator(""),
		WithColumnSeparator(""),
		WithRowSeparator(""),
		WithHeaderLine(false),
		WithAllBorders(false),
		WithPadding("\t"), // pad with tabs
		WithNoWhiteSpace(true),
		WithRows(data),
	)

	table.Render()

	want := `DATE    	DESCRIPTION        	CV2 	AMOUNT
1/1/2014	jan_hosting        	2233	$10.98
1/1/2014	feb_hosting        	2233	$54.95
1/4/2014	feb_extra_bandwidth	2233	$51.00
1/4/2014	mar_hosting        	2233	$30.00
`

	checkEqual(t, buf.String(), want, "kube format rendering failed")
}
