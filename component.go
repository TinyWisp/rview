package rview

import (
	"github.com/gdamore/tcell"
)

const WIDTH_AUTO = -1
const HEIGHT_AUTO = -1
const RUNE_EMPTY = ' '

const (
	POSITION_ABSOLUTE = 1
	POSITION_RELATIVE = 2
	POSITION_FIXED = 3
)

type Cell struct {
	style tcell.Style
	ch rune
}

type Layout struct {
	marginLeft int
	marginTop int
	marginRight int
	marginBottom int

	borderLeft int
	borderTop int
	borderRight int
	borderBottom int

	paddingLeft int
	paddingTop int
	paddingRight int
	paddingBottom int

	width int
	height int

	scrollWidth int
	scrollHeight int
	scrollLeft int
	scrollTop int

	offsetLeft int
	offsetTop int
	pageLeft int
	pageTop int

	position int
	left int
	top int
}

type Style struct {
	display string
	position string
	width string
	height string
	backgroundColor string
	textColor string
}

type ComputedStyle struct {
	display int
	position int
	left int
	top int
}

type Component struct {
	name string
	components map[string]Component
	template string
	parent *Component
	children []*Component
	app *Component
	buffer [][]Cell
	width int
	height int
	layout Layout
	style Style
}

func (c *Component) Render() {
}

func (c *Component) Init() {

}

func (c *Component) InitBuffer() {
	colNum := 0
	rowNum := 0

	if (c.style.height != "auto") {
		rowNum = c.layout.height
	}

	if (c.style.width != "auto") {
		colNum = c.layout.width
	}

	c.buffer = make([][]Cell, rowNum)
	for y := range c.buffer {
		c.buffer[y] = make([]Cell, colNum)
		for x := range c.buffer[y] {
			c.buffer[y][x] = c.EmptyCell()
		}
	}
}

func (c *Component) ExtendBuffer(rowNum int, colNum int) {
	if (rowNum >= len(c.buffer)) {
		for i:=len(c.buffer); i<=rowNum; i++ {
			c.buffer[i] = make([]Cell, rowNum)
			for j:=0; j<rowNum; j++ {
				c.buffer[i][j] = c.EmptyCell()
			}
		}
	}
}

func (c *Component) calcWidth() {
	
}

func (c *Component) calcHeight() {

}

func (c *Component) calcPageX() {
	if c.style.position == "fixed" {
		c.layout.pageLeft = c.app.layout.scrollLeft + c.layout.offsetLeft
		c.layout.pageTop = c.app.layout.scrollTop + c.layout.offsetTop
	} else if c.style.position == "absolute" {
		c.layout.pageLeft = c.parent.layout.pageLeft + c.layout.offsetLeft
		c.layout.pageTop = c.parent.layout.pageTop + c.layout.offsetTop
	} else if c.style.position == "relative" {
		c.layout.pageLeft = c.parent.layout.pageLeft
	}
}

func (c *Component) calcPageY() {

}

func (c *Component) Merge() {

}

func (c *Component) EmptyCell() Cell {
	return Cell{
		ch: RUNE_EMPTY,
	}
}

func (c *Component) SetCell(x int, y int, cell Cell) {
	if (y >= len(c.buffer)) {
		appendRows := make([][]Cell, y - len(c.buffer) + 1)
		for i := range appendRows {
			appendRows[i] = make([]Cell, c.layout.width)
			if c.style.width != "auto" {
				for j:=0; j<c.layout.width; j++ {
					appendRows[i][j] = c.EmptyCell()
				}
			}
		}
		c.buffer = append(c.buffer, appendRows...)
	}

	if (x >= len(c.buffer[y])) {
		appendCols := make([]Cell, x - len(c.buffer[y]) + 1)
		for i := range appendCols {
			appendCols[i] = c.EmptyCell()
		}
		c.buffer[y] = append(c.buffer[y], appendCols...)
	}

	c.buffer[y][x] = cell

	if (c.layout.scrollWidth < x + 1) {
		c.layout.scrollWidth = x + 1;
	}
	if (c.layout.scrollHeight < y + 1) {
		c.layout.scrollHeight = y + 1;
	}
}

func (c *Component) SetContent(x int, y int, cnt string, style tcell.Style) {

}