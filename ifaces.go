package tablewriter

type (
	// Wrapper knows how to wrap a string into multiple lines,
	// under the constraint of a maximum display width.
	Wrapper interface {
		WrapString(input string, maxWidth int) []string
	}

	// Titler knows how to format an input string, suitable to display headings.
	Titler interface {
		Title(string) string
	}
)
