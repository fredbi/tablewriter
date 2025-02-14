package tablewriter_test

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/fredbi/tablewriter"
	"github.com/logrusorgru/aurora/v4"
)

func sampleData() [][]string {
	return [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
		{"D", "The Gopher", "800"},
	}
}

func ExampleTable() {
	data := sampleData()
	table := tablewriter.New(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithRows(data),
	)

	table.Render()

	// Output:
	// +------+-----------------------+--------+
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
	// *=============================*===============================*===============================*==========*
	// |            NAME             |             SIGN              |            RATING             |          |
	// *=============================*===============================*===============================*==========*
	// | Learn East has computers    | Some Data                     | Another Data                  |          |
	// | with adapted keyboards with |                               |                               |          |
	// | enlarged print etc          |                               |                               |          |
	// | Instead of lining up the    | the way across, he splits the | Like most ergonomic keyboards | See Data |
	// | letters all                 | keyboard in two               |                               |          |
	// *=============================*===============================*===============================*==========*
}

func ExampleWithNoWhiteSpace() {
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
	// *=============================*===============================*===============================*==========*
	// |            NAME             |             SIGN              |            RATING             |          |
	// *=============================*===============================*===============================*==========*
	// | Learn East has computers    | Some Data                     | Another Data                  |          |
	// | with adapted keyboards with |                               |                               |          |
	// | enlarged print etc          |                               |                               |          |
	// | Instead of lining up the    | the way across, he splits the | Like most ergonomic keyboards | See Data |
	// | letters all                 | keyboard in two               |                               |          |
	// *=============================*===============================*===============================*==========*
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
	const multiline = `A multiline
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
		for _, titled := range []bool{false, true} {
			if mode == testRow && titled {
				// Nothing special to test, skip
				continue
			}

			for _, wrapped := range []bool{false, true} {
				fmt.Println("mode:", mode, "titled:", titled, "wrapped:", wrapped)

				options := []tablewriter.Option{
					tablewriter.WithWriter(os.Stdout),
					tablewriter.WithTitledHeader(titled),
					tablewriter.WithWrap(wrapped),
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

	// Output:
	// mode: 0 titled: false wrapped: false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | A multiline                               |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	//
	// mode: 0 titled: false wrapped: true
	// +-----+------------------------------+
	// | woo |             waa              |
	// +-----+------------------------------+
	// | woo | A multiline string with some |
	// |     | lines being really long.     |
	// +-----+------------------------------+
	// | woo |             waa              |
	// +-----+------------------------------+
	//
	// mode: 1 titled: false wrapped: false
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	//
	// mode: 1 titled: false wrapped: true
	// +-----+------------------------------+
	// | woo | A multiline string with some |
	// |     |   lines being really long.   |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// | woo |             waa              |
	// +-----+------------------------------+
	//
	// mode: 1 titled: true wrapped: false
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	//
	// mode: 1 titled: true wrapped: true
	// +-----+------------------------------+
	// | WOO | A MULTILINE STRING WITH SOME |
	// |     |   LINES BEING REALLY LONG    |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// | WOO |             WAA              |
	// +-----+------------------------------+
	//
	// mode: 2 titled: false wrapped: false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | woo |                A multiline                |
	// |     | string with some lines being really long. |
	// +-----+-------------------------------------------+
	//
	// mode: 2 titled: false wrapped: true
	// +-----+------------------------------+
	// | woo |             waa              |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// | woo | A multiline string with some |
	// |     |   lines being really long.   |
	// +-----+------------------------------+
	//
	// mode: 2 titled: true wrapped: false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// | WOO |                A MULTILINE                |
	// |     | STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	//
	// mode: 2 titled: true wrapped: true
	// +-----+------------------------------+
	// | WOO |             WAA              |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// | WOO | A MULTILINE STRING WITH SOME |
	// |     |   LINES BEING REALLY LONG    |
	// +-----+------------------------------+
	//
	// mode: 3 titled: false wrapped: false
	// +-----+-------------------------------------------+
	// | woo |                    waa                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A multiline                |
	// |       string with some lines being really long. |
	// +-----+-------------------------------------------+
	//
	// mode: 3 titled: false wrapped: true
	// +-----+------------------------------+
	// | woo |             waa              |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// |       A multiline string with some |
	// |         lines being really long.   |
	// +-----+------------------------------+
	//
	// mode: 3 titled: true wrapped: false
	// +-----+-------------------------------------------+
	// | WOO |                    WAA                    |
	// +-----+-------------------------------------------+
	// | woo | waa                                       |
	// +-----+-------------------------------------------+
	// |                      A MULTILINE                |
	// |       STRING WITH SOME LINES BEING REALLY LONG  |
	// +-----+-------------------------------------------+
	//
	// mode: 3 titled: true wrapped: true
	// +-----+------------------------------+
	// | WOO |             WAA              |
	// +-----+------------------------------+
	// | woo | waa                          |
	// +-----+------------------------------+
	// |       A MULTILINE STRING WITH SOME |
	// |         LINES BEING REALLY LONG    |
	// +-----+------------------------------+
}

func ExampleFormatter() {
	data := sampleData()
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

	// Output:
	// +------+-----------------------+--------+
	// | [31mNAME[0m | [34m        SIGN         [0m | [1mRATING[0m |
	// +------+-----------------------+--------+
	// | A    | The Good              |    500 |
	// | B    | The Very very Bad Man |    288 |
	// | C    | The Ugly              |    120 |
	// | D    | The Gopher            |    800 |
	// +------+-----------------------+--------+
}

func ExampleWithRowSeparator() {
	data := sampleData()
	red := func(in string) string {
		return aurora.Sprintf(aurora.Red(in))
	}

	// prints a colorized grid
	table := tablewriter.New(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithRows(data),
		tablewriter.WithRowSeparator(red(tablewriter.ROW)),
		tablewriter.WithColumnSeparator(red(tablewriter.COLUMN)),
		tablewriter.WithCenterSeparator(red(tablewriter.CENTER)),
	)

	table.Render()

	// Output:
	// [31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m
	// [31m|[0m NAME [31m|[0m         SIGN          [31m|[0m RATING [31m|[0m
	// [31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m
	// [31m|[0m A    [31m|[0m The Good              [31m|[0m    500 [31m|[0m
	// [31m|[0m B    [31m|[0m The Very very Bad Man [31m|[0m    288 [31m|[0m
	// [31m|[0m C    [31m|[0m The Ugly              [31m|[0m    120 [31m|[0m
	// [31m|[0m D    [31m|[0m The Gopher            [31m|[0m    800 [31m|[0m
	// [31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m-[0m[31m+[0m
}

func ExampleWithMaxTableWidth() {
	data := sampleData()

	// adapt the width of the table to the maximum display size
	table := tablewriter.New(
		tablewriter.WithHeader([]string{"Name", "Sign", "Rating"}),
		tablewriter.WithRows(data),
		tablewriter.WithMaxTableWidth(30),
	)
	table.Render()

	// Output:
	// 	+------+------------+--------+
	// | NAME |    SIGN    | RATING |
	// +------+------------+--------+
	// | A    | The Good   |    500 |
	// | B    | The Very   |    288 |
	// |      | very Bad   |        |
	// |      | Man        |        |
	// | C    | The Ugly   |    120 |
	// | D    | The Gopher |    800 |
	// +------+------------+--------+
}
