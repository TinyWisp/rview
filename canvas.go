package rview

import (
	"github.com/TinyWisp/rview/ddl"
	"github.com/gdamore/tcell"
)

const WIDTH_AUTO = -1
const HEIGHT_AUTO = -1
const RUNE_EMPTY = 0

const (
	POSITION_ABSOLUTE = 1
	POSITION_RELATIVE = 2
	POSITION_FIXED    = 3
)

type Cell struct {
	style tcell.Style
	ch    rune
}

type Layout struct {
	MarginLeft   int
	MarginTop    int
	MarginRight  int
	MarginBottom int

	BorderLeft   int
	BorderTop    int
	BorderRight  int
	BorderBottom int

	PaddingLeft   int
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int

	Width  int
	Height int

	ScrollWidth  int
	ScrollHeight int
	ScrollLeft   int
	ScrollTop    int
	ScrollbarX  bool
	ScrollbarY  bool

	OffsetLeft int
	OffsetTop  int

	Position int
	Left     int
	Top      int
}

type Style struct {
	display         ddl.CSSToken
	position        ddl.CSSToken
	width           ddl.CSSToken
	height          ddl.CSSToken
	backgroundColor ddl.CSSToken
	textColor       ddl.CSSToken
}

type ComputedStyle struct {
	display  int
	position int
	left     int
	top      int
}

type Canvas struct {
	parent     *Canvas
	children   []*Canvas
	Layout     Layout
	buffer     [][]Cell
	fullBuffer [][]Cell
	style      Style
	cstyle     ComputedStyle
}

func (c *Canvas) InitBuffer() {
	colNum := c.calcWidth()
	rowNum := c.calcHeight()

	c.buffer = make([][]Cell, rowNum)
	for y := range c.buffer {
		c.buffer[y] = make([]Cell, colNum)
		for x := range c.buffer[y] {
			c.buffer[y][x] = c.EmptyCell()
		}
	}
}

func (c *Canvas) ExtendBuffer(rowNum int, colNum int) {
	if rowNum >= len(c.buffer) {
		for i := len(c.buffer); i <= rowNum; i++ {
			c.buffer[i] = make([]Cell, rowNum)
			for j := 0; j < rowNum; j++ {
				c.buffer[i][j] = c.EmptyCell()
			}
		}
	}
}

func (c *Canvas) calcWidth() int {
	width := c.style.width
	pLayout := c.parent.Layout
	if width.Type == ddl.CSSTokenNum {
		switch width.Unit {
		case ddl.CH:
			return int(width.Num)
		case ddl.PFW:
			return int(width.Num * float64(pLayout.ScrollWidth) / 100.0)
		case ddl.PCW:
			pcw := pLayout.ScrollWidth - pLayout.PaddingLeft - pLayout.PaddingRight - pLayout.BorderLeft - pLayout.BorderRight
			return int(width.Num * float64(pcw) / 100.0)
		case ddl.VW:
			vw := ScreenWidth
			return int(width.Num * float64(vw) / 100.0)
		}
	}

	return 0
}

func (c *Canvas) calcHeight() int {
	height := c.style.height
	pLayout := c.parent.Layout
	if height.Type == ddl.CSSTokenNum {
		switch height.Unit {
		case ddl.CH:
			return int(height.Num)
		case ddl.PFW:
			return int(height.Num * float64(pLayout.ScrollHeight) / 100.0)
		case ddl.PCW:
			pcw := pLayout.ScrollHeight - pLayout.PaddingTop - pLayout.PaddingBottom - pLayout.BorderTop - pLayout.BorderBottom
			return int(height.Num * float64(pcw) / 100.0)
		case ddl.VH:
			vh := ScreenHeight
			return int(height.Num * float64(vh) / 100.0)
		}
	}

	return 0
}

func (c *Canvas) calcPageY() {

}

func (c *Canvas) Merge() {

}

func (c *Canvas) EmptyCell() Cell {
	return Cell{
		ch: RUNE_EMPTY,
	}
}

func (c *Canvas) SetCell(x int, y int, cell Cell) {
	if y >= len(c.buffer) {
		appendRows := make([][]Cell, y-len(c.buffer)+1)
		for i := range appendRows {
			appendRows[i] = make([]Cell, c.Layout.Width)
			for j := 0; j < c.Layout.Width; j++ {
				appendRows[i][j] = c.EmptyCell()
			}
		}
		c.buffer = append(c.buffer, appendRows...)
	}

	if x >= len(c.buffer[y]) {
		appendCols := make([]Cell, x-len(c.buffer[y])+1)
		for i := range appendCols {
			appendCols[i] = c.EmptyCell()
		}
		c.buffer[y] = append(c.buffer[y], appendCols...)
	}

	c.buffer[y][x] = cell

	if c.Layout.ScrollWidth < x+1 {
		c.Layout.ScrollWidth = x + 1
	}
	if c.Layout.ScrollHeight < y+1 {
		c.Layout.ScrollHeight = y + 1
	}
}

func (c *Canvas) SetContent(x int, y int, cnt string, style tcell.Style) {

}
