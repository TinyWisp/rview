package comp

import (
	"github.com/rivo/tview"
)

type Box struct {
	Base[*tview.Box]
}

func CreateBox() Component {
	box := &Box{
		Base: Base[*tview.Box]{
			name:      "box",
			tviewInst: tview.NewBox(),
		},
	}

	box.Base.outerInst = box

	return box
}
