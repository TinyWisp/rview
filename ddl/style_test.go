package ddl

import (
	"fmt"
	"math"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type tokenizeCssTestCase struct {
	str    string
	tokens []CSSToken
}

type parseCssTestCase struct {
	str      string
	classMap CSSClassMap
	err      string
}

var (
	tokenizeCssTestCases = []tokenizeCssTestCase{
		{
			str: "+",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "+",
				},
			},
		},
		{
			str: "-",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "-",
				},
			},
		},
		{
			str: "*",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "*",
				},
			},
		},
		{
			str: "/",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "/",
				},
			},
		},
		{
			str: "{",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "{",
				},
			},
		},
		{
			str: "}",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "}",
				},
			},
		},
		{
			str: "(",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: "(",
				},
			},
		},
		{
			str: ")",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: ")",
				},
			},
		},
		{
			str: ":",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: ":",
				},
			},
		},
		{
			str: ";",
			tokens: []CSSToken{
				{
					Type:     CSSTokenOperator,
					Operator: ";",
				},
			},
		},
		{
			str: "10",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  10,
					Unit: noUnit,
				},
			},
		},
		{
			str: "5ch",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: ch,
				},
			},
		},
		{
			str: "5ch",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: ch,
				},
			},
		},
		{
			str: "5vw",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: vw,
				},
			},
		},
		{
			str: "5vh",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: vh,
				},
			},
		},
		{
			str: "5pfw",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: pfw,
				},
			},
		},
		{
			str: "5pfh",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: pfh,
				},
			},
		},
		{
			str: "5pcw",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: pcw,
				},
			},
		},
		{
			str: "5pch",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  5,
					Unit: pch,
				},
			},
		},
		{
			str: "#ababab",
			tokens: []CSSToken{
				{
					Type:  CSSTokenColor,
					Color: "#ababab",
				},
			},
		},
		{
			str: "left",
			tokens: []CSSToken{
				{
					Type: CSSTokenStr,
					Str:  "left",
				},
			},
		},
		{
			str: "\"hello\"",
			tokens: []CSSToken{
				{
					Type: CSSTokenStr,
					Str:  "hello",
				},
			},
		},
		{
			str: "'hello'",
			tokens: []CSSToken{
				{
					Type: CSSTokenStr,
					Str:  "hello",
				},
			},
		},
		{
			str: ".abc",
			tokens: []CSSToken{
				{
					Type:  CSSTokenClass,
					Class: "abc",
				},
			},
		},
		{
			str: "calc(70vw - 30ch)",
			tokens: []CSSToken{
				{
					Type:     CSSTokenFunc,
					FuncName: "calc",
					Arguments: [][]CSSToken{
						{
							{
								Type: CSSTokenNum,
								Num:  70,
								Unit: vw,
							},
							{
								Type:     CSSTokenOperator,
								Operator: "-",
							},
							{
								Type: CSSTokenNum,
								Num:  30,
								Unit: ch,
							},
						},
					},
				},
			},
		},
		{
			str: "10ch #ffffff",
			tokens: []CSSToken{
				{
					Type: CSSTokenNum,
					Num:  10,
					Unit: ch,
				},
				{
					Type:  CSSTokenColor,
					Color: "#ffffff",
				},
			},
		},
		{
			str: ".class1 {\nmargin-left: 10ch;\n}",
			tokens: []CSSToken{
				{
					Type:  CSSTokenClass,
					Class: "class1",
				},
				{
					Type:     CSSTokenOperator,
					Operator: "{",
				},
				{
					Type: CSSTokenProp,
					Prop: "margin-left",
				},
				{
					Type:     CSSTokenOperator,
					Operator: ":",
				},
				{
					Type: CSSTokenNum,
					Num:  10,
					Unit: ch,
				},
				{
					Type:  CSSTokenOperator,
					Color: ";",
				},
				{
					Type:  CSSTokenOperator,
					Color: "}",
				},
			},
		},
		{
			str: ".class1 {\nmargin-left: 10ch;\nmargin-right: 11ch;\n}",
			tokens: []CSSToken{
				{
					Type:  CSSTokenClass,
					Class: "class1",
				},
				{
					Type:     CSSTokenOperator,
					Operator: "{",
				},
				{
					Type: CSSTokenProp,
					Prop: "margin-left",
				},
				{
					Type:     CSSTokenOperator,
					Operator: ":",
				},
				{
					Type: CSSTokenNum,
					Num:  10,
					Unit: ch,
				},
				{
					Type:  CSSTokenOperator,
					Color: ";",
				},
				{
					Type: CSSTokenProp,
					Prop: "margin-right",
				},
				{
					Type:     CSSTokenOperator,
					Operator: ":",
				},
				{
					Type: CSSTokenNum,
					Num:  11,
					Unit: ch,
				},
				{
					Type:  CSSTokenOperator,
					Color: ";",
				},
				{
					Type:  CSSTokenOperator,
					Color: "}",
				},
			},
		},
		{
			str: ".class1 .class2 {\nprop1: 10ch 0 #ffffff;\n}",
			tokens: []CSSToken{
				{
					Type:  CSSTokenClass,
					Class: "class1",
				},
				{
					Type:  CSSTokenClass,
					Class: "class2",
				},
				{
					Type:     CSSTokenOperator,
					Operator: "{",
				},
				{
					Type: CSSTokenProp,
					Prop: "prop1",
				},
				{
					Type:     CSSTokenOperator,
					Operator: ":",
				},
				{
					Type: CSSTokenNum,
					Num:  10,
					Unit: ch,
				},
				{
					Type: CSSTokenNum,
					Num:  0,
					Unit: noUnit,
				},
				{
					Type:  CSSTokenColor,
					Color: "#ffffff",
				},
				{
					Type:  CSSTokenOperator,
					Color: ";",
				},
				{
					Type:  CSSTokenOperator,
					Color: "}",
				},
			},
		},
	}

	parseCssTestCases = []parseCssTestCase{
		// -------------- errors ---------------
		{
			str: ".class1",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 {",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 {margin",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 {margin:",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 {margin: 1ch",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 {margin: 1ch;",
			err: "css.unexpectedEnd",
		},
		{
			str: ".class1 .class2",
			err: "css.unexpectedToken",
		},
		{
			str: ".class1 abc",
			err: "css.unexpectedToken",
		},
		{
			str: ".class1 {margin margin",
			err: "css.unexpectedToken",
		},
		{
			str: ".class1 {margin: }",
			err: "css.unexpectedToken",
		},
		{
			str: ".class1 {abc:15ch;}",
			err: "css.unsupportedProp",
		},
		{
			str: ".class1 {abc:15ch;}",
			err: "css.unsupportedProp",
		},
		{
			str: ".class1 {margin: #ff0000;}",
			err: "css.invalidPropVal",
		},
		{
			str: ".class1 {width: calc(3ch+",
			err: "css.mismatchedParenthesis",
		},

		// -------------- margin ---------------
		{
			str: ".class1 {\nmargin: 1ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"margin": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nmargin: 1ch 2ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"margin": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"margin-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"margin-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"margin-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nmargin: 1ch 2ch 3ch 4ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"margin": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"margin-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"margin-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"margin-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"margin-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nmargin: 1ch 2ch 3ch 4ch;\nmargin-left: 5ch;\nmargin-right: 6ch;\nmargin-top: 7ch;\nmargin-bottom: 8ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"margin": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"margin-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  5,
							Unit: ch,
						},
					},
					"margin-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  6,
							Unit: ch,
						},
					},
					"margin-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  7,
							Unit: ch,
						},
					},
					"margin-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  8,
							Unit: ch,
						},
					},
				},
			},
		},
		// -------------- padding ---------------
		{
			str: ".class1 {\npadding: 1ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"padding": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\npadding: 1ch 2ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"padding": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"padding-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"padding-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"padding-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\npadding: 1ch 2ch 3ch 4ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"padding": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"padding-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"padding-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"padding-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"padding-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\npadding: 1ch 2ch 3ch 4ch;\npadding-left: 5ch;\npadding-right: 6ch;\npadding-top: 7ch;\npadding-bottom: 8ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"padding": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"padding-left": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  5,
							Unit: ch,
						},
					},
					"padding-right": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  6,
							Unit: ch,
						},
					},
					"padding-top": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  7,
							Unit: ch,
						},
					},
					"padding-bottom": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  8,
							Unit: ch,
						},
					},
				},
			},
		},

		// -------------- border-width ---------------
		{
			str: ".class1 {\nborder-width: 1ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-left-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-right-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-top-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-bottom-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-width: 1ch 2ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"border-left-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"border-right-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"border-top-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-bottom-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-width: 1ch 2ch 3ch 4ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"border-left-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"border-right-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
					},
					"border-top-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
					},
					"border-bottom-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-width: 1ch 2ch 3ch 4ch;\nborder-left-width: 5ch;\nborder-right-width: 6ch;\nborder-top-width: 7ch;\nborder-bottom-width: 8ch;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  1,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  2,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  3,
							Unit: ch,
						},
						{
							Type: CSSTokenNum,
							Num:  4,
							Unit: ch,
						},
					},
					"border-left-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  5,
							Unit: ch,
						},
					},
					"border-right-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  6,
							Unit: ch,
						},
					},
					"border-top-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  7,
							Unit: ch,
						},
					},
					"border-bottom-width": []CSSToken{
						{
							Type: CSSTokenNum,
							Num:  8,
							Unit: ch,
						},
					},
				},
			},
		},

		// -------------- border-color ---------------
		{
			str: ".class1 {\nborder-color: #ffffff;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
					"border-left-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
					"border-right-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
					"border-top-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
					"border-bottom-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-color: #ab0000 #00ab00;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ab0000",
						},
						{
							Type:  CSSTokenColor,
							Color: "#00ab00",
						},
					},
					"border-left-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#00ab00",
						},
					},
					"border-right-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#00ab00",
						},
					},
					"border-top-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ab0000",
						},
					},
					"border-bottom-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ab0000",
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-color: #aaaaaa #bbbbbb #cccccc #dddddd;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#aaaaaa",
						},
						{
							Type:  CSSTokenColor,
							Color: "#bbbbbb",
						},
						{
							Type:  CSSTokenColor,
							Color: "#cccccc",
						},
						{
							Type:  CSSTokenColor,
							Color: "#dddddd",
						},
					},
					"border-left-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#dddddd",
						},
					},
					"border-right-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#bbbbbb",
						},
					},
					"border-top-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#aaaaaa",
						},
					},
					"border-bottom-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#cccccc",
						},
					},
				},
			},
		},
		{
			str: ".class1 {\nborder-color: #aaaaaa #bbbbbb #cccccc #dddddd;\nborder-left-color: #aa0000;\nborder-right-color: #bb0000;\nborder-top-color: #cc0000;\nborder-bottom-color: #dd0000;\n}",
			classMap: CSSClassMap{
				"class1": {
					"border-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#aaaaaa",
						},
						{
							Type:  CSSTokenColor,
							Color: "#bbbbbb",
						},
						{
							Type:  CSSTokenColor,
							Color: "#cccccc",
						},
						{
							Type:  CSSTokenColor,
							Color: "#dddddd",
						},
					},
					"border-left-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#aa0000",
						},
					},
					"border-right-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#bb0000",
						},
					},
					"border-top-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#cc0000",
						},
					},
					"border-bottom-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#dd0000",
						},
					},
				},
			},
		},
		// ----------------- bordr-char -------------------
		{
			str: ".class1 {\nborder-char: ╔ ╗ ╝ ╚ ═ ║;\n}",
			classMap: CSSClassMap{
				"class1": CSSClass{
					"border-char": []CSSToken{
						{
							Type: CSSTokenStr,
							Str:  "╔",
						},
						{
							Type: CSSTokenStr,
							Str:  "╗",
						},
						{
							Type: CSSTokenStr,
							Str:  "╝",
						},
						{
							Type: CSSTokenStr,
							Str:  "╚",
						},
						{
							Type: CSSTokenStr,
							Str:  "═",
						},
						{
							Type: CSSTokenStr,
							Str:  "║",
						},
					},
				},
			},
		},
		// ----------------- background-color -------------------
		{
			str: ".class1 {\nbackground-color: #ff0000;\n}",
			classMap: CSSClassMap{
				"class1": CSSClass{
					"background-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ffffff",
						},
					},
				},
			},
		},
		// ----------------- multiple classes -------------------
		{
			str: ".class1 {\nbackground-color: #ff0000;\n}\n.class2 {\nbackground-color: #00ff00;\n}",
			classMap: CSSClassMap{
				"class1": CSSClass{
					"background-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#ff0000",
						},
					},
				},
				"class2": CSSClass{
					"background-color": []CSSToken{
						{
							Type:  CSSTokenColor,
							Color: "#00ff00",
						},
					},
				},
			},
		},
	}
)

