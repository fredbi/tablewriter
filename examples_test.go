package tablewriter_test

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/logrusorgru/aurora/v4"
	"github.com/olekukonko/tablewriter"
)

func ExampleTable() {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table := tablewriter.New(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithRows(data),
	)

	table.Render()

	// Output: +------+-----------------------+--------+
	// | NAME |         SIGN          | RATING |
	// +------+-----------------------+--------+
	// | A    | The Good              |    500 |
	// | B    | The Very very Bad Man |    288 |
	// | C    | The Ugly              |    120 |
	// | D    | The Gopher            |    800 |
	// +------+-----------------------+--------+
}

func ExampleOption() {
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "Some Data    ", "Another Data "},
		{"Instead of lining up the letters all ", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards", "See Data"},
	}

	table := tablewriter.New(
		tablewriter.WithWriter(os.Stdout), // default is os.Stdout
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithCenterSeparator("*"), // default is '+'
		tablewriter.WithRowSeparator("="),    // default is '-'
	)

	for _, v := range data {
		table.Append(v) // an alternative to WithRows(data)
	}

	table.Render()

	// Output:
	// *================================*================================*===============================*==========*
	// |              NAME              |              SIGN              |            RATING             |          |
	// *================================*================================*===============================*==========*
	// | Learn East has computers       | Some Data                      | Another Data                  |          |
	// | with adapted keyboards with    |                                |                               |          |
	// | enlarged print etc             |                                |                               |          |
	// | Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards | See Data |
	// | letters all                    | keyboard in two                |                               |          |
	// *================================*================================*===============================*==========*
}

func ExampleWithNoWhiteSpace() { // TODO: NoWhiteSpace option is still full of bugs
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "Some Data    ", "Another Data "},
		{"Instead of lining up the letters all ", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards", "See Data"},
	}

	options := []tablewriter.Option{
		tablewriter.WithWriter(os.Stdout), // default is os.Stdout
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithCenterSeparator("*"), // default is '+'
		tablewriter.WithRowSeparator("="),    // default is '-'
		tablewriter.WithRows(data),
	}

	// packed without extraneous blank space
	table := tablewriter.New(append(options,
		tablewriter.WithNoWhiteSpace(true),
	)...)
	table.Render()
	fmt.Println()

	// default layout
	table = tablewriter.New(append(options,
		tablewriter.WithNoWhiteSpace(false),
	)...)
	table.Render()
}

// disabled for now

// Output:
// *================================*================================*===============================*==========*
//             NAME                           SIGN                         RATING
// *================================*================================*===============================*==========*
// Learn East has computers       Some Data                      Another Data
// with adapted keyboards with
// enlarged print etc
// Instead of lining up the       the way across, he splits the  Like most ergonomic keyboards See Data
// letters all                    keyboard in two
// *================================*================================*===============================*==========*
//
// *================================*================================*===============================*==========*
// |              NAME              |              SIGN              |            RATING             |          |
// *================================*================================*===============================*==========*
// | Learn East has computers       | Some Data                      | Another Data                  |          |
// | with adapted keyboards with    |                                |                               |          |
// | enlarged print etc             |                                |                               |          |
// | Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards | See Data |
// | letters all                    | keyboard in two                |                               |          |
// *================================*================================*===============================*==========*

func ExampleNewBuffered() {
	data := [][]string{
		{"Learn East has computers with adapted keyboards with enlarged print etc", "Some Data    ", "Another Data "},
		{"Instead of lining up the letters all ", "the way across, he splits the keyboard in two", "Like most ergonomic keyboards", "See Data"},
	}

	table, buf := tablewriter.NewBuffered(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithCenterSeparator("*"), // default is '+'
		tablewriter.WithRowSeparator("="),    // default is '-'
	)

	for _, v := range data {
		table.Append(v)
	}

	table.Render() // writes to buffer

	fmt.Println(buf)

	// Output:
	// *================================*================================*===============================*==========*
	// |              NAME              |              SIGN              |            RATING             |          |
	// *================================*================================*===============================*==========*
	// | Learn East has computers       | Some Data                      | Another Data                  |          |
	// | with adapted keyboards with    |                                |                               |          |
	// | enlarged print etc             |                                |                               |          |
	// | Instead of lining up the       | the way across, he splits the  | Like most ergonomic keyboards | See Data |
	// | letters all                    | keyboard in two                |                               |          |
	// *================================*================================*===============================*==========*
}

