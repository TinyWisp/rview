package template

import (
	"fmt"
	"math"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
			str: "obj.key",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: ".",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "obj",
				},
				Right: &TplExp{
					Type: TplExpStr,
					Str:  "key",
				},
			},
		},
		{
			str: "obj.key.subkey",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: ".",
				Left: &TplExp{
					Type:     TplExpCalc,
					Operator: ".",
					Left: &TplExp{
						Type:     TplExpVar,
						Variable: "obj",
					},
					Right: &TplExp{
						Type: TplExpStr,
						Str:  "key",
					},
				},
				Right: &TplExp{
					Type: TplExpStr,
					Str:  "subkey",
				},
			},
		},
		{
			str: "var1[attr1]",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "[",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type:     TplExpVar,
					Variable: "attr1",
				},
			},
		},
		{
			str: "var1['attr1']",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "[",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type: TplExpStr,
					Str:  "attr1",
				},
			},
		},
		{
			str: "var1[3]",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "[",
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
			str: "var1[3]['attr1']",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "[",
				Left: &TplExp{
					Type:     TplExpCalc,
					Operator: "[",
					Left: &TplExp{
						Type:     TplExpVar,
						Variable: "var1",
					},
					Right: &TplExp{
						Type: TplExpInt,
						Int:  3,
					},
				},
				Right: &TplExp{
					Type: TplExpStr,
					Str:  "attr1",
				},
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
			str: `func1("param1", param2, 333, 555.5, true, false, nil, var1.attr)`,
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
					{
						Type:     TplExpCalc,
						Operator: ".",
						Left: &TplExp{
							Type:     TplExpVar,
							Variable: "var1",
						},
						Right: &TplExp{
							Type: TplExpStr,
							Str:  "attr",
						},
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
			str: "-var1",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "-",
				Right: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
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
			str: "var1 >=   3",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: ">=",
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
			str: "var1 <3",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "<",
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
			str: "var1<=3",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "<=",
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
			str: "var1 != 3.14159",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "!=",
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
			str: "var1 && !var2",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "&&",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type:     TplExpCalc,
					Operator: "!",
					Left:     nil,
					Right: &TplExp{
						Type:     TplExpVar,
						Variable: "var2",
					},
				},
			},
		},
		{
			str: "var1 || !var2",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "||",
				Left: &TplExp{
					Type:     TplExpVar,
					Variable: "var1",
				},
				Right: &TplExp{
					Type:     TplExpCalc,
					Operator: "!",
					Left:     nil,
					Right: &TplExp{
						Type:     TplExpVar,
						Variable: "var2",
					},
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
		{
			str: "(((a + b)*c)+d)+e",
			exp: TplExp{
				Type:     TplExpCalc,
				Operator: "+",
				Left: &TplExp{
					Type:     TplExpCalc,
					Operator: "+",
					Left: &TplExp{
						Type:     TplExpCalc,
						Operator: "*",
						Left: &TplExp{
							Type:     TplExpCalc,
							Operator: "+",
							Left: &TplExp{
								Type:     TplExpVar,
								Variable: "a",
							},
							Right: &TplExp{
								Type:     TplExpVar,
								Variable: "b",
							},
						},
						Right: &TplExp{
							Type:     TplExpVar,
							Variable: "c",
						},
					},
					Right: &TplExp{
						Type:     TplExpVar,
						Variable: "d",
					},
				},
				Right: &TplExp{
					Type:     TplExpVar,
					Variable: "e",
				},
			},
		},
		{
			str: "{}",
			exp: TplExp{
				Type: TplExpMap,
				Map:  map[string]*TplExp{},
			},
		},
		{
			str: "{a:1}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
				},
			},
		},
		{
			str: "{a:1;}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
				},
			},
		},
		{
			str: "{a:1;\nb:2}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
					"b": {
						Type: TplExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;\nb:2;}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
					"b": {
						Type: TplExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;\n\rb:2;}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
					"b": {
						Type: TplExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;   b:1.1 ;  c:\"str\";d:'str'; e: true;\n f: false;\n\r g:nil; h:var1; i:a+b+3;  j:  func1();  k:func2(a,b);}",
			exp: TplExp{
				Type: TplExpMap,
				Map: map[string]*TplExp{
					"a": {
						Type: TplExpInt,
						Int:  1,
					},
					"b": {
						Type:  TplExpFloat,
						Float: 1.1,
					},
					"c": {
						Type: TplExpStr,
						Str:  "str",
					},
					"d": {
						Type: TplExpStr,
						Str:  "str",
					},
					"e": {
						Type: TplExpBool,
						Bool: true,
					},
					"f": {
						Type: TplExpBool,
						Bool: false,
					},
					"g": {
						Type: TplExpNil,
					},
					"h": {
						Type:     TplExpVar,
						Variable: "var1",
					},
					"i": {
						Type:     TplExpCalc,
						Operator: "+",
						Left: &TplExp{
							Type:     TplExpCalc,
							Operator: "+",
							Left: &TplExp{
								Type:     TplExpVar,
								Variable: "a",
							},
							Right: &TplExp{
								Type:     TplExpVar,
								Variable: "b",
							},
						},
						Right: &TplExp{
							Type: TplExpInt,
							Int:  3,
						},
					},
					"j": {
						Type:       TplExpFunc,
						FuncName:   "func1",
						FuncParams: []*TplExp{},
					},
					"k": {
						Type:     TplExpFunc,
						FuncName: "func2",
						FuncParams: []*TplExp{
							{
								Type:     TplExpVar,
								Variable: "a",
							},
							{
								Type:     TplExpVar,
								Variable: "b",
							},
						},
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

	case TplExpMap:
		if len(a.Map) != len(b.Map) {
			return false
		}
		for k, v := range a.Map {
			if b.Map[k] == nil {
				return false
			}
			if !isTplExpEqual(*v, *b.Map[k]) {
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
		fmt.Printf("expression: %s\n", str)
		realExp, err = ParseTplExp(str)
		if err != nil {
			t.Fatalf("error: %s", err)
		} else if !isTplExpEqual(*realExp, testCase.exp) {
			spew.Dump(*realExp)
			spew.Dump(testCase.exp)
			t.Fatalf("the expression is not parsed as expected\n")
		}
	}
}
