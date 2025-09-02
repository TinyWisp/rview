package comp

import (
	"fmt"

	"github.com/TinyWisp/rview"
	"github.com/TinyWisp/rview/ddl"
	"github.com/rivo/tview"
)

type Page struct {
	Tpl               string
	tplNode           *ddl.TplNode
	TagCompCreatorMap map[string]func() Component
	root              *ComponentNode
}

func (p *Page) genCompNode(tplNode *ddl.TplNode) (*ComponentNode, error) {
	compNode := ComponentNode{}
	if tagCompCreator, ok := p.TagCompCreatorMap[tplNode.TagName]; ok {
		compNode.Comp = tagCompCreator()
		for prop, val := range tplNode.Attrs {
			err := SetProp(compNode.Comp.Primitive(), prop, val)
			if err != nil {
				return nil, err
			}
		}
		for prop, attr := range tplNode.Attrs {
			prim := compNode.Comp.Primitive()
			val := interface{}(nil)
			exp := attr.Exp
			switch exp.Type {
			case ddl.ExpBool:
				val = exp.Bool

			case ddl.ExpFloat:
				val = exp.Float

			case ddl.ExpInt:
				val = exp.Int

			case ddl.ExpStr:
				val = exp.Str
			}
			err := SetProp(prim, prop, val)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		return nil, fmt.Errorf("tag:%s", tplNode.TagName)
	}

	if len(tplNode.Children) > 0 {
		for _, childTplNode := range tplNode.Children {
			childCompNode, err := p.genCompNode(childTplNode)
			if err != nil {
				return nil, err
			}
			childCompNode.Parent = &compNode
			compNode.Children = append(compNode.Children, childCompNode)
		}
	}

	return &compNode, nil
}

// get a variable for a node
func (p *Page) getVarForNode(node *ComponentNode, varName string) (interface{}, error) {
	structVar, err := rview.GetStructField(p, varName)
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
		curNode = node.Parent
	}

	return nil, rview.NewTypedError("comp.undefinedVariable", varName)
}

func (p *Page) createCompNode(tplNode *ddl.TplNode, parent *ComponentNode) (*ComponentNode, error) {
	// define a function that can get variables accessible for the parent node
	getParentVariable := func(varName string) (interface{}, error) {
		return p.getVarForNode(parent, varName)
	}

	// create node
	keyPrefix := ""
	if parent != nil {
		keyPrefix = parent.Key
	}
	key := fmt.Sprintf("%s-%d", keyPrefix, tplNode.Idx)
	node := &ComponentNode{
		Key:         key,
		TplNode:     tplNode,
		Parent:      parent,
		IsVoid:      false,
		Vars:        make(map[string]interface{}),
		InheritVars: true,
		HasIf:       false,
		HasElseIf:   false,
		HasElse:     false,
	}

	// v-if
	if tplNode.If != nil {
		node.HasIf = true
		res, err := rview.CalcExp(tplNode.If.Exp, getParentVariable)
		if err != nil {
			return nil, err
		}
		if res.Type != ddl.ExpBool {
			return nil, ddl.NewDdlError(p.Tpl, tplNode.If.Pos, "comp.vifDirectriveMustBeBool")
		}
		if res.Type == ddl.ExpBool && res.Bool == false {
			node.If = false
			node.IsVoid = true
			return nil, nil
		}
		if res.Type == ddl.ExpBool && res.Bool == true {
			node.If = true
			node.IsVoid = false
		}

		// v-else-if
	} else if tplNode.ElseIf != nil {
		node.ElseIf = true
		if len(parent.Children) == 0 {
			return nil, rview.NewTypedError("comp.noCorrespondingIf")
		}
		prevCompNode := parent.Children[len(parent.Children)-1]
		if !prevCompNode.HasIf && !prevCompNode.HasElseIf {
			return nil, rview.NewTypedError("comp.noCorrespondingIf")
		}
		if prevCompNode.HasIf {
			if prevCompNode.If == true {
				node.IsVoid = true
			} else {
				res, err := rview.CalcExp(node.TplNode.ElseIf.Exp, getParentVariable)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	tagCompCreator, ok := p.TagCompCreatorMap[tplNode.TagName]
	if !ok {
		return nil, fmt.Errorf("tag:%s", tplNode.TagName)
	}
	compNode.Comp = tagCompCreator()

	// define the function to get variables
	getVariable := func(name string) (interface{}, error) {
		structVar, err := rview.GetStructField(p, name)
		nodeVar, ok := compNode.InheritableVars[name]

		if err == nil && ok {
			return nil, rview.NewTypedError("comp.duplicateVariable", name)
		}

		if err != nil && !ok {
			return nil, rview.NewTypedError("comp.undefinedVariable", name)
		}

		if err == nil {
			return structVar, nil
		}

		return nodeVar, nil
	}

	// set the common props
	for prop, val := range tplNode.Attrs {
		err := SetProp(compNode.Comp.Primitive(), prop, val)
		if err != nil {
			fmt.Println(err)
		}
	}

	// set the props defined with v-bind: or :
	for prop, exp := range tplNode.Binds {
		val := interface{}(nil)
		switch exp.Type {
		case ddl.ExpBool:
			val = exp.Bool

		case ddl.ExpFloat:
			val = exp.Float

		case ddl.ExpInt:
			val = exp.Int

		case ddl.ExpStr:
			val = exp.Str

		case ddl.ExpNil:
			val = nil

		case ddl.ExpVar:
			resVal, err := rview.CalcExp(exp, getVariable)
			if err != nil {
				return nil, err
			}
			val = resVal

		case ddl.ExpFunc:
			resVal, err := rview.CalcExp(exp, getVariable)
			if err != nil {
				return nil, err
			}
			val = resVal
		}
		err := compNode.Comp.SetProp(prop, val)
		if err != nil {
			return nil, err
		}
	}

	if len(tplNode.Children) > 0 {
		for _, childTplNode := range tplNode.Children {
			childCompNode, err := p.genCompNode(childTplNode)
			if err != nil {
				return nil, err
			}
			childCompNode.Parent = &compNode
			compNode.Children = append(compNode.Children, childCompNode)
		}
	}

	return &compNode, nil
}

func (p *Page) Init() error {
	def, err := ddl.ParseDdl(p.Tpl)
	if err != nil {
		return err
	}
	mainTpl := def.TplMap["main"]
	p.tplNode = mainTpl

	rootCompNode, err2 := p.genCompNode(p.tplNode.Children[0])
	if err2 != nil {
		return err2
	}
	p.root = rootCompNode

	return nil
}

func (p *Page) Primitive() tview.Primitive {
	return p.root.Comp.Primitive()
}
