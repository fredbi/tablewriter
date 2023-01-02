package hyphenator

import (
	"golang.org/x/text/language"
)

type (
	Option func(*options)

	options struct {
		lang                 language.Tag
		minRunesBeforeHyphen int
	}
)

func WithLanguageTag(tag language.Tag) Option {
	return func(o *options) {
		o.lang = tag
	}
}

func WithLanguage(lang string) Option {
	tag, _ := language.MatchStrings(langMatcher, lang)

	return func(o *options) {
		o.lang = tag
	}
}

func defaultOptions(opts []Option) *options {
	o := &options{
		lang:                 language.AmericanEnglish,
		minRunesBeforeHyphen: 3,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}
