package linebreak

import (
	"container/list"
	"log"
	"math"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type (
	LineBreaker struct {
		nodes       []nodeT
		lineLengths []int
		sum         *sums
		activeNodes *list.List
		spaceWidth  int
		hyphenWidth int
		// spaceStretch int
		//spaceShrink  int

		*options
	}
)

const (
	infinity   = 10000.00
	maxDemerit = math.MaxFloat64

	hyphen         = "-"
	space          = " "
	flaggedPenalty = 1.0
	noWidth        = 0
	noShrink       = 0
)

func New(opts ...Option) *LineBreaker {
	l := &LineBreaker{
		options: defaultOptions(opts),
	}

	l.spaceWidth = l.scale(l.measurer(space))
	if l.wordBreak {
		l.hyphenWidth = l.scale(l.measurer(hyphen))
	}

	// NOTE: this is for center & justify (not implemented for now)
	// l.spaceStretch = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.stretch)))
	// l.spaceShrink = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.shrink)))

	return l
}

// LeftAlignUniform left-align a series of tokens that compose a paragraph,
// rendering multiple lines of uniform length maxLength.
func (l *LineBreaker) LeftAlignUniform(tokens []string, maxLength int) ([]string, error) {
	nodes := l.leftAlignedNodes(tokens)
	spew.Dump(nodes)
	lineLengths := l.buildUniformLengths(maxLength)

	breaks := l.breakPoints(nodes, lineLengths)
	if len(breaks) == 0 {
		return nil, ErrCannotBeSet
	}

	return l.render(nodes, breaks), nil
}

// buildUniformLengths fills 1 line of line length constraints,
// making this constraint uniform across all lines.
func (l *LineBreaker) buildUniformLengths(width int) []int {
	lengths := make([]int, 1)
	for i := range lengths {
		lengths[i] = l.scale(width)
	}

	return lengths
}

// scale a text measure to convert the measurement unit to a value compatible
// with the algorithm's parameters.
func (l *LineBreaker) scale(in int) int {
	return int(math.Round(float64(in) * l.scaleFactor))
}

func (l *LineBreaker) downScale(in int) int {
	return int(math.Round(float64(in) / l.scaleFactor))
}

func skipNodes(start int, nodes []nodeT) int {
	for j, node := range nodes[start:] {
		// skip nodes after a line break
		if node.nodeType == nodeTypeBox || (node.nodeType == nodeTypePenalty && node.penalty == -infinity) {
			start += j

			break
		}
	}

	return start
}

// render the nodes with the provided line breaks .
//
// TODO: state handling for ANSI escape sequences.
func (l *LineBreaker) render(nodes []nodeT, breaks []breakPoint) []string {
	lines := make([]lineT, 0, len(breaks))
	result := make([]string, 0, len(breaks))

	lineStart := 0
	for _, brk := range breaks[1:] {
		lineStart = skipNodes(lineStart, nodes)

		lines = append(lines, lineT{
			ratio:    brk.ratio,
			nodes:    nodes[lineStart : brk.position+1],
			position: brk.position,
		})

		lineStart = brk.position
	}

	for _, line := range lines {
		// x := 0
		//prevx := 0
		lineResult := new(strings.Builder)

		for index, node := range line.nodes {
			switch {
			case node.nodeType == nodeTypeBox:
				// render a box node
				//l.fillText(lineResult, node.value, prevx, x)
				lineResult.WriteString(node.value)
				//x += node.width
				//prevx = x

			case node.nodeType == nodeTypeGlue:
				// render a glue node
				//x += node.width
				// ignore ratio for now
				lineResult.WriteString(strings.Repeat(space, l.downScale(node.width))) // ici
				/*
					if line.ratio < 0 {
						lineResult.WriteString(strings.Repeat(space, max(l.downScale(node.width-node.shrink), 1))) // ici
						//		x += node.width + int(line.ratio*float64(node.shrink))
					} else {
						lineResult.WriteString(strings.Repeat(space, max(l.downScale(node.width+node.stretch), 1))) // ici
						//		x += node.width  + int(line.ratio*float64(node.stretch))
					}
				*/

			case node.nodeType == nodeTypePenalty:
				if node.penalty == l.hyphenPenalty {
					log.Printf("saw hypen penalty: %#v", node)
				}
				if node.penalty == l.hyphenPenalty && index == len(line.nodes)-1 {
					// render a hyphen penalty node
					if l.renderHyphens {
						//l.fillText(lineResult, hyphen, prevx, x)
						// prevx = 0
						lineResult.WriteString(hyphen)
					}
				}
			}
		}

		result = append(result, lineResult.String())
	}

	return result
}

// fillText appends a string at position x, padding with blank space from prevx up to x.
func (l *LineBreaker) fillText(b *strings.Builder, value string, prevx, x int) {
	if x == 0 {
		b.WriteString(value)

		return
	}

	pad := l.downScale(x - prevx)
	b.WriteString(strings.Repeat(space, pad))
	b.WriteString(value)
}

// boxNodes produces box nodes with possible word breaks (suited for left-aligned output).
// TODO(fredbi): process punctuation marks, word breaking on natural separators (e.g. /|-_)
func (l *LineBreaker) boxNodes(word string) []nodeT {
	if !l.wordBreak || len(word) <= l.minHyphenate {
		return []nodeT{newBox(l.scale(l.measurer(word)), word)}
	}

	hyphenated := l.hyphenator(word)
	if len(hyphenated) == 0 { // this rule is based on the # runes, not the width
		return []nodeT{newBox(l.scale(l.measurer(word)), word)}
	}

	boxNodes := make([]nodeT, 0, len(hyphenated))

	// word break points are associated with a penalty
	for _, part := range hyphenated[:len(hyphenated)-1] {
		boxNodes = append(boxNodes,
			newBox(l.scale(l.measurer(part)), part),
		)

		if l.renderHyphens {
			// when rendering hyphens, the penalty incurs some consumed width
			boxNodes = append(boxNodes,
				newPenalty(l.hyphenWidth, l.hyphenPenalty, flaggedPenalty),
			)
		} else {
			// when hyphens are not rendered (words are just broken), there is no width associated to the penalty
			boxNodes = append(boxNodes,
				newPenalty(noWidth, l.hyphenPenalty, flaggedPenalty),
			)
		}
	}

	lastPart := hyphenated[len(hyphenated)-1]
	boxNodes = append(boxNodes,
		newBox(l.scale(l.measurer(lastPart)), lastPart),
	)

	return boxNodes
}

// leftAlignedNodes prepares nodes for left-alignment.
func (l *LineBreaker) leftAlignedNodes(tokens []string) []nodeT {
	nodes := make([]nodeT, 0, 4*(len(tokens)-1)+3)

	// transform tokens into a list of nodes of type (box|glue|penalty)
	for _, word := range tokens[:len(tokens)-1] {
		nodes = append(nodes, l.boxNodes(word)...) // a word token, possibly broken in parts
		nodes = append(nodes, newGlue(noWidth, l.glueStretch, noShrink))
		nodes = append(nodes, newPenalty(noWidth, 0, 0))
		nodes = append(nodes, newGlue(l.spaceWidth, -l.glueStretch, noShrink))
	}

	// last token: complete the list of nodes with a final infinite glue and penalty.
	nodes = append(nodes, l.boxNodes(tokens[len(tokens)-1])...)
	nodes = append(nodes, newGlue(noWidth, infinity, noShrink))
	nodes = append(nodes, newPenalty(noWidth, -infinity, flaggedPenalty))

	return nodes
}

// TODO: center(), justify()?

// yields the ordered slice of breakpoints
func (l *LineBreaker) breakPoints(nodes []nodeT, lineLengths []int) []breakPoint {
	// reset state
	l.lineLengths = lineLengths
	l.nodes = nodes
	l.sum = new(sums)
	l.activeNodes = list.New()
	l.activeNodes.PushBack(newBreakPoint(0, 0, 0, 0, 0, sums{}, nil)) // first empty node starting a paragraph

	for index, node := range l.nodes {
		log.Printf("breakpoints list: %d", l.activeNodes.Len())
		switch {
		case node.nodeType == nodeTypeBox:
			// accumulate the total width of word
			log.Printf("seeing box node %d [%q]", index, node.value)
			l.sum.width += node.width

		case node.nodeType == nodeTypeGlue:
			if index > 0 && l.nodes[index-1].nodeType == nodeTypeBox {
				// explore a glue following a word
				log.Printf("explore glue node %d", index)
				l.exploreForNode(node, index)
			}

			l.sum.Add(node.sums)

		case node.nodeType == nodeTypePenalty:
			// explore a penalty
			log.Printf("explore penalty node %d", index)
			l.exploreForNode(node, index)
		}
	}

	if l.activeNodes.Len() > 0 {
		nodeWithMinDemerits := l.findBestBreak()

		return makeBreakPoints(nodeWithMinDemerits)
	}

	return []breakPoint{}
}

func (l *LineBreaker) findBestBreak() *breakPoint {
	nodeWithMinDemerits := &breakPoint{
		demerits: maxDemerit,
	}

	for element := l.activeNodes.Front(); element != nil; element = element.Next() {
		node := element.Value.(*breakPoint)

		if node.demerits < nodeWithMinDemerits.demerits {
			nodeWithMinDemerits = node
		}
	}

	return nodeWithMinDemerits
}

// attention: currentLine starts with 1
func (l *LineBreaker) costRatio(sum *sums, index int, active breakPoint, currentLine int) float64 {
	actualWidth := sum.width - active.totals.width
	idealWidth := getIdealWidth(currentLine, l.lineLengths)

	if l.nodes[index].nodeType == nodeTypePenalty {
		// for penalties with a width, e.g. extra hyphen
		actualWidth += l.nodes[index].width
	}

	switch {
	case actualWidth < idealWidth:
		// need to stretch
		stretch := sum.stretch - active.totals.stretch

		if stretch > 0 {
			return float64(idealWidth-actualWidth) / float64(stretch)
		}

		return infinity

	case actualWidth > idealWidth:
		// need to shrink
		shrink := sum.shrink - active.totals.shrink

		if shrink > 0 {
			return float64(idealWidth-actualWidth) / float64(shrink)
		}

		return -infinity

	default:
		// perfect match
		return 0.00
	}
}

func (l *LineBreaker) sumFromNode(index int) sums {
	sum := *l.sum
	for i, node := range l.nodes[index:] {
		if node.nodeType == nodeTypeGlue {
			sum.Add(node.sums)

			continue
		}

		if node.nodeType == nodeTypeBox || (node.nodeType == nodeTypePenalty && node.penalty == -infinity && i > 0) {
			// stop summing up when a new word or a forced line break is found
			break
		}
	}

	return sum
}

func (l *LineBreaker) exploreForNode(node nodeT, index int) {
	activeElement := l.activeNodes.Front()

	var (
		currentLine int // will range over lines starting from 1
		candidates  []candidateT
	)

	for activeElement != nil {
		candidates = defaultCandidates()

		// break points up to the current line
		for activeElement != nil {
			active := activeElement.Value.(*breakPoint)
			currentLine = active.line + 1
			ratio := l.costRatio(l.sum, index, *active, currentLine)
			log.Printf("node[%d][%v](%q) ratio=%v", index, l.nodes[index].nodeType, l.nodes[index].value, ratio)
			next := activeElement.Next()

			var removed bool
			if ratio < -1 || (node.nodeType == nodeTypePenalty && node.penalty == -infinity) {
				log.Printf("node[%d][%v](%q) remove", index, l.nodes[index].nodeType, l.nodes[index].value)
				// undesirable node or forced line break
				l.activeNodes.Remove(activeElement)
				removed = true
			}

			if ratio >= -1 && ratio <= l.tolerance {
				if removed {
					log.Printf("WARN: removed node but computing demerits")
				}
				demerits := l.demeritsForRatio(node.nodeType, ratio, node.penalty)

				if node.nodeType == nodeTypePenalty && l.nodes[active.position].nodeType == nodeTypePenalty {
					// penalize flagged penalty nodes
					demerits += l.demerits.flagged * node.flagged * l.nodes[active.position].flagged
				}

				currentClass := fitnessClassForRatio(ratio)

				// add demerits due to fitness class whenever the fitness of 2 adjacent lines differ too much
				if abs(currentClass-active.fitnessClass) > 1 {
					demerits += l.demerits.fitness
				}

				demerits += active.demerits

				if demerits < candidates[currentClass].demerits {
					// we have a better candidate for this class
					candidates[currentClass] = candidateT{
						active:   active,
						demerits: demerits,
						ratio:    ratio,
					}
					log.Printf("we have a better candidate: %#v", candidates[currentClass])
				}
			}

			activeElement = next

			if activeElement != nil && active.line >= currentLine {
				// stop iterating to add new candidate break points
				break
			}
		}

		sum := l.sumFromNode(index)
		for fitnessClass, candidate := range candidates {
			if candidate.demerits >= maxDemerit {
				// skip default candidate
				continue
			}

			newBreak := newBreakPoint(
				index,                               // break at node index
				candidate.demerits, candidate.ratio, // ratings for this break point
				candidate.active.line+1, fitnessClass,
				sum,              // totals from the node
				candidate.active, // link to the previous candidate breakpoint
			)

			if activeElement != nil {
				log.Printf("node[%d][%v](%q) insert break before", index, l.nodes[index].nodeType, l.nodes[index].value)
				_ = l.activeNodes.InsertBefore(newBreak, activeElement)

				continue
			}

			log.Printf("node[%d][%v](%q) push break", index, l.nodes[index].nodeType, l.nodes[index].value)
			_ = l.activeNodes.PushBack(newBreak)
		}
	}
}

func (l *LineBreaker) demeritsForRatio(typ nodeType, ratio, penalty float64) float64 {
	badness := l.badness * math.Pow(math.Abs(ratio), 3)
	baseDemerit := math.Pow(l.demerits.line+badness, 2)

	switch {
	case typ == nodeTypePenalty && penalty > 0:
		return baseDemerit + math.Pow(penalty, 2)
	case typ == nodeTypePenalty && penalty != -infinity:
		return baseDemerit - math.Pow(penalty, 2)
	default:
		return baseDemerit
	}
}

// fitnessClassForRatio establish a coarse fitness classification according to the cost ratio.
func fitnessClassForRatio(ratio float64) breakClass {
	switch {
	case ratio < -0.5:
		return breakClassZero
	case ratio <= 0.5:
		return breakClassOne
	case ratio <= 1:
		return breakClassTwo
	default:
		return breakClassThree
	}
}

// makeBreakPoints walks backwards the list of breakpoints,
// then returns ordered breakpoints as a slice.
func makeBreakPoints(brk *breakPoint) []breakPoint {
	breaks := []breakPoint{}

	for brk != nil {
		breaks = append(breaks, *brk)

		brk = brk.previous
	}

	return reverseBreakPoints(breaks)
}

func reverseBreakPoints(breaks []breakPoint) []breakPoint {
	reverse := make([]breakPoint, len(breaks))
	for i := range breaks {
		reverse[i] = breaks[len(breaks)-i-1]
	}

	return reverse
}

// getIdealWidth retrieves the constraint on the line length.
//
// NOTE: in this implementation, currentLine starts at 1.
func getIdealWidth(currentLine int, lengths []int) int {
	if currentLine < len(lengths)+1 {
		return lengths[currentLine-1]
	}

	return lengths[len(lengths)-1]
}
