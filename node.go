package rview

import (
	"github.com/TinyWisp/rview/comp"
	"github.com/TinyWisp/rview/ddl"
)

type ComponentNode struct {
	Comp        comp.Component
	Parent      *ComponentNode
	Children    []*ComponentNode
	Key         string
	TplNode     *ddl.TplNode
	Vars        map[string]interface{}
	InheritVars bool
	Ignore      bool
	HasIf       bool
	If          bool
	HasElseIf   bool
	ElseIf      bool
	HasElse     bool
	Else        bool
	HasFor      bool
}
