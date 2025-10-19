package comp

import (
	"github.com/rivo/tview"
)

type Modal struct {
	Base[*tview.Modal]
}

func CreateModal() Component {
	modal := &Modal{
		Base: Base[*tview.Modal]{
			name:      "modal",
			tviewInst: tview.NewModal(),
		},
	}

	modal.Base.outerInst = modal

	return modal
}
