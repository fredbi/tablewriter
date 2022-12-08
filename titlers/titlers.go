package titlers

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type (
	// DefaultTitler formats headers and footers.
	//
	// It replaces _ , . and spaces.
	DefaultTitler struct {
	}

	// CaseTitler returns a titler based on golang.org/x/text/cases.Caser.
	CaseTitler struct {
		cases.Caser
	}
)

func NewDefault() *DefaultTitler {
	return new(DefaultTitler)
}

func NewCaseTitler(tag language.Tag, opts ...cases.Option) *CaseTitler {
	return &CaseTitler{
		Caser: cases.Title(tag, opts...),
	}
}

func (t *DefaultTitler) Title(name string) string {
	origLen := len(name)
	rs := []rune(name)

	for i, r := range rs {
		switch r {
		case '_':
			rs[i] = ' '
		case '.':
			// ignore floating number 0.0
			if (i > 0 && !isNumOrSpace(rs[i-1])) || (i < len(rs)-1 && !isNumOrSpace(rs[i+1])) {
				rs[i] = ' '
			}
		}
	}
	name = strings.TrimSpace(string(rs))

	if len(name) == 0 && origLen > 0 {
		// Keep at least one character. This is important to preserve
		// empty lines in multi-line headers/footers.
		name = " "
	}

	return strings.ToUpper(name)
}

func isNumOrSpace(r rune) bool {
	return ('0' <= r && r <= '9') || r == ' '
}

func (t *CaseTitler) Title(name string) string {
	return t.String(name)
}
