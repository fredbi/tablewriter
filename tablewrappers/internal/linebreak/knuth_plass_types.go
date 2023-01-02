package linebreak

type (
	sums struct {
		width   int
		stretch int
		shrink  int
	}

	breakPoint struct {
		position     int
		demerits     float64
		ratio        float64
		line         int
		fitnessClass breakClass
		totals       sums
		previous     *breakPoint
	}

	nodeT struct {
		nodeType nodeType
		penalty  float64
		flagged  float64
		value    string

		sums
	}

	nodeType uint8

	demeritsT struct {
		line    float64
		flagged float64
		fitness float64
	}

	candidateT struct {
		active   *breakPoint
		demerits float64
		ratio    float64
	}

	lineT struct {
		ratio    float64
		nodes    []nodeT
		position int
	}

	breakClass = int

	err string
)

const (
	// ErrCannotBeSet indicates that the display constraints cannot be met with the given tolerance parameter.
	ErrCannotBeSet err = "paragraph cannot be set with the given tolerance"
)

const (
	nodeTypePenalty nodeType = iota
	nodeTypeGlue
	nodeTypeBox
)

const (
	breakClassZero breakClass = iota
	breakClassOne
	breakClassTwo
	breakClassThree
)

func newBreakPoint(position int, demerits float64, ratio float64, line int, fitnessClass int, totals sums, previous *breakPoint) *breakPoint {
	return &breakPoint{
		position: position,
		demerits: demerits,
		ratio:    ratio,
		line:     line,
		totals:   totals,
		previous: previous,
	}
}

func newGlue(width, stretch, shrink int) nodeT {
	return nodeT{
		nodeType: nodeTypeGlue,
		sums: sums{
			width:   width,
			stretch: stretch,
			shrink:  shrink,
		},
	}
}

func newBox(width int, value string) nodeT {
	return nodeT{
		nodeType: nodeTypeBox,
		value:    value,
		sums: sums{
			width: width,
		},
	}
}

func newPenalty(width int, penalty float64, flagged float64) nodeT {
	return nodeT{
		nodeType: nodeTypePenalty,
		sums: sums{
			width: width,
		},
		penalty: penalty,
		flagged: flagged,
	}
}

func defaultCandidates() []candidateT {
	return []candidateT{
		{demerits: maxDemerit},
		{demerits: maxDemerit},
		{demerits: maxDemerit},
		{demerits: maxDemerit},
	}
}

func (t nodeType) String() string {
	switch t {
	case nodeTypePenalty:
		return "penalty"
	case nodeTypeGlue:
		return "glue"
	case nodeTypeBox:
		return "box"
	default:
		return ""
	}
}

func (e err) Error() string {
	return string(e)
}

func (s *sums) Add(t sums) {
	s.width += t.width
	s.stretch += t.stretch
	s.shrink += t.shrink
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}
