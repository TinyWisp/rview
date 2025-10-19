package comp

import (
	"github.com/rivo/tview"
)

type Grid struct {
	Base[*tview.Grid]
}

func CreateGrid() Component {
	grid := &Grid{
		Base: Base[*tview.Grid]{
			name:      "grid",
			tviewInst: tview.NewGrid(),
		},
	}

	grid.Base.outerInst = grid

	return grid
}
