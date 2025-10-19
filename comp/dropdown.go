package comp

import (
	"github.com/rivo/tview"
)

type Dropdown struct {
	Base[*tview.DropDown]
}

func CreateDropdown() Component {
	dropdown := &Dropdown{
		Base: Base[*tview.DropDown]{
			name:      "dropdown",
			tviewInst: tview.NewDropDown(),
		},
	}

	dropdown.Base.outerInst = dropdown

	return dropdown
}
