package template

import (
	"fmt"
	"math"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type readCssTestCase struct {
	str    string
	tokens []CSSToken
}

var (
	readCssTestCases = []readCssTestCase{
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

func TestTokenizeCss(t *testing.T) {
	var realTokens []CSSToken
	var err error
	var str string
	for _, testCase := range readCssTestCases {
		str = testCase.str
		fmt.Printf("css rule: %s\n", str)
		realTokens, err = tokenizeCss(str)
		if err != nil {
			t.Fatalf("error: %s", err)
		} else if !areCssTokensEqual(realTokens, testCase.tokens) {
			spew.Dump(realTokens)
			spew.Dump(testCase.tokens)
			t.Fatalf("the css rule is not readd as expected\n")
		}
	}

}
