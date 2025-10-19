package comp

import (
	"github.com/rivo/tview"
)

type Button struct {
	Base[*tview.Button]
}

func CreateButton() Component {
	button := &Button{
		Base: Base[*tview.Button]{
			name:      "button",
			tviewInst: tview.NewButton(""),
		},
	}
	button.Base.outerInst = button

	return button
}
