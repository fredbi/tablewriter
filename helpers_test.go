package tablewriter

import (
	"io"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func checkEqual(t *testing.T, got, want interface{}, msgs ...interface{}) {
	t.Helper()

	if !assert.EqualValues(t, want, got, msgs...) {
		wantStr, wantString := want.(string)
		gotStr, gotString := got.(string)
		if wantString && gotString {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(wantStr, gotStr, true)
			t.Logf("Diff:\n%s", dmp.DiffPrettyText(diffs))
		}
	}
}

func newCustomizedTable(out io.Writer) *Table {
	return New(
		WithWriter(out),
		WithCenterSeparator(""),
		WithColumnSeparator(""),
		WithRowSeparator(""),
		WithAllBorders(false),
		WithCellAlignment(AlignLeft),
		WithHeader([]string{}),
	)
}
