package comp

import (
	"github.com/rivo/tview"
)

type Component interface {
	GetPrimitive() tview.Primitive
	CanAddItem() bool
}

type ComponentNode struct {
	Comp     Component
	Parent   *ComponentNode
	Children []*ComponentNode
}
