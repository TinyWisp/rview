package rview

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview/comp"
	"github.com/TinyWisp/rview/ddl"
	"github.com/TinyWisp/rview/tperr"
	"github.com/rivo/tview"
)

type Page struct {
	Tpl               string
	tplNode           *ddl.TplNode
	TagCompCreatorMap map[string]func() comp.Component
	root              *ComponentNode
	def               interface{}
}

// get a variable for a node
func (p *Page) getVarForNode(node *ComponentNode, varName string) (interface{}, error) {
	structVar, err := GetStructField(p.def, varName)
	if err == nil {
		return structVar, nil
	}

	curNode := node
	for curNode != nil {
		if nodeVar, ok := curNode.Vars[varName]; ok {
			return nodeVar, nil
		}
		if !curNode.InheritVars {
			break
		}
		curNode = curNode.Parent
	}

	return nil, tperr.NewTypedError("page.undefinedVariable", varName)
}

func (p *Page) createCompNode(tplNode *ddl.TplNode, parent *ComponentNode) ([]*ComponentNode, error) {
	fmt.Println("aaaaaaaaaaaa")
	empty := []*ComponentNode{}

	// define a function that can get variables accessible for the parent node
	getParentVariable := func(varName string) (interface{}, error) {
		return p.getVarForNode(parent, varName)
	}

	fmt.Println("bbbbbbbbb")
	// create node
	keyPrefix := ""
	if parent != nil {
		keyPrefix = parent.Key
	}
	key := fmt.Sprintf("%s-%d", keyPrefix, tplNode.Idx)
	compNode := &ComponentNode{
		Key:         key,
		TplNode:     tplNode,
		Parent:      parent,
		Ignore:      false,
		Vars:        make(map[string]interface{}),
		InheritVars: true,
		HasIf:       false,
		HasElseIf:   false,
		HasElse:     false,
	}

	// v-for
	if tplNode.For != nil {
		iterateExp, err := CalcExp(tplNode.For.Range, getParentVariable)
		if err != nil {
			return empty, err
		}
		if iterateExp.Type != ddl.ExpInterface {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.For.Pos, "page.cannotRangeOverTheVar")
		}
		iterateVal := reflect.ValueOf(iterateExp.Interface)
		kind := iterateVal.Kind()
		itemVarName := tplNode.For.Item
		idxVarName := tplNode.For.Idx

		if kind == reflect.Array || kind == reflect.Slice {
			forNodes := []*ComponentNode{}
			for i := 0; i < iterateVal.Len(); i++ {
				itemVal := iterateVal.Index(i)
				copyCompNode := *compNode
				copyCompNode.Vars[itemVarName] = itemVal
				copyCompNode.Vars[idxVarName] = i
				copyTplNode := *tplNode
				copyTplNode.For = nil
				comp, cerr := p.createComponentAndSetProps(&copyCompNode, &copyTplNode)
				if cerr != nil {
					return empty, cerr
				}
				copyCompNode.Comp = comp
				forNodes = append(forNodes, &copyCompNode)
			}

			return forNodes, nil
		}
	}

	fmt.Println("cccccccccccccc")
	// v-if
	if tplNode.If != nil {
		compNode.HasIf = true
		fmt.Println("kkkkkkkkkkkkkkkkk")
		res, err := CalcExp(tplNode.If.Exp, getParentVariable)
		fmt.Println("llllllllllllll")
		if err != nil {
			return empty, err
		}
		if res.Type != ddl.ExpBool {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.If.Pos, "page.vifDirectriveMustBeBool", res.ActualTypeName())
		}
		compNode.If = res.Bool
		compNode.Ignore = !res.Bool

		// v-else-if
	} else if tplNode.ElseIf != nil {
		compNode.HasElseIf = true
		if tplNode.Idx == 0 {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Pos, "page.velseifHasNoCorrespondingIf")
		}
		prevCompNode := parent.Children[len(parent.Children)-1]
		if !prevCompNode.HasIf && !prevCompNode.HasElseIf {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Pos, "page.velseifHasNoCorrespondingIf")
		}
		if (prevCompNode.HasIf && prevCompNode.If) || (prevCompNode.HasElseIf && prevCompNode.ElseIf) {
			compNode.ElseIf = false
			compNode.Ignore = true
			return empty, nil
		}
		res, err := CalcExp(compNode.TplNode.ElseIf.Exp, getParentVariable)
		if err != nil {
			return empty, err
		}
		if res.Type != ddl.ExpBool {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.If.Pos, "page.velseifDirectriveMustBeBool", res.ActualTypeName())
		}
		compNode.ElseIf = res.Bool
		compNode.Ignore = !res.Bool

		// v-else
	} else if tplNode.Else != nil {
		compNode.HasElse = true
		if tplNode.Idx == 0 {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Pos, "page.velseHasNoCorrespondingIf")
		}
		prevCompNode := parent.Children[len(parent.Children)-1]
		if !prevCompNode.HasIf && !prevCompNode.HasElseIf {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Pos, "page.velseHasNoCorrespondingIf")
		}
		if (prevCompNode.HasIf && prevCompNode.If) || (prevCompNode.HasElseIf && prevCompNode.ElseIf) {
			compNode.Else = false
			compNode.Ignore = true
			return []*ComponentNode{compNode}, nil
		}
		compNode.Else = true
		compNode.Ignore = false
	}

	fmt.Println("dddddddddddddddd")
	if compNode.Ignore {
		return []*ComponentNode{compNode}, nil
	}

	fmt.Println("eeeeeeeeeeeee")
	// create component instance
	comp, cerr := p.createComponentAndSetProps(compNode, tplNode)
	if cerr != nil {
		return empty, cerr
	}

	compNode.Comp = comp

	fmt.Println("ffffffffff")
	// children
	for _, childTplNode := range tplNode.Children {
		childCompNodes, cerr := p.createCompNode(childTplNode, compNode)
		if cerr != nil {
			return empty, cerr
		}
		compNode.Children = append(compNode.Children, childCompNodes...)
	}

	fmt.Println("ggggggggggg")
	return []*ComponentNode{compNode}, nil
}

