// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package wrap

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/require"
)

var text = "The quick brown fox jumps over the lazy dog."

func TestWrap(t *testing.T) {
	exp := []string{
		"The", "quick", "brown", "fox",
		"jumps", "over", "the", "lazy", "dog."}

	w := New()
	got := w.WrapString(text, 6)
	require.EqualValues(t, len(exp), len(got))
}

func TestWrapOneLine(t *testing.T) {
	exp := "The quick brown fox jumps over the lazy dog."
	w := New()
	words := w.WrapString(text, 500)
	require.EqualValues(t, exp, strings.Join(words, string(sp)))

}

func TestUnicode(t *testing.T) {
	input := "Česká řeřicha"
	w := New()
	var wordsUnicode []string
	if runewidth.IsEastAsian() {
		wordsUnicode = w.WrapString(input, 14)
	} else {
		wordsUnicode = w.WrapString(input, 13)
	}
	// input contains 13 (or 14 for CJK) runes, so it fits on one line.
	require.Len(t, wordsUnicode, 1)
}

func TestDisplayWidth(t *testing.T) {
	input := "Česká řeřicha"
	want := 13
	if runewidth.IsEastAsian() {
		want = 14
	}
	if n := displayWidth(input); n != want {
		t.Errorf("Wants: %d Got: %d", want, n)
	}
	input = "\033[43;30m" + input + "\033[00m"
	require.EqualValues(t, want, DisplayWidth(input))
}

func TestWrapString(t *testing.T) {
	want := []string{"ああああああああああああああああああああああああ", "あああああああ"}
	w := New()
	got := w.WrapString("ああああああああああああああああああああああああ あああああああ", 55)

	require.EqualValues(t, want, got)
}