func isCssTokenEqual(a CSSToken, b CSSToken) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case CSSTokenNum:
		if a.Unit != b.Unit || math.Abs(a.Num-b.Num) > 1e-9 {
			return false
		}

	case CSSTokenStr:
		if a.Str != b.Str {
			return false
		}

	case CSSTokenFunc:
		if a.FuncName != b.FuncName {
			return false
		}
		if len(a.Arguments) != len(b.Arguments) {
			return false
		}
		for aidx, arga := range a.Arguments {
			argb := b.Arguments[aidx]
			if !areCssTokensEqual(arga, argb) {
				return false
			}
		}
	}

	return true
}

func areCssTokensEqual(a []CSSToken, b []CSSToken) bool {
	if len(a) != len(b) {
		return false
	}

	for tidx, ta := range a {
		tb := b[tidx]
		if !isCssTokenEqual(ta, tb) {
			return false
		}
	}

	return true
}

func isCssClassEqual(a CSSClass, b CSSClass) bool {
	if len(a) != len(b) {
		return false
	}

	for aprop, atokens := range a {
		if _, ok := b[aprop]; !ok {
			return false
		}

		btokens := b[aprop]
		if len(atokens) != len(btokens) {
			return false
		}

		for idx, atoken := range atokens {
			btoken := btokens[idx]
			if !isCssTokenEqual(atoken, btoken) {
				return false
			}
		}
	}

	return true
}

