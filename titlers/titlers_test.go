package titlers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestDefaultTitler(t *testing.T) {
	titler := NewDefault()

	ts := []struct {
		text string
		want string
	}{
		{"", ""},
		{"foo", "FOO"},
		{"Foo", "FOO"},
		{"foO", "FOO"},
		{".foo", "FOO"},
		{"foo.", "FOO"},
		{".foo.", "FOO"},
		{".foo.bar.", "FOO BAR"},
		{"_foo", "FOO"},
		{"foo_", "FOO"},
		{"_foo_", "FOO"},
		{"_foo_bar_", "FOO BAR"},
		{" foo", "FOO"},
		{"foo ", "FOO"},
		{" foo ", "FOO"},
		{" foo bar ", "FOO BAR"},
		{"0.1", "0.1"},
		{"FOO 0.1", "FOO 0.1"},
		{".1 0.1", ".1 0.1"},
		{"1. 0.1", "1. 0.1"},
		{"1. 0.", "1. 0."},
		{".1. 0.", ".1. 0."},
		{".$ . $.", "$ . $"},
		{".$. $.", "$  $"},
	}
	for _, tt := range ts {
		require.Equal(t, tt.want, titler.Title(tt.text))
	}
}

func TestCaseTitler(t *testing.T) {
	titler := NewCaseTitler(language.AmericanEnglish)

	require.Equal(t, "This Is Titled", titler.Title("this is TITLED"))
}
