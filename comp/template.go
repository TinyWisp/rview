package comp

import (
	"github.com/rivo/tview"
)

type Template struct {
}

func (t *Template) GetName() string {
	return "template"
}

func (t *Template) Primitive() tview.Primitive {
	return nil
}

func (t *Template) SetProp(prop string, val interface{}) error {
	return nil
}

func (t *Template) GetProp(prop string) (interface{}, error) {
	return nil, nil
}

func (t *Template) CanAddItem() bool {
	return true
}

func (t *Template) AddItem(comp *Component, props map[string]interface{}) error {
	return nil
}

func CreateTemplate() Component {
	return &Template{}
}
