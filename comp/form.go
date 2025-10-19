package comp

import (
	"github.com/rivo/tview"
)

type Form struct {
	Base[*tview.Form]
}

func CreateForm() Component {
	form := &Form{
		Base: Base[*tview.Form]{
			name:      "form",
			tviewInst: tview.NewForm(),
		},
	}

	form.Base.outerInst = form

	return form
}
