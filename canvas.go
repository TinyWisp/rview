package rview

import (
	"github.com/TinyWisp/rview/ddl"
	"github.com/gdamore/tcell"
)

const WIDTH_AUTO = -1
const HEIGHT_AUTO = -1
const RUNE_EMPTY = 0
const UNLIMITED = -1

const (
	POSITION_ABSOLUTE = 1
	POSITION_RELATIVE = 2
	POSITION_FIXED    = 3
)

type Cell struct {
	Style tcell.Style
	Char  rune
}

var EmptyCell = Cell{
	Char: RUNE_EMPTY,
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

	ContentWidth        int
	ContentHeight       int
	ActualContentWidth  int
	ActualContentHeight int

	ScrollWidth  int
	ScrollHeight int
	ScrollLeft   int
	ScrollTop    int
	ScrollbarX   bool
	ScrollbarY   bool

	OffsetLeft int
	OffsetTop  int

	Position int
	Left     int
	Top      int
}

type Style struct {
	Display           ddl.CSSToken
	Position          ddl.CSSToken
	Width             ddl.CSSToken
	Height            ddl.CSSToken
	BackgroundColor   ddl.CSSToken
	TextColor         ddl.CSSToken
	BorderLeftColor   ddl.CSSToken
	BorderRightColor  ddl.CSSToken
	BorderTopColor    ddl.CSSToken
	BorderBottomColor ddl.CSSToken
	BorderLeftWidth   ddl.CSSToken
	BorderRightWidth  ddl.CSSToken
	BorderTopWidth    ddl.CSSToken
	BorderBottomWidth ddl.CSSToken
	BorderChar        []ddl.CSSToken
	PaddingLeft       ddl.CSSToken
	PaddingTop        ddl.CSSToken
	PaddingRight      ddl.CSSToken
	PaddingBottom     ddl.CSSToken
}

type ComputedStyle struct {
	display  int
	position int
	left     int
	top      int
}

type Canvas struct {
	parent    *Canvas
	children  []*Canvas
	Layout    Layout
	buffer    [][]Cell
	cntBuffer [][]Cell
	style     Style
	cstyle    ComputedStyle
}

func (c *Canvas) Init() {
	c.initLayout()
	c.InitBuffer()
	c.InitContentBuffer()
}

func (c *Canvas) InitBuffer() {
	colNum := c.Layout.Width
	rowNum := c.Layout.Height

	c.buffer = make([][]Cell, rowNum)
	for y := range c.buffer {
		c.buffer[y] = make([]Cell, colNum)
		for x := range c.buffer[y] {
			c.buffer[y][x] = EmptyCell
		}
	}
}

func (c *Canvas) InitContentBuffer() {
	colNum := c.Layout.ContentWidth
	rowNum := c.Layout.ContentHeight

	c.cntBuffer = make([][]Cell, rowNum)
	for y := range c.cntBuffer {
		c.cntBuffer[y] = make([]Cell, colNum)
		for x := range c.buffer[y] {
			c.cntBuffer[y][x] = EmptyCell
		}
	}
}

func (c *Canvas) initLayout() {
	c.Layout.Width, _ = c.calcChNum(c.style.Width)
	c.Layout.Height, _ = c.calcChNum(c.style.Height)
	c.Layout.BorderLeft, _ = c.calcChNum(c.style.BorderLeftWidth)
	c.Layout.BorderRight, _ = c.calcChNum(c.style.BorderRightWidth)
	c.Layout.BorderTop, _ = c.calcChNum(c.style.BorderTopWidth)
	c.Layout.BorderBottom, _ = c.calcChNum(c.style.BorderBottomWidth)
	c.Layout.PaddingLeft, _ = c.calcChNum(c.style.PaddingLeft)
	c.Layout.PaddingRight, _ = c.calcChNum(c.style.PaddingRight)
	c.Layout.PaddingTop, _ = c.calcChNum(c.style.PaddingTop)
	c.Layout.PaddingBottom, _ = c.calcChNum(c.style.PaddingBottom)
	c.Layout.ContentWidth = c.Layout.Width - c.Layout.BorderLeft - c.Layout.BorderRight - c.Layout.PaddingLeft - c.Layout.PaddingRight
	c.Layout.ContentHeight = c.Layout.Height - c.Layout.BorderTop - c.Layout.BorderBottom - c.Layout.PaddingTop - c.Layout.PaddingBottom

	if c.Layout.ScrollbarX {
		c.Layout.ContentHeight -= 1
	}
	if c.Layout.ContentHeight < 0 {
		c.Layout.ContentHeight = 0
	}

	if c.Layout.ScrollbarY {
		c.Layout.ContentWidth -= 1
	}
	if c.Layout.ContentWidth < 0 {
		c.Layout.ContentWidth = 0
	}
}

