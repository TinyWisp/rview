package comp

import (
	"github.com/rivo/tview"
)

type Flex struct {
	inst *tview.Flex
}

func (b *Flex) GetPrimitive() tview.Primitive {
	return b.inst
}

func (b *Flex) CanAddItem() bool {
	return true
}

func (b *Flex) AddItem(comp *Component) {
	b.inst.AddItem((*comp).GetPrimitive(), 1, 1, false)
}

func CreateFlex() Component {
	return &Flex{
		inst: tview.NewFlex(),
	}
}
