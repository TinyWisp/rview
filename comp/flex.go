package comp

import (
	"github.com/rivo/tview"
)

type Flex struct {
	Base[*tview.Flex]
}

func (b *Flex) CanAddItem() bool {
	return true
}

func (b *Flex) AddItem(comp *Component, props map[string]interface{}) error {
	fixedSize, err := getFromPropMap(props, "fixed-size", false, 1)
	if err != nil {
		return err
	}

	proportion, err := getFromPropMap(props, "proportion", false, 1)
	if err != nil {
		return err
	}

	focus, err := getFromPropMap(props, "focus", false, false)
	if err != nil {
		return err
	}

	b.inst.AddItem((*comp).Primitive(), fixedSize.(int), proportion.(int), focus.(bool))
	return nil
}

func CreateFlex() Component {
	return &Flex{
		Base[*tview.Flex]{
			inst: tview.NewFlex(),
		},
	}
}