func isCssClassMapEqual(a CSSClassMap, b CSSClassMap) bool {
	if len(a) != len(b) {
		return false
	}

	for acname, aclass := range a {
		if _, ok := b[acname]; !ok {
			return false
		}
		bclass := b[acname]

		if !isCssClassEqual(aclass, bclass) {
			return false
		}
	}

	return true
}

func TestTokenizeCss(t *testing.T) {
	var realTokens []CSSToken
	var err error
	var str string
	for _, testCase := range tokenizeCssTestCases {
		str = testCase.str
		fmt.Printf("css: %s\n", str)
		realTokens, err = tokenizeCss(str)
		if err != nil {
			t.Fatalf("error: %s", err)
		} else if !areCssTokensEqual(realTokens, testCase.tokens) {
			spew.Dump(realTokens)
			spew.Dump(testCase.tokens)
			t.Fatalf("the css string is not as expected\n")
		}
	}
}

func TestParseCss(t *testing.T) {
	var realClassMap CSSClassMap
	var err error
	var str string
	for _, testCase := range parseCssTestCases {
		str = testCase.str
		fmt.Printf("css: %s\n", str)
		realClassMap, err = parseCss(str)
		if err != nil {
			if tpe, ok := err.(*DdlParseError); ok {
				if testCase.err != "" && tpe.err == testCase.err {
					continue
				}
			}
			t.Fatalf("error: %s", err)
		} else if !isCssClassMapEqual(realClassMap, testCase.classMap) {
			spew.Dump(realClassMap)
			spew.Dump(testCase.classMap)
			t.Fatalf("the css is not parsed as expected\n")
		}
	}

}
