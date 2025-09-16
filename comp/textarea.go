package comp

import (
	"github.com/rivo/tview"
)

type TextArea struct {
	Base[*tview.TextArea]
}

func CreateTextArea() Component {
	return &TextArea{
		Base[*tview.TextArea]{
			name: "textarea",
			inst: tview.NewTextArea(),
		},
	}
}
