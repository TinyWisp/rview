package template

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type parseTplTestCase struct {
	str string
	tpl []*TplNode
}

var (
	parseTplTestCases = []parseTplTestCase{
		{
			str: `<template></template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
				},
			},
		},
		{
			str: `<template>    hello, world!  </template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type: TplNodeText,
							Text: "hello, world!",
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				hello, world! {{ var1 +var2 }}
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type: TplNodeText,
							Text: "hello, world!",
						},
						{
							Type: TplNodeExp,
							Exp: &TplExp{
								Type:     TplExpCalc,
								Operator: "+",
								Left: &TplExp{
									Type:     TplExpVar,
									Variable: "var1",
								},
								Right: &TplExp{
									Type:     TplExpVar,
									Variable: "var2",
								},
							},
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				<comp-a></comp-a>
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type:    TplNodeTag,
							TagName: "comp-a",
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				<comp-a/>
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type:    TplNodeTag,
							TagName: "comp-a",
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				<comp-a attr1="val1" attr2="val2"/>
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type:    TplNodeTag,
							TagName: "comp-a",
							Attrs: map[string]string{
								"attr1": "val1",
								"attr2": "val2",
							},
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				<comp-a v-if="var1"></comp-a>
				<comp-b v-else-if="var2" />
				<comp-c v-else /> 
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type:    TplNodeTag,
							TagName: "comp-a",
							If: &TplExp{
								Type:     TplExpVar,
								Variable: "var1",
							},
						},
						{
							Type:    TplNodeTag,
							TagName: "comp-b",
							ElseIf: &TplExp{
								Type:     TplExpVar,
								Variable: "var2",
							},
						},
						{
							Type:    TplNodeTag,
							TagName: "comp-c",
							Else: &TplExp{
								Type:     TplExpVar,
								Variable: "",
							},
						},
					},
				},
			},
		},
		{
			str: `
			<template>
				<comp-a v-for="idx, item := range items"></comp-a>
			</template>`,
			tpl: []*TplNode{
				{
					Type:    TplNodeTag,
					TagName: "template",
					Children: []*TplNode{
						{
							Type:    TplNodeTag,
							TagName: "comp-a",
							For: &TplFor{
								Idx:  "idx",
								Item: "item",
								Range: TplExp{
									Type:     TplExpVar,
									Variable: "items",
								},
							},
						},
					},
				},
			},
		},
	}
)

func isTplForEqual(a TplFor, b TplFor) bool {
	return a.Idx == b.Idx && a.Item == b.Item && isTplExpEqual(a.Range, b.Range)
}

func isTplEqual(a []*TplNode, b []*TplNode) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		node1 := *a[idx]
		node2 := *b[idx]

		if node1.Type != node2.Type {
			return false
		}

		if node1.Type == TplNodeTag {
			if node1.TagName != node2.TagName {
				return false
			}

			if (node1.If != nil && node2.If == nil) ||
				(node1.If == nil && node2.If != nil) ||
				(node1.If != nil && node2.If != nil && !isTplExpEqual(*node1.If, *node2.If)) {
				return false
			}

			if (node1.ElseIf != nil && node2.ElseIf == nil) ||
				(node1.ElseIf == nil && node2.ElseIf != nil) ||
				(node1.ElseIf != nil && node2.ElseIf != nil && !isTplExpEqual(*node1.ElseIf, *node2.ElseIf)) {
				return false
			}

			if (node1.Else != nil && node2.Else == nil) ||
				(node1.Else == nil && node2.Else != nil) ||
				(node1.Else != nil && node2.Else != nil && !isTplExpEqual(*node1.Else, *node2.Else)) {
				return false
			}

			if (node1.For != nil && node2.For == nil) ||
				(node1.For == nil && node2.For != nil) ||
				(node1.For != nil && node2.For != nil && !isTplForEqual(*node1.For, *node2.For)) {
				return false
			}

			if (node1.Binds == nil && node2.Binds != nil) ||
				(node1.Binds != nil && node2.Binds == nil) ||
				(node1.Binds != nil && node2.Binds != nil && len(node1.Binds) != len(node2.Binds)) {
				return false
			}

			if node1.Binds != nil {
				for bkey, bval := range node1.Binds {
					if _, ok := node2.Binds[bkey]; !ok {
						return false
					}
					if !isTplExpEqual(*bval, *node2.Binds[bkey]) {
						return false
					}
				}
			}

			if (node1.Events == nil && node2.Binds != nil) ||
				(node1.Events != nil && node2.Events == nil) ||
				(node1.Events == nil && node2.Events == nil && len(node1.Events) != len(node2.Events)) {
				return false
			}

			if node1.Events != nil {
				for ekey, eval := range node1.Events {
					if _, ok := node2.Events[ekey]; !ok {
						return false
					}
					if !isTplExpEqual(*eval, *node2.Events[ekey]) {
						return false
					}
				}
			}

			if (node1.Attrs == nil && node2.Attrs != nil) ||
				(node1.Attrs != nil && node2.Attrs == nil) ||
				(node1.Attrs != nil && node2.Attrs != nil && len(node1.Attrs) != len(node2.Attrs)) {
				return false
			}

			if node1.Attrs != nil {
				for akey, aval := range node1.Attrs {
					if _, ok := node2.Attrs[akey]; !ok {
						return false
					}
					if aval != node2.Attrs[akey] {
						return false
					}
				}
			}

			if !isTplEqual(node1.Children, node2.Children) {
				return false
			}
		}
	}

	return true
}

func TestParseTpl(t *testing.T) {
	var parsedTpl []*TplNode
	var err error
	var str string
	for _, testCase := range parseTplTestCases {
		str = testCase.str
		fmt.Printf("==========================\ntemplate:\n%s\n", str)
		parsedTpl, err = ParseTpl(str)
		if err != nil {
			t.Fatalf("error: %s", err)
		} else if !isTplEqual(parsedTpl, testCase.tpl) {
			spew.Dump(parsedTpl)
			spew.Dump(testCase.tpl)
			t.Fatalf("the template is not parsed as expected\n")
		}
	}
}