func (p *Page) createComponentAndSetProps(node *ComponentNode, tplNode *ddl.TplNode) (comp.Component, error) {
	tagCompCreator, ok := p.TagCompCreatorMap[tplNode.TagName]
	if !ok {
		err := ddl.NewDdlError(p.Tpl, tplNode.Pos, "page.compNotFound", tplNode.TagName)
		return nil, err
	}
	comp := tagCompCreator()

	// define a function to get variables
	getVariable := func(name string) (interface{}, error) {
		return p.getVarForNode(node, name)
	}

	// set the props
	for prop, attr := range tplNode.Attrs {
		exp, err := CalcExp(attr.Exp, getVariable)
		if err != nil {
			return nil, err
		}

		switch exp.Type {
		case ddl.ExpStr:
			err = comp.SetProp(prop, exp.Str)

		case ddl.ExpInt:
			err = comp.SetProp(prop, exp.Int)

		case ddl.ExpFloat:
			err = comp.SetProp(prop, exp.Float)

		case ddl.ExpBool:
			err = comp.SetProp(prop, exp.Bool)

		case ddl.ExpInterface:
			err = comp.SetProp(prop, exp.Interface)
		}

		if err != nil {
			return nil, err
		}
	}

	return comp, nil
}

func (p *Page) Primitive() tview.Primitive {
	return (p.root.Comp).Primitive()
}

func NewPage(def interface{}) (*Page, error) {
	p := &Page{
		def: def,
	}

	// Tpl
	itpl, err := GetStructField(p.def, "Tpl")
	if err != nil {
		return nil, tperr.NewTypedError("page.tplFieldIsRequired")
	}
	tpl, ok := itpl.(string)
	if !ok {
		return nil, tperr.NewTypedError("page.tplMustBeString")
	}
	p.Tpl = tpl

	// tplNode
	pddl, err := ddl.ParseDdl(p.Tpl)
	if err != nil {
		return nil, err
	}
	tplNode, ok := pddl.TplMap["main"]
	if !ok {
		return nil, tperr.NewTypedError("page.mainTemplateBeEssential")
	}
	p.tplNode = tplNode

	// TagCompCreatorMap
	icomponents, err := GetStructField(p.def, "Components")
	p.TagCompCreatorMap = map[string]func() comp.Component{
		"box":      comp.CreateBox,
		"button":   comp.CreateButton,
		"textarea": comp.CreateTextArea,
		"flex":     comp.CreateFlex,
		"template": comp.CreateTemplate,
	}
	if err == nil {
		tagCompCreatorMap, ok := icomponents.(map[string]func() comp.Component)
		if !ok {
			return nil, tperr.NewTypedError("page.invalidTypeOfComponentsField")
		}
		for k, v := range tagCompCreatorMap {
			p.TagCompCreatorMap[k] = v
		}
	}

	// root
	nodes, err := p.createCompNode(p.tplNode, nil)
	if err != nil {
		return nil, err
	}
	validRootNodeCount := 0
	for _, cnode := range nodes[0].Children {
		if !cnode.Ignore {
			validRootNodeCount += 1
		}
	}
	if validRootNodeCount == 0 {
		return nil, tperr.NewTypedError("page.tplMustContainOneRootNode")
	}
	if validRootNodeCount > 1 {
		return nil, tperr.NewTypedError("page.tplMustContainExactlyOneRootNode")
	}
	p.root = nodes[0]

	return p, nil
}
