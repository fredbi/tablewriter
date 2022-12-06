// Package tablewriter exposes a utility to render tabular data as text.
package tablewriter

type Wrapper interface {
	WrapString(string, int) []string
}
