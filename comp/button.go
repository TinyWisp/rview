package comp

import (
	"github.com/rivo/tview"
)

type Button struct {
	Base[*tview.Button]
}

func CreateButton() Component {
	return &Button{
		Base[*tview.Button]{
			name: "button",
			inst: tview.NewButton(""),
		},
	}
}
