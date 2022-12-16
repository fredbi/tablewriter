// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

type (
	transformer    func(string) string
	colPadder      func(string, int, int) string
	colTransformer func(int) transformer
)

func identity(in string) string { return in }

func stringIf(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}

	return ifFalse
}

// enforce all cells in a row to have the same number of lines
func normalizeRowHeight(columns [][]string, height int) [][]string {
	for i, rowLines := range columns {
		currentHeight := len(rowLines)
		padHeight := height - currentHeight

		for n := 0; n < padHeight; n++ {
			columns[i] = append(columns[i], "")
		}
	}

	return columns
}