func ExampleNewCSV() {
	file, err := os.Open("testdata/test.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)

	table, err := tablewriter.NewCSV(reader, true,
		tablewriter.WithCenterSeparator("*"),
		tablewriter.WithRowSeparator("="),
	)
	if err != nil {
		log.Fatal(err)
	}

	table.Render()

	// Output:
	// *============*===========*=========*
	// | FIRST NAME | LAST NAME |   SSN   |
	// *============*===========*=========*
	// | John       | Barry     |  123456 |
	// | Kathy      | Smith     |  687987 |
	// | Bob        | McCornick | 3979870 |
	// *============*===========*=========*
}

func ExampleWithWrap() {
	var multiline = `A multiline
string with some lines being really long.`

	type testMode uint8
	const (
		// test mode
		testRow testMode = iota
		testHeader
		testFooter
		testFooter2
	)

	for mode := testRow; mode <= testFooter2; mode++ {
		for _, autoFmt := range []bool{false, true} {
			if mode == testRow && autoFmt {
				// Nothing special to test, skip
				continue
			}

			for _, autoWrap := range []bool{false, true} {
				for _, reflow := range []bool{false, true} {
					if !autoWrap && reflow {
						// Invalid configuration, skip
						continue
					}

					fmt.Println("mode", mode, "autoFmt", autoFmt, "autoWrap", autoWrap, "reflow", reflow)

					options := []tablewriter.Option{
						tablewriter.WithWriter(os.Stdout),
						tablewriter.WithTitledHeader(autoFmt),
						tablewriter.WithWrap(autoWrap),
						tablewriter.WithWrapReflow(reflow), // TODO: wrap option
					}

					switch mode {
					case testHeader:
						options = append(options, tablewriter.WithHeader([]string{"woo", multiline}))
						options = append(options, tablewriter.WithFooter([]string{"woo", "waa"}))
						options = append(options, tablewriter.WithRows([][]string{{"woo", "waa"}}))
					case testRow:
						options = append(options, tablewriter.WithHeader([]string{"woo", "waa"}))
						options = append(options, tablewriter.WithFooter([]string{"woo", "waa"}))
						options = append(options, tablewriter.WithRows([][]string{{"woo", multiline}}))
					case testFooter:
						options = append(options, tablewriter.WithHeader([]string{"woo", "waa"}))
						options = append(options, tablewriter.WithFooter([]string{"woo", multiline}))
						options = append(options, tablewriter.WithRows([][]string{{"woo", "waa"}}))
					case testFooter2:
						options = append(options, tablewriter.WithHeader([]string{"woo", "waa"}))
						options = append(options, tablewriter.WithFooter([]string{"", multiline}))
						options = append(options, tablewriter.WithRows([][]string{{"woo", "waa"}}))
					}

					t := tablewriter.New(options...)
					t.Render()
					fmt.Println()
				}
			}
		}
	}

	// Output:
	// mode 0 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | A multiline                               |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	//
	// mode 0 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | A multiline                    |
	// |     |                                |
	// |     | string with some lines being   |
	// |     | really long.                   |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	//
	// mode 0 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | A multiline string with some   |
	// |     | lines being really long.       |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	//
	// mode 1 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	//
	// mode 1 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |          A multiline           |
	// |     |                                |
	// |     |  string with some lines being  |
	// |     |          really long.          |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	//
	// mode 1 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |  A multiline string with some  |
	// |     |    lines being really long.    |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	//
	// mode 1 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	//
	// mode 1 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |          A MULTILINE           |
	// |     |                                |
	// |     |  STRING WITH SOME LINES BEING  |
	// |     |          REALLY LONG           |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	//
	// mode 1 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |  A MULTILINE STRING WITH SOME  |
	// |     |    LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	//
	// mode 2 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	//
	// mode 2 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |          A multiline           |
	// |     |                                |
	// |     |  string with some lines being  |
	// |     |          really long.          |
	// +-----+--------------------------------+
	//
	// mode 2 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | woo |  A multiline string with some  |
	// |     |    lines being really long.    |
	// +-----+--------------------------------+
	//
	// mode 2 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	//
	// mode 2 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |          A MULTILINE           |
	// |     |                                |
	// |     |  STRING WITH SOME LINES BEING  |
	// |     |          REALLY LONG           |
	// +-----+--------------------------------+
	//
	// mode 2 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// | WOO |  A MULTILINE STRING WITH SOME  |
	// |     |    LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
	//
	// mode 3 autoFmt false autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A multiline                |
	// |       string with some lines being really long. |
	// +-----+-------------------------------------------+
	//
	// mode 3 autoFmt false autoWrap true reflow false
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |                A multiline           |
	// |                                      |
	// |        string with some lines being  |
	// |                really long.          |
	// +-----+--------------------------------+
	//
	// mode 3 autoFmt false autoWrap true reflow true
	// +-----+--------------------------------+
	// | woo |              waa               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |        A multiline string with some  |
	// |          lines being really long.    |
	// +-----+--------------------------------+
	//
	// mode 3 autoFmt true autoWrap false reflow false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A MULTILINE                |
	// |       STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	//
	// mode 3 autoFmt true autoWrap true reflow false
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |                A MULTILINE           |
	// |                                      |
	// |        STRING WITH SOME LINES BEING  |
	// |                REALLY LONG           |
	// +-----+--------------------------------+
	//
	// mode 3 autoFmt true autoWrap true reflow true
	// +-----+--------------------------------+
	// | WOO |              WAA               |
	// +-----+--------------------------------+
	// | woo | waa                            |
	// +-----+--------------------------------+
	// |        A MULTILINE STRING WITH SOME  |
	// |          LINES BEING REALLY LONG     |
	// +-----+--------------------------------+
}

func ExampleFormatter() {
	data := [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}

	table := tablewriter.New(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithRows(data),
		tablewriter.WithHeaderFormatters(map[int]tablewriter.Formatter{
			0: aurora.Red,
			1: aurora.Blue,
			2: aurora.Bold,
		}),
	)

	table.Render()

	// Output: +------+-----------------------+--------+
	// | [31mNAME[0m | [34m        SIGN         [0m | [1mRATING[0m |
	// +------+-----------------------+--------+
	// | A    | The Good              |    500 |
	// | B    | The Very very Bad Man |    288 |
	// | C    | The Ugly              |    120 |
	// | D    | The Gopher            |    800 |
	// +------+-----------------------+--------+
}
