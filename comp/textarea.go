package comp

import (
	"github.com/rivo/tview"
)

type Textarea struct {
	Base[*tview.TextArea]
}

func (t *Textarea) SetText(text string) {
	t.tviewInst.SetText(text, false)
}

func CreateTextArea() Component {
	textarea := &Textarea{
		Base: Base[*tview.TextArea]{
			name:      "textarea",
			tviewInst: tview.NewTextArea(),
		},
	}

	textarea.Base.outerInst = textarea

	return textarea
}
