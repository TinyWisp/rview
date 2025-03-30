package rview

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview/ddl"
	"github.com/gdamore/tcell"
)

const WIDTH_AUTO = -1
const HEIGHT_AUTO = -1
const RUNE_EMPTY = ' '

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
	marginLeft   int
	marginTop    int
	marginRight  int
	marginBottom int

	borderLeft   int
	borderTop    int
	borderRight  int
	borderBottom int

	paddingLeft   int
	paddingTop    int
	paddingRight  int
	paddingBottom int

	width  int
	height int

	scrollWidth  int
	scrollHeight int
	scrollLeft   int
	scrollTop    int

	offsetLeft int
	offsetTop  int
	pageLeft   int
	pageTop    int

	position int
	left     int
	top      int
}

type Style struct {
	display         string
	position        string
	width           string
	height          string
	backgroundColor string
	textColor       string
}

type ComputedStyle struct {
	display  int
	position int
	left     int
	top      int
}

type ComponentInstance struct {
	Comp        interface{}
	compMap     map[string]interface{}
	ddl         string
	Parent      *ComponentInstance
	Children    []*ComponentInstance
	App         *ComponentInstance
	tplMap      map[string]*ddl.TplNode
	cssClassMap ddl.CSSClassMap
	prevTree    *ddl.TplNode
	curTree     *ddl.TplNode
	root        *ComponentInstance
	defaultSlot []*ComponentInstance
	cache       map[string]*ComponentInstance
	buffer      [][]Cell
	width       int
	height      int
	layout      Layout
	style       Style
}

type Component struct {
	Name string
	Ddl  string
	Inst *ComponentInstance
}

type Initable interface {
	Init()
}

type Renderable interface {
	Render()
}

func (c *ComponentInstance) Render() {
}

func (c *ComponentInstance) GetCompProp(field string, expectType interface{}) (interface{}, error) {
	comp := reflect.ValueOf(c.Comp).Elem()
	val := comp.FieldByName(field)

	if !val.IsValid() {
		return nil, NewError("comp.propNotExist")
	}

	if reflect.TypeOf(val) != reflect.TypeOf(expectType) {
		return nil, NewError("comp.typeNotAsExpected")
	}

	return val.Interface(), nil
}

func (c *ComponentInstance) SetCompProp(field string, val interface{}) error {
	comp := reflect.ValueOf(c.Comp).Elem()
	curVal := comp.FieldByName(field)

	if !curVal.IsValid() {
		return NewError("comp.SetCompProp.propNotExist")
	}

	if curVal.Type() != reflect.TypeOf(val) {
		fmt.Println(reflect.TypeOf(curVal))
		fmt.Println(reflect.TypeOf(val))
		return NewError("comp.SetCompProp.typeMismatch %s %s", reflect.TypeOf(curVal), reflect.TypeOf(val))
	}

	/*
		if !curVal.CanSet() {
			return NewError("comp.SetCompProp.cannotSetFieldValue")
		}
	*/

	curVal.Set(reflect.ValueOf(val))
	return nil
}

func (c *ComponentInstance) CreateChildInstance(node *ddl.TplNode, parent *ComponentInstance) (*ComponentInstance, error) {
	tagName := node.TagName
	instance := &ComponentInstance{}
	if comp, ok := c.compMap[tagName]; ok {
		ncomp := reflect.New(reflect.TypeOf(comp)).Interface()
		instance.Comp = ncomp
		instance.Parent = parent
		instance.Init()
	} else {
		return nil, NewError("comp.cannotResolveComponent", tagName)
	}

	cinsts := []*ComponentInstance{}
	for _, child := range node.Children {
		if cinst, err := c.CreateChildInstance(child, instance); err == nil {
			cinsts = append(cinsts, cinst)
		} else {
			return nil, err
		}
	}
	instance.defaultSlot = cinsts
	instance.Children = append(instance.Children, cinsts...)

	return instance, nil
}

func (c *ComponentInstance) Init() error {
	if initableComp, ok := c.Comp.(Initable); ok {
		initableComp.Init()
	}

	if err := c.cloneCompMap(); err != nil {
		return err
	}

	if err2 := c.parseDdl(); err2 != nil {
		return err2
	}

	_, err3 := c.CreateChildInstance(c.tplMap["main"], c)
	if err3 != nil {
		return err3
	}
	c.SetCompProp("Inst", c)

	return nil
}

func (c *ComponentInstance) cloneCompMap() error {
	components, err := c.GetCompProp("components", map[string]interface{}{})
	if err != nil {
		return err
	}

	c.compMap = components.(map[string]interface{})

	return nil
}

func (c *ComponentInstance) parseDdl() error {
	ddlInterface, err := c.GetCompProp("Ddl", "")
	if err != nil {
		return err
	}

	ddlStr := ddlInterface.(string)
	def, err2 := ddl.ParseDdl(ddlStr)
	if err2 != nil {
		return err2
	}

	c.tplMap = def.TplMap
	c.cssClassMap = def.CssClassMap

	return nil
}

func (c *ComponentInstance) InitBuffer() {
	colNum := 0
	rowNum := 0

	if c.style.height != "auto" {
		rowNum = c.layout.height
	}

	if c.style.width != "auto" {
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

func (c *ComponentInstance) ExtendBuffer(rowNum int, colNum int) {
	if rowNum >= len(c.buffer) {
		for i := len(c.buffer); i <= rowNum; i++ {
			c.buffer[i] = make([]Cell, rowNum)
			for j := 0; j < rowNum; j++ {
				c.buffer[i][j] = c.EmptyCell()
			}
		}
	}
}

func (c *ComponentInstance) calcWidth() {

}

func (c *ComponentInstance) calcHeight() {

}

func (c *ComponentInstance) calcPageX() {
	if c.style.position == "fixed" {
		c.layout.pageLeft = c.App.layout.scrollLeft + c.layout.offsetLeft
		c.layout.pageTop = c.App.layout.scrollTop + c.layout.offsetTop
	} else if c.style.position == "absolute" {
		c.layout.pageLeft = c.Parent.layout.pageLeft + c.layout.offsetLeft
		c.layout.pageTop = c.Parent.layout.pageTop + c.layout.offsetTop
	} else if c.style.position == "relative" {
		c.layout.pageLeft = c.Parent.layout.pageLeft
	}
}

func (c *ComponentInstance) calcPageY() {

}

func (c *ComponentInstance) Merge() {

}

func (c *ComponentInstance) EmptyCell() Cell {
	return Cell{
		ch: RUNE_EMPTY,
	}
}

func (c *ComponentInstance) SetCell(x int, y int, cell Cell) {
	if y >= len(c.buffer) {
		appendRows := make([][]Cell, y-len(c.buffer)+1)
		for i := range appendRows {
			appendRows[i] = make([]Cell, c.layout.width)
			if c.style.width != "auto" {
				for j := 0; j < c.layout.width; j++ {
					appendRows[i][j] = c.EmptyCell()
				}
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

	if c.layout.scrollWidth < x+1 {
		c.layout.scrollWidth = x + 1
	}
	if c.layout.scrollHeight < y+1 {
		c.layout.scrollHeight = y + 1
	}
}

func (c *ComponentInstance) SetContent(x int, y int, cnt string, style tcell.Style) {

}
