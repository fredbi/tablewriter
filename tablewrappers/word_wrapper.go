package tablewrappers

import (
	"math"
	"strings"
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
//
// Notice an alternative approach there: https://github.com/mitchellh/go-wordwrap.
//
// NOTE(fred): the absolute reference ever written on that topic may be found here:
// https://tug.org/TUGboat/tb21-3/tb68fine.pdf.
// https://fdocuments.net/document/breaking-paragraphs-into-lines-github-pages-donald-e-knuth-and-michael-f-plass.html?page=9
func wrapWords(words []string, spc, limit, penalty int) [][]string {
	lengths, maxWordLength, _ := buildLengthsMatrix(words, spc)
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

// stripEmpty prunes empty strings for a list of words.
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
func buildLengthsMatrix(words []string, spc int) ([][]int, int, int) {
	n := len(words)
	length := make([][]int, n)
	var (
		maxWordLength int
		minWordLength int
	)

	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = displayWidth(words[i])

		maxWordLength = max(maxWordLength, length[i][i])
		if minWordLength == 0 {
			minWordLength = length[i][i]
		} else {
			minWordLength = min(minWordLength, length[i][i])
		}
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

	return length, maxWordLength, minWordLength
}

func initCosts(n int) []int {
	costVector := make([]int, n)

	for i := range costVector {
		costVector[i] = math.MaxInt32
	}

	return costVector
}

func wrapMultiline(words []string, limit, spc int) []string {
	var lines []string
	/*
		startEsc string
		endEsc   string
	*/

	pad := strings.Repeat(space, spc)
	/*
				for _, word := range wordList {
		//word_wrappero_test.go:107: ["\x1b[43;30mABC" "XYZ" "123\x1b[00m"]
					stripped, start, end := stripANSI(word)

					switch {
					case start != "" && end == "":
						if endEsc != "" {
							stripped = endEsc + stripped
							endEsc = ""
						}
						if startEsc != "" {
							stripped = startEsc + stripped
						}
						startEsc = start
						wordList[i] = stripped
					case start == "" && end != "":
						if endEsc != "" {
							stripped = endEsc + stripped
						}
						if startEsc != "" {
							stripped = startEsc + stripped
							startEsc = ""
						}
						endEsc = end
						wordList[i] = stripped
					}

				}

				if startEsc != "" {
					wordList[0]
				}
	*/
	for _, wordList := range wrapWords(words, 1, limit, defaultPenalty) {
		lines = append(lines, strings.Join(wordList, pad))
	}

	if len(lines) == 0 {
		// always ensure at least one line
		lines = []string{""}
	}

	return lines
}
