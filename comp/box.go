package comp

import (
	"github.com/rivo/tview"
)

/*
type Box struct {
	inst *tview.Box
}

func (b *Box) GetPrimitive() tview.Primitive {
	return b.inst
}

func (b *Box) CanAddItem() bool {
	return false
}
*/

type Box struct {
	Base[*tview.Box]
}

func CreateBox() Component {
	return &Box{
		Base[*tview.Box]{
			inst: tview.NewBox(),
		},
	}
}
