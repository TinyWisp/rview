package ddl

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type parseDdlTestCase struct {
	str string
	def DDLDef
	err string
}

var (
	parseDdlTestCases = []parseDdlTestCase{
		{
			str: `
						<template>
							<div></div>
						</template>`,
			def: DDLDef{
				TplMap: map[string]*TplNode{
					"main": {
						Type:    TplNodeTag,
						TagName: "template",
						Children: []*TplNode{
							{
								Type:    TplNodeTag,
								TagName: "div",
							},
						},
					},
				},
			},
		},
		{
			str: `
						<template>
							<div></div>
						</template>

						<template def='link-item(param1,param2)'>
							<div></div>
						</template>`,
			def: DDLDef{
				TplMap: map[string]*TplNode{
					"main": {
						Type:    TplNodeTag,
						TagName: "template",
						Children: []*TplNode{
							{
								Type:    TplNodeTag,
								TagName: "div",
							},
						},
					},
					"link-item": {
						Type:    TplNodeTag,
						TagName: "template",
						Def: &Exp{
							Type:     ExpFunc,
							FuncName: "link-item",
							FuncParams: []*Exp{
								{
									Type:     ExpVar,
									Variable: "param1",
								},
								{
									Type:     ExpVar,
									Variable: "param2",
								},
							},
						},
						Children: []*TplNode{
							{
								Type:    TplNodeTag,
								TagName: "div",
							},
						},
					},
				},
			},
		},
		{
			str: `
						<template>
							<div></div>
						</template>

						<template def='link-item(param1,param2)'>
							<div></div>
						</template>
						
						<style>
							.class1 {
								margin-top: 10ch;
							}
						</style>
						`,
			def: DDLDef{
				TplMap: map[string]*TplNode{
					"main": {
						Type:    TplNodeTag,
						TagName: "template",
						Children: []*TplNode{
							{
								Type:    TplNodeTag,
								TagName: "div",
							},
						},
					},
					"link-item": {
						Type:    TplNodeTag,
						TagName: "template",
						Def: &Exp{
							Type:     ExpFunc,
							FuncName: "link-item",
							FuncParams: []*Exp{
								{
									Type:     ExpVar,
									Variable: "param1",
								},
								{
									Type:     ExpVar,
									Variable: "param2",
								},
							},
						},

						Children: []*TplNode{
							{
								Type:    TplNodeTag,
								TagName: "div",
							},
						},
					},
				},
				CssClassMap: CSSClassMap{
					"class1": {
						"margin-top": []CSSToken{
							{
								Type: CSSTokenNum,
								Num:  10,
								Unit: ch,
							},
						},
					},
				},
			},
		},
	}
)

func isDdlDefEqual(a DDLDef, b DDLDef) bool {
	if len(a.TplMap) != len(b.TplMap) {
		return false
	}

	for tname, tna := range a.TplMap {
		tnb, ok := b.TplMap[tname]
		if !ok {
			return false
		}
		if !isTplEqual([]*TplNode{tna}, []*TplNode{tnb}) {
			return false
		}
	}

	if len(a.CssClassMap) != len(b.CssClassMap) {
		return false
	}

	if !isCssClassMapEqual(a.CssClassMap, b.CssClassMap) {
		return false
	}

	return true
}

func TestParseDdl(t *testing.T) {
	var realDef DDLDef
	var err error
	var str string
	for _, testCase := range parseDdlTestCases {
		str = testCase.str
		fmt.Printf("ddl: %s\n", str)
		realDef, err = ParseDdl(str)
		if err != nil {
			if dpe, ok := err.(*DdlParseError); ok {
				if testCase.err != "" && dpe.err == testCase.err {
					continue
				}
			}
			t.Fatalf("error: %s", err)
		} else if !isDdlDefEqual(realDef, testCase.def) {
			spew.Dump(realDef)
			spew.Dump(testCase.def)
			t.Fatalf("the ddl is not parsed as expected\n")
		}
	}
}
