package comp

import (
	"github.com/rivo/tview"
)

type TreeView struct {
	Base[*tview.TreeView]
}

func CreateTreeView() Component {
	treeview := &TreeView{
		Base: Base[*tview.TreeView]{
			name:      "treeview",
			tviewInst: tview.NewTreeView(),
		},
	}

	treeview.Base.outerInst = treeview

	return treeview
}