func (c *Canvas) calcChNum(token ddl.CSSToken) (int, error) {
	if token.Type == ddl.CSSTokenNum {
		pLayout := c.parent.Layout
		switch token.Unit {
		case ddl.CH:
			return int(token.Num), nil
		case ddl.PFW:
			return int(token.Num * float64(pLayout.ScrollWidth) / 100.0), nil
		case ddl.PFH:
			return int(token.Num * float64(pLayout.ScrollHeight) / 100.0), nil
		case ddl.PCW:
			return int(token.Num * float64(pLayout.ContentWidth) / 100.0), nil
		case ddl.PCH:
			return int(token.Num * float64(pLayout.ContentHeight) / 100.0), nil
		case ddl.VW:
			return int(token.Num * float64(ScreenWidth) / 100.0), nil
		case ddl.VH:
			return int(token.Num * float64(ScreenHeight) / 100.0), nil
		}
	}

	return 0, NewError("canvas.NotANumberToken")
}

func (c *Canvas) calcPageY() {

}

func (c *Canvas) Merge() {

}

func (c *Canvas) SetContentCell(x int, y int, cell Cell) {
	if y >= len(c.cntBuffer) {
		appendRows := make([][]Cell, y-len(c.cntBuffer)+1)
		c.cntBuffer = append(c.cntBuffer, appendRows...)
	}

	if x >= len(c.cntBuffer[y]) {
		appendColNum := (x - len(c.cntBuffer[y]) + 1 + 19) / 20 * 20
		appendCols := make([]Cell, appendColNum)
		for i, _ := range appendCols {
			appendCols[i] = EmptyCell
		}
		c.cntBuffer[y] = append(c.cntBuffer[y], appendCols...)
	}

	c.cntBuffer[y][x] = cell

	if x+1 > c.Layout.ActualContentWidth {
		c.Layout.ActualContentWidth = x + 1
	}

	if y+1 > c.Layout.ActualContentHeight {
		c.Layout.ActualContentHeight = y + 1
	}
}

func (c *Canvas) WriteContentText(x int, y int, cnt string, style tcell.Style) {
	if c.Layout.ContentWidth == 0 || c.Layout.ContentHeight == 0 {
		return
	}

	if c.Layout.ContentWidth == UNLIMITED {
		for idx, ch := range cnt {
			cell := Cell{
				Style: style,
				Char:  ch,
			}
			c.SetContentCell(x+idx, y, cell)
		}
		return
	}

	for idx, ch := range cnt {
		cell := Cell{
			Style: style,
			Char:  ch,
		}
		nx := (x + idx) % c.Layout.ContentWidth
		ny := y + (x+idx)/c.Layout.ContentWidth
		c.SetContentCell(nx, ny, cell)
	}
}

func (c *Canvas) SetCell(x int, y int, cell Cell) {
	if y >= len(c.buffer) {
		appendRows := make([][]Cell, y-len(c.buffer)+1)
		c.buffer = append(c.buffer, appendRows...)
	}

	if x >= len(c.buffer[y]) {
		appendColNum := (x - len(c.buffer[y]) + 1 + 19) / 20 * 20
		appendCols := make([]Cell, appendColNum)
		for i, _ := range appendCols {
			appendCols[i] = EmptyCell
		}
		c.buffer[y] = append(c.buffer[y], appendCols...)
	}

	c.buffer[y][x] = cell

	if x+1 > c.Layout.Width {
		c.Layout.Width = x + 1
	}

	if y+1 > c.Layout.Height {
		c.Layout.Height = y + 1
	}
}

