package comp

import (
	"github.com/rivo/tview"
)

type Box struct {
	inst *tview.Box
}

func (b *Box) GetPrimitive() tview.Primitive {
	return b.inst
}

func (b *Box) CanAddItem() bool {
	return false
}

func CreateBox() Component {
	return &Box{
		inst: tview.NewBox(),
	}
}
