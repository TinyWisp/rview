package comp

import (
	"github.com/rivo/tview"
)

type Table struct {
	Base[*tview.Table]
}

func CreateTable() Component {
	table := &Table{
		Base: Base[*tview.Table]{
			name:      "table",
			tviewInst: tview.NewTable(),
		},
	}

	table.Base.outerInst = table

	return table
}
