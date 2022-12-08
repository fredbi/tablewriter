package titlers

import "testing"

func TestDefaultTitler(t *testing.T) {
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
		got := DefaultTitler(tt.text)
		if got != tt.want {
			t.Errorf("want %q, bot got %q", tt.want, got)
		}
	}
}
