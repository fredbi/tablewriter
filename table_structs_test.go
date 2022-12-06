package tablewriter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStringerType struct{}

func (t testStringerType) String() string { return "testStringerType" }

func TestStructs(t *testing.T) {
	type testType struct {
		A string
		B int
		C testStringerType
		D bool `tablewriter:"DD"`
	}
	type testType2 struct {
		A *string
		B *int
		C *testStringerType
		D *bool `tablewriter:"DD"`
	}
	type testType3 struct {
		A **string
		B **int
		C **testStringerType
		D **bool `tablewriter:"DD"`
	}
	a := "a"
	b := 1
	c := testStringerType{}
	d := true

	ap := &a
	bp := &b
	cp := &c
	dp := &d

	tests := []struct {
		name    string
		values  interface{}
		wantErr bool
		want    string
	}{
		{
			name: "slice of struct",
			values: []testType{
				{A: "AAA", B: 11, D: true},
				{A: "BBB", B: 22},
			},
			want: `
+-----+----+------------------+-------+
|  A  | B  |        C         |  DD   |
+-----+----+------------------+-------+
| AAA | 11 | testStringerType | true  |
| BBB | 22 | testStringerType | false |
+-----+----+------------------+-------+
`,
		},
		{
			name: "slice of struct pointer",
			values: []*testType{
				{A: "AAA", B: 11, D: true},
				{A: "BBB", B: 22},
			},
			want: `
+-----+----+------------------+-------+
|  A  | B  |        C         |  DD   |
+-----+----+------------------+-------+
| AAA | 11 | testStringerType | true  |
| BBB | 22 | testStringerType | false |
+-----+----+------------------+-------+
`,
		},
		{
			name: "pointer field",
			values: []*testType2{
				{A: &a, B: &b, C: &c, D: &d},
			},
			want: `
+---+---+------------------+------+
| A | B |        C         |  DD  |
+---+---+------------------+------+
| a | 1 | testStringerType | true |
+---+---+------------------+------+
`,
		},
		{
			name: "nil pointer field",
			values: []*testType2{
				{A: nil, B: nil, C: nil, D: nil},
			},
			want: `
+-----+-----+-----+-----+
|  A  |  B  |  C  | DD  |
+-----+-----+-----+-----+
| nil | nil | nil | nil |
+-----+-----+-----+-----+
`,
		},
		{
			name: "typed nil pointer field",
			values: []*testType2{
				{A: (*string)(nil), B: (*int)(nil), C: (*testStringerType)(nil), D: (*bool)(nil)},
			},
			want: `
+-----+-----+-----+-----+
|  A  |  B  |  C  | DD  |
+-----+-----+-----+-----+
| nil | nil | nil | nil |
+-----+-----+-----+-----+
`,
		},
		{
			name: "pointer of pointer field",
			values: []*testType3{
				{A: &ap, B: &bp, C: &cp, D: &dp},
			},
			want: `
+---+---+------------------+------+
| A | B |        C         |  DD  |
+---+---+------------------+------+
| a | 1 | testStringerType | true |
+---+---+------------------+------+
`,
		},
		{
			name:    "invalid input",
			values:  interface{}(1),
			wantErr: true,
		},
		{
			name:    "invalid input",
			values:  testType{},
			wantErr: true,
		},
		{
			name:    "invalid input",
			values:  &testType{},
			wantErr: true,
		},
		{
			name:    "nil value",
			values:  nil,
			wantErr: true,
		},
		{
			name:    "the first element is nil",
			values:  []*testType{nil, nil},
			wantErr: true,
		},
		{
			name:    "empty slice",
			values:  []testType{},
			wantErr: true,
		},
		{
			name: "mixed slice", // TODO: Should we support this case?
			values: []interface{}{
				testType{A: "a", B: 2, C: c, D: false},
				testType2{A: &a, B: &b, C: &c, D: &d},
				testType3{A: &ap, B: &bp, C: &cp, D: &dp},
			},
			wantErr: true,
		},
		{
			name: "skip nil element",
			values: []*testType{
				{A: "a", B: 1, D: true},
				nil,
				nil,
				{A: "A", B: 3, D: false},
			},
			want: `
+---+---+------------------+-------+
| A | B |        C         |  DD   |
+---+---+------------------+-------+
| a | 1 | testStringerType | true  |
| A | 3 | testStringerType | false |
+---+---+------------------+-------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, buf := NewBuffered()

			err := table.SetStructs(tt.values)
			if tt.wantErr {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)

			table.Render()
			checkEqual(t, buf.String(), strings.TrimPrefix(tt.want, "\n"))
		})
	}
}
