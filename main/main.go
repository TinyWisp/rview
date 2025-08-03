package main

import (
	"errors"

	"github.com/TinyWisp/rview/comp"
	"github.com/TinyWisp/rview/ddl"
	"github.com/davecgh/go-spew/spew"
	"github.com/rivo/tview"
)

var (
	tagCompCreatorMap = map[string]func() comp.Component{
		"box":  comp.CreateBox,
		"flex": comp.CreateFlex,
	}
)

func GenerateComponentNode(tplNode *ddl.TplNode) (*comp.ComponentNode, error) {
	compNode := comp.ComponentNode{}
	if tagCompCreator, ok := tagCompCreatorMap[tplNode.TagName]; ok {
		compNode.Comp = tagCompCreator()
	} else {
		return nil, errors.New("hello world")
	}

	if len(tplNode.Children) > 0 {
		for _, childTplNode := range tplNode.Children {
			childCompNode, err := GenerateComponentNode(childTplNode)
			if err != nil {
				return nil, err
			}
			childCompNode.Parent = &compNode
			compNode.Children = append(compNode.Children, childCompNode)
		}
	}

	return &compNode, nil
}

func GeneratePrimitiveFromDdl(ddlStr string) (tview.Primitive, error) {
	def, err := ddl.ParseDdl(ddlStr)
	if err != nil {
		return nil, err
	}
	mainTpl := def.TplMap["main"]

	rootCompNode, err2 := GenerateComponentNode(mainTpl)
	if err2 != nil {
		return nil, err2
	}

	return rootCompNode.Comp.GetPrimitive(), nil
}

func main() {
	ddlStr := `
		<template>
			<box></box>
		</template>
	`
	prim, err := GeneratePrimitiveFromDdl(ddlStr)
	if err != nil {
		spew.Println(err)
	}

	spew.Println(prim)
}