func (c *Canvas) DrawBorder() {
	baseStyle := tcell.StyleDefault.Background(tcell.Color(c.style.BackgroundColor.IntColor))
	topStyle := baseStyle.Foreground(tcell.Color(c.style.BorderTopColor.IntColor))
	bottomStyle := baseStyle.Foreground(tcell.Color(c.style.BorderBottomColor.IntColor))
	leftStyle := baseStyle.Foreground(tcell.Color(c.style.BorderLeftColor.IntColor))
	rightStyle := baseStyle.Foreground(tcell.Color(c.style.BorderRightColor.IntColor))
	leftTopStyle := topStyle
	rightTopStyle := topStyle
	leftBottomStyle := bottomStyle
	rightBottomStyle := bottomStyle

	leftTopChar := []rune(c.style.BorderChar[0].Str)[0]
	topChar := []rune(c.style.BorderChar[1].Str)[0]
	rightTopChar := []rune(c.style.BorderChar[2].Str)[0]
	rightChar := []rune(c.style.BorderChar[3].Str)[0]
	rightBottomChar := []rune(c.style.BorderChar[4].Str)[0]
	bottomChar := []rune(c.style.BorderChar[5].Str)[0]
	leftBottomChar := []rune(c.style.BorderChar[6].Str)[0]
	leftChar := []rune(c.style.BorderChar[7].Str)[0]

	if c.style.BorderLeftWidth.Num > 0 && c.style.BorderTopWidth.Num == 0 {
		leftTopChar = leftChar
		leftTopStyle = leftStyle
	} else if c.style.BorderLeftWidth.Num == 0 && c.style.BorderTopWidth.Num > 0 {
		leftTopChar = topChar
		leftTopStyle = topStyle
	}

	if c.style.BorderTopWidth.Num > 0 && c.style.BorderRightWidth.Num == 0 {
		rightTopChar = topChar
		rightTopStyle = topStyle
	} else if c.style.BorderTopWidth.Num == 0 && c.style.BorderRightWidth.Num > 0 {
		rightTopChar = rightChar
		rightTopStyle = rightStyle
	}

	if c.style.BorderBottomWidth.Num > 0 && c.style.BorderRightWidth.Num == 0 {
		rightBottomChar = bottomChar
		rightBottomStyle = bottomStyle
	} else if c.style.BorderBottomWidth.Num == 0 && c.style.BorderRightWidth.Num > 0 {
		rightBottomChar = rightChar
		rightBottomStyle = rightStyle
	}

	if c.style.BorderBottomWidth.Num > 0 && c.style.BorderLeftWidth.Num == 0 {
		leftBottomChar = bottomChar
		leftBottomStyle = bottomStyle
	} else if c.style.BorderBottomWidth.Num == 0 && c.style.BorderLeftWidth.Num > 0 {
		leftBottomChar = leftChar
		leftBottomStyle = leftStyle
	}

	c.SetCell(0, 0, Cell{
		Style: leftTopStyle,
		Char:  leftTopChar,
	})
	if c.style.BorderTopWidth.Num > 0 {
		for i := 1; i < c.Layout.Width-1; i++ {
			c.SetCell(i, 0, Cell{
				Style: topStyle,
				Char:  topChar,
			})
		}
	}
	c.SetCell(c.Layout.Width-1, 0, Cell{
		Style: rightTopStyle,
		Char:  rightTopChar,
	})
	if c.style.BorderRightWidth.Num > 0 {
		for i := 1; i < c.Layout.Height-1; i++ {
			c.SetCell(c.Layout.Width-1, i, Cell{
				Style: rightStyle,
				Char:  rightChar,
			})
		}
	}
	c.SetCell(c.Layout.Width-1, c.Layout.Height-1, Cell{
		Style: rightBottomStyle,
		Char:  rightBottomChar,
	})
	if c.style.BorderBottomWidth.Num > 0 {
		for i := 1; i < c.Layout.Width-1; i++ {
			c.SetCell(i, c.Layout.Height-1, Cell{
				Style: bottomStyle,
				Char:  bottomChar,
			})
		}
	}
	c.SetCell(0, c.Layout.Height-1, Cell{
		Style: leftBottomStyle,
		Char:  leftBottomChar,
	})
	if c.style.BorderLeftWidth.Num > 0 {
		for i := 1; i < c.Layout.Height-1; i++ {
			c.SetCell(0, i, Cell{
				Style: leftStyle,
				Char:  leftChar,
			})
		}
	}
}

func (c *Canvas) calcScrollWidth() int {
	width := c.Layout.ActualContentWidth + c.Layout.BorderLeft + c.Layout.BorderRight + c.Layout.PaddingLeft + c.Layout.PaddingRight
	if c.Layout.ScrollbarY {
		width += 1
	}
	if width < c.Layout.Width {
		width = c.Layout.Width
	}
	return width
}

func (c *Canvas) calcScrollHeight() int {
	height := c.Layout.ActualContentHeight + c.Layout.BorderTop + c.Layout.BorderBottom + c.Layout.PaddingTop + c.Layout.PaddingBottom
	if c.Layout.ScrollbarY {
		height += 1
	}
	if height < c.Layout.Height {
		height = c.Layout.Height
	}
	return height
}
