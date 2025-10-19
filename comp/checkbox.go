package comp

import (
	"github.com/rivo/tview"
)

type Checkbox struct {
	Base[*tview.Checkbox]
}

func CreateCheckbox() Component {
	checkbox := &Checkbox{
		Base: Base[*tview.Checkbox]{
			name:      "checkbox",
			tviewInst: tview.NewCheckbox(),
		},
	}

	checkbox.Base.outerInst = checkbox

	return checkbox
}
