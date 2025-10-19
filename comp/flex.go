package comp

import (
	"github.com/rivo/tview"
)

type Flex struct {
	Base[*tview.Flex]
}

func CreateFlex() Component {
	flex := &Flex{
		Base: Base[*tview.Flex]{
			name:      "flex",
			tviewInst: tview.NewFlex(),
		},
	}

	flex.Base.outerInst = flex

	return flex
}
