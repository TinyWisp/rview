package main

import (
	"fmt"

	"github.com/TinyWisp/rview/comp"
	"github.com/davecgh/go-spew/spew"
	"github.com/rivo/tview"
)

func main() {
	page := comp.Page{
		Tpl: `
			<template>
				<box :border="true" title="hello world"></box>
			</template>
		`,
		TagCompCreatorMap: map[string]func() comp.Component{
			"box":  comp.CreateBox,
			"flex": comp.CreateFlex,
		},
	}
	err1 := page.Init()
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	prim := page.Primitive()
	err2 := tview.NewApplication().SetRoot(prim, true).Run()
	if err2 != nil {
		fmt.Println(err2)
	}
	spew.Print(prim)
}
