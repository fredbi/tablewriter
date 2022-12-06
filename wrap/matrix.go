package wrap

import "sort"

type (
	columns []column
	cells   []cell

	column struct {
		i        int
		maxWidth int
		cells    cells
	}

	cell struct {
		i       int
		j       int
		content *string
		pvalues []int
		width   int
		passNo  int
	}
)

func (c columns) Less(i, j int) bool {
	return c[i].maxWidth > c[j].maxWidth
}

func (c columns) Len() int {
	return len(c)
}

func (c columns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c cells) Less(i, j int) bool {
	return c[i].width > c[j].width
}

func (c cells) Len() int {
	return len(c)
}

func (c cells) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c column) SortRows() {
	sort.Sort(c.cells)
}
