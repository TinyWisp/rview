package rview

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/TinyWisp/rview/comp"
	"github.com/TinyWisp/rview/ddl"
	"github.com/TinyWisp/rview/tperr"
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
	empty := []*ComponentNode{}

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
			return empty, ddl.NewDdlError(p.Tpl, tplNode.For.RangePos, "page.cannotIterateOverTheVar")
		}

		iterateVal := reflect.ValueOf(iterateExp.Interface)
		kind := iterateVal.Kind()
		itemVarName := tplNode.For.Val
		idxVarName := tplNode.For.Idx

		// array, slice
		if kind == reflect.Array || kind == reflect.Slice {
			forNodes := []*ComponentNode{}
			for i := 0; i < iterateVal.Len(); i++ {
				itemVal := iterateVal.Index(i)
				copyTplNode := *tplNode
				copyTplNode.For = nil
				copyCompNode := *compNode
				copyCompNode.Key = fmt.Sprintf("%s-%d", keyPrefix, i)
				copyCompNode.Vars = map[string]interface{}{}
				copyCompNode.Vars[itemVarName] = itemVal.Interface()
				copyCompNode.Vars[idxVarName] = i
				comp, cerr := p.createComponentAndSetProps(&copyCompNode, &copyTplNode)
				if cerr != nil {
					return empty, cerr
				}
				copyCompNode.Comp = comp

				if len(copyTplNode.Children) > 0 {
					for _, childTplNode := range copyTplNode.Children {
						childComponentNode, ccerr := p.createCompNode(childTplNode, &copyCompNode)
						if ccerr != nil {
							return empty, ccerr
						}
						copyCompNode.Children = append(copyCompNode.Children, childComponentNode...)
					}
				}

				forNodes = append(forNodes, &copyCompNode)
			}

			return forNodes, nil
		}

		// map
		if kind == reflect.Map {
			forNodes := []*ComponentNode{}
			mapKeys := iterateVal.MapKeys()

			// golang's maps are unordered.
			// to avoid inconsistencies in the order of generated nodes each time, sorting is necessary.
			sort.SliceStable(mapKeys, func(i int, j int) bool {
				switch mapKeys[i].Kind() {
				case reflect.String:
					return mapKeys[i].String() < mapKeys[j].String()

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					return mapKeys[i].Int() < mapKeys[j].Int()

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return mapKeys[i].Uint() < mapKeys[j].Uint()

				case reflect.Float32, reflect.Float64:
					return mapKeys[i].Float() < mapKeys[j].Float()

				default:
					return false
				}
			})

			for i := 0; i < len(mapKeys); i++ {
				mkey := mapKeys[i]
				mval := iterateVal.MapIndex(mkey)
				copyTplNode := *tplNode
				copyTplNode.For = nil
				copyCompNode := *compNode
				copyCompNode.Key = fmt.Sprintf("%s-%d", keyPrefix, i)
				copyCompNode.Vars = map[string]interface{}{}
				copyCompNode.Vars[itemVarName] = mval.Interface()
				copyCompNode.Vars[idxVarName] = mkey.Interface()
				comp, cerr := p.createComponentAndSetProps(&copyCompNode, &copyTplNode)
				if cerr != nil {
					return empty, cerr
				}
				copyCompNode.Comp = comp

				if len(copyTplNode.Children) > 0 {
					for _, childTplNode := range copyTplNode.Children {
						childComponentNode, ccerr := p.createCompNode(childTplNode, &copyCompNode)
						if ccerr != nil {
							return empty, ccerr
						}
						copyCompNode.Children = append(copyCompNode.Children, childComponentNode...)
					}
				}

				forNodes = append(forNodes, &copyCompNode)
			}

			return forNodes, nil
		}
	}

	// v-if
	if tplNode.If != nil {
		compNode.HasIf = true
		res, err := CalcExp(tplNode.If.Exp, getParentVariable)
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
			return empty, ddl.NewDdlError(p.Tpl, tplNode.ElseIf.Pos, "page.velseifHasNoCorrespondingIf")
		}
		prevCompNode := parent.Children[len(parent.Children)-1]
		if !prevCompNode.HasIf && !prevCompNode.HasElseIf {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.ElseIf.Pos, "page.velseifHasNoCorrespondingIf")
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
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Else.Pos, "page.velseHasNoCorrespondingIf")
		}
		prevCompNode := parent.Children[len(parent.Children)-1]
		if !prevCompNode.HasIf && !prevCompNode.HasElseIf {
			return empty, ddl.NewDdlError(p.Tpl, tplNode.Else.Pos, "page.velseHasNoCorrespondingIf")
		}
		if (prevCompNode.HasIf && prevCompNode.If) || (prevCompNode.HasElseIf && prevCompNode.ElseIf) {
			compNode.Else = false
			compNode.Ignore = true
			return []*ComponentNode{compNode}, nil
		}
		compNode.Else = true
		compNode.Ignore = false
	}

	if compNode.Ignore {
		return []*ComponentNode{compNode}, nil
	}

	// create component instance
	comp, cerr := p.createComponentAndSetProps(compNode, tplNode)
	if cerr != nil {
		return empty, cerr
	}

	compNode.Comp = comp

	// children
	for _, childTplNode := range tplNode.Children {
		childCompNodes, cerr := p.createCompNode(childTplNode, compNode)
		if cerr != nil {
			return empty, cerr
		}
		compNode.Children = append(compNode.Children, childCompNodes...)
	}

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
			if terr, ok := err.(*tperr.TypedError); ok {
				derr := ddl.NewDdlError(p.Tpl, attr.Pos, terr.GetEtype(), terr.GetVars()...)
				return nil, derr
			}

			return nil, err
		}
	}

	return comp, nil
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
