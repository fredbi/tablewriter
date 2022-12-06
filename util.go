// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"strings"
)

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

// String value based on condition
func conditionString(cond bool, valid, inValid string) string {
	if cond {
		return valid
	}

	return inValid
}

func isNumOrSpace(r rune) bool {
	return ('0' <= r && r <= '9') || r == ' '
}

// Format Table Header
// Replace _ , . and spaces
func title(name string) string {
	origLen := len(name)
	rs := []rune(name)
	for i, r := range rs {
		switch r {
		case '_':
			rs[i] = ' '
		case '.':
			// ignore floating number 0.0
			if (i != 0 && !isNumOrSpace(rs[i-1])) || (i != len(rs)-1 && !isNumOrSpace(rs[i+1])) {
				rs[i] = ' '
			}
		}
	}
	name = string(rs)
	name = strings.TrimSpace(name)
	if len(name) == 0 && origLen > 0 {
		// Keep at least one character. This is important to preserve
		// empty lines in multi-line headers/footers.
		name = " "
	}

	return strings.ToUpper(name)
}

// enforce all cells in the row to have the same number of lines
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

// getLines decomposes a multiline string into a slice of strings.
func getLines(s string) []string {
	return strings.Split(s, NEWLINE) // TODO: what if CRLF
}
