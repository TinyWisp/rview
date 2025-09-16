package comp

import (
	"github.com/rivo/tview"
)

type Box struct {
	Base[*tview.Box]
}

func CreateBox() Component {
	return &Box{
		Base[*tview.Box]{
			name: "box",
			inst: tview.NewBox(),
		},
	}
}
