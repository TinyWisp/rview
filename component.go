package rview

import (
	"reflect"

	"github.com/TinyWisp/rview/ddl"
)

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

func (c *ComponentInstance) GetCompField(field string) (interface{}, error) {
	val, err := GetStructField(c.Comp, field)

	return val, err
}

func (c *ComponentInstance) SetCompField(field string, val interface{}) error {
	return SetStructField(c.Comp, field, val)
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
	c.SetCompField("Inst", c)

	return nil
}

func (c *ComponentInstance) cloneCompMap() error {
	components, err := c.GetCompField("components")
	if err != nil {
		return err
	}

	c.compMap = components.(map[string]interface{})

	return nil
}

func (c *ComponentInstance) parseDdl() error {
	ddlInterface, err := c.GetCompField("Ddl")
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
