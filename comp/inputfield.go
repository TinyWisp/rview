package comp

import (
	"github.com/rivo/tview"
)

type InputField struct {
	Base[*tview.InputField]
}

func CreateInputField() Component {
	inputField := &InputField{
		Base: Base[*tview.InputField]{
			name:      "inputfield",
			tviewInst: tview.NewInputField(),
		},
	}

	inputField.Base.outerInst = inputField

	return inputField
}
