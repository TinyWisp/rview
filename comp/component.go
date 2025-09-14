package comp

import (
	"github.com/rivo/tview"
)

type Component interface {
	GetName() string
	Primitive() tview.Primitive
	CanAddItem() bool
	SetProp(string, interface{}) error
	GetProp(string) (interface{}, error)
}
