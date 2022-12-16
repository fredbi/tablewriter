package tablewrappers

import (
	// "log"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	space          = " "
	defaultPenalty = 1e5
)

// wrapWords is the low-level line-breaking algorithm, useful if you need more
// control over the details of the text wrapping process.
//
// wrapWords splits a list of words into lines with minimal "raggedness",
// treating each rune as one unit, accounting for spc units between adjacent
// words on each line, and attempting to limit lines to lim units. Raggedness
// is the total error over all lines, where error is the square of the
// difference of the length of the line and lim. Too-long lines (which only
// happen when a single word is longer than lim units) have pen penalty units
// added to the error.
func wrapWords(words []string, spc, limit, penalty int) [][]string {
	lengths, maxWordLength := buildLengthsMatrix(words, spc)
	n := len(lengths)
	nbrk := make([]int, n)
	costVector := initCosts(n)

	// guard: if any word is larger than the limit:
	// there is no point in trying to abide by this limit: adjust the limit
	// to result in a best effort.
	limit = max(limit, maxWordLength)

	for i := n - 1; i >= 0; i-- {
		if lengths[i][n-1] <= limit {
			costVector[i] = 0
			nbrk[i] = n

			continue
		}

		for j := i + 1; j < n; j++ {
			d := limit - lengths[i][j-1]
			c := d*d + costVector[j]

			if d < 0 {
				c += penalty // too-long lines get a worse penalty
			}

			if c < costVector[i] {
				costVector[i] = c
				nbrk[i] = j // add break point
			}
		}

		// safeguard: no break point was found
		if nbrk[i] == 0 {
			nbrk[i] = n
		}
	}

	var lines [][]string
	i := 0
	for i < n { // walk break points
		paragraph := stripEmpty(words[i:nbrk[i]])
		if len(paragraph) > 0 {
			lines = append(lines, paragraph)
		}
		i = nbrk[i]
	}

	return lines
}

func stripEmpty(words []string) []string {
	out := make([]string, 0, len(words))
	for _, word := range words {
		if len(word) == 0 {
			continue // that's an empty string, not a string that has a zero-width
		}
		out = append(out, word)
	}

	return out
}

// buildLengthsMatrix builds an upper triangular matrix of increasing lengths when assembling words.
func buildLengthsMatrix(words []string, spc int) ([][]int, int) {
	n := len(words)
	length := make([][]int, n)
	var maxWordLength int

	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = runewidth.StringWidth(words[i]) // TODO: use displayWidth(str string) to strip control sequences
		maxWordLength = max(maxWordLength, length[i][i])
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if length[i][j-1] > 0 {
				length[i][j] = length[i][j-1] + spc + length[j][j]
			} else {
				length[i][j] = length[j][j]
			}
		}
	}

	return length, maxWordLength
}

func initCosts(n int) []int {
	costVector := make([]int, n)

	for i := range costVector {
		costVector[i] = math.MaxInt32
	}

	return costVector
}

func wrapMultiline(words []string, limit int) []string {
	var lines []string

	for _, words := range wrapWords(words, 1, limit, defaultPenalty) {
		lines = append(lines, strings.Join(words, space))
	}

	if len(lines) == 0 {
		// always ensure at least one line
		lines = []string{""}
	}

	return lines
}
