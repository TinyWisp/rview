package comp

import (
	"github.com/rivo/tview"
)

type List struct {
	Base[*tview.List]
}

func CreateList() Component {
	list := &List{
		Base: Base[*tview.List]{
			name:      "list",
			tviewInst: tview.NewList(),
		},
	}

	list.Base.outerInst = list

	return list
}
