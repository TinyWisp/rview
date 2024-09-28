package template

import (
	"fmt"
	"math"
	"testing"
)

type parseTplExpTestCase struct {
	str string
	exp TplExp
}

var (
	parseTplExpTestCases = []parseTplExpTestCase{
		{
			str: "true",
			exp: TplExp{
				Type: TplExpBool,
				Bool: true,
			},
		},
		{
			str: "false",
			exp: TplExp{
				Type: TplExpBool,
				Bool: false,
			},
		},
		{
			str: "nil",
			exp: TplExp{
				Type: TplExpNil,
			},
		},
		{
			str: "0",
			exp: TplExp{
				Type: TplExpInt,
				Int:  0,
			},
		},
		{
			str: "333",
			exp: TplExp{
				Type: TplExpInt,
				Int:  333,
			},
		},
		{
			str: "333.333",
			exp: TplExp{
				Type:  TplExpFloat,
				Float: 333.333,
			},
		},
		{
			str: `""`,
			exp: TplExp{
				Type: TplExpStr,
				Str:  "",
			},
		},
		{
			str: `''`,
			exp: TplExp{
				Type: TplExpStr,
				Str:  "",
			},
		},
		{
			str: `"hello, \"world\""`,
			exp: TplExp{
				Type: TplExpStr,
				Str:  `hello, "world"`,
			},
		},
		{
			str: `'hello, \'world\''`,
			exp: TplExp{
				Type: TplExpStr,
				Str:  `hello, 'world'`,
			},
		},
		{
			str: "v",
			exp: TplExp{
				Type:     TplExpVar,
				Variable: "v",
			},
		},
		{
			str: "var1",
			exp: TplExp{
				Type:     TplExpVar,
				Variable: "var1",
			},
		},
		{
			str: "func1()",
			exp: TplExp{
				Type:       TplExpFunc,
				FuncName:   "func1",
				FuncParams: make([]*TplExp, 0),
			},
		},
		{
			str: "func1()",
			exp: TplExp{
				Type:       TplExpFunc,
				FuncName:   "func1",
				FuncParams: make([]*TplExp, 0),
			},
		},
		{
			str: `func1("param1", param2, 333, 555.5, true, false, nil)`,
			exp: TplExp{
				Type:     TplExpFunc,
				FuncName: "func1",
				FuncParams: []*TplExp{
					{
						Type: TplExpStr,
						Str:  "param1",
					},
					{
						Type:     TplExpVar,
						Variable: "param2",
					},
					{
						Type: TplExpInt,
						Int:  333,
					},
					{
						Type:  TplExpFloat,
						Float: 555.5,
					},
					{
						Type: TplExpBool,
						Bool: true,
					},
					{
						Type: TplExpBool,
						Bool: false,
					},
					{
						Type: TplExpNil,
					},
				},
			},
		},
		{
			str: "-1",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "-",
				Right: &TplExp{
					Type: TplExpInt,
					Int:  1,
				},
			},
		},
		{
			str: "-func1()",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "-",
				Right: &TplExp{
					Type:       TplExpFunc,
					FuncName:   "func1",
					FuncParams: make([]*TplExp, 0),
				},
			},
		},
		{
			str: "!var1",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "!",
				Right: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
			},
		},
		{
			str: "var1 >3",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: ">",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type: TplExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1 == 3.14159",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "==",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type:  TplExpFloat,
					Float: 3.14159,
				},
			},
		},
		{
			str: "var1 == 3.14159 && var2 >= var1",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "&&",
				Left: &TplExp{
					Type:     TplExpCalc,
					Operator: "==",
					Left: &TplExp{
						Type:     TplExpVar,
						Variable: "var1",
					},
					Right: &TplExp{
						Type:  TplExpFloat,
						Float: 3.14159,
					},
				},
				Right: &TplExp{
					Type:     TplExpCalc,
					Operator: ">=",
					Left: &TplExp{
						Type:     TplExpVar,
						Variable: "var2",
					},
					Right: &TplExp{
						Type:     TplExpVar,
						Variable: "var1",
					},
				},
			},
		},
		{
			str: "var1+ 111 + (var2 -var3)*5",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "+",
				Left: &TplExp{
					Type:     TplExpCalc,
					Operator: "+",
					Left: &TplExp{
						Type:     TplExpVar,
						Variable: "var1",
					},
					Right: &TplExp{
						Type: TplExpInt,
						Int:  111,
					},
				},
				Right: &TplExp{
					Type:     TplExpCalc,
					Operator: "*",
					Left: &TplExp{
						Type:     TplExpCalc,
						Operator: "-",
						Left: &TplExp{
							Type:     TplExpVar,
							Variable: "var2",
						},
						Right: &TplExp{
							Type:     TplExpVar,
							Variable: "var3",
						},
					},
					Right: &TplExp{
						Type: TplExpInt,
						Int:  5,
					},
				},
			},
		},
	}
)

func isTplExpEqual(a TplExp, b TplExp) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case TplExpInt:
		if a.Int != b.Int {
			return false
		}

	case TplExpFloat:
		if math.Abs(a.Float-b.Float) > 1e-9 {
			return false
		}

	case TplExpBool:
		if a.Bool != b.Bool {
			return false
		}

	case TplExpStr:
		if a.Str != b.Str {
			return false
		}

	case TplExpVar:
		if a.Variable != b.Variable {
			return false
		}

	case TplExpOperator:
		if a.Operator != b.Operator {
			return false
		}

	case TplExpFunc:
		if a.FuncName != b.FuncName || len(a.FuncParams) != len(b.FuncParams) {
			return false
		}
		for i := 0; i < len(a.FuncParams); i++ {
			if !isTplExpEqual(*a.FuncParams[i], *b.FuncParams[i]) {
				return false
			}
		}

	case TplExpCalc:
		if a.Operator != b.Operator ||
			(a.Left == nil && b.Left != nil) ||
			(a.Left != nil && b.Left == nil) ||
			(a.Right == nil && b.Right != nil) ||
			(a.Right != nil && b.Right == nil) ||
			(a.Left != nil && b.Left != nil && !isTplExpEqual(*a.Left, *b.Left)) ||
			(a.Right != nil && b.Right != nil && !isTplExpEqual(*a.Right, *b.Right)) {
			return false
		}
	}

	return true
}

func TestParseTplExp(t *testing.T) {
	var realExp *TplExp
	var err error
	var str string
	for _, testCase := range parseTplExpTestCases {
		str = testCase.str
		fmt.Printf("template:\n%s\n", str)
		realExp, err = ParseTplExp(str)
		if err != nil {
			t.Fatalf("error: %s", err)
		} else if !isTplExpEqual(*realExp, testCase.exp) {
			t.Fatalf("the template is not parsed as expected\n%+v\n%+v", *realExp, testCase.exp)
		}
	}
}
