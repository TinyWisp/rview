package comp

import (
	"github.com/TinyWisp/rview/ddl"
	"github.com/rivo/tview"
)

type Component interface {
	Primitive() tview.Primitive
	CanAddItem() bool
	SetProp(string, interface{}) error
}

type ComponentNode struct {
	Comp        Component
	Parent      *ComponentNode
	Children    []*ComponentNode
	Key         string
	TplNode     *ddl.TplNode
	Vars        map[string]interface{}
	InheritVars bool
	IsVoid      bool
	HasIf       bool
	If          bool
	HasElseIf   bool
	ElseIf      bool
	HasElse     bool
	Else        bool
}
