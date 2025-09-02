package ddl

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type parseExpTestCase struct {
	str string
	exp Exp
	err string
}

var (
	parseExpTestCases = []parseExpTestCase{
		{
			str: "true",
			exp: Exp{
				Type: ExpBool,
				Bool: true,
			},
		},
		{
			str: "false",
			exp: Exp{
				Type: ExpBool,
				Bool: false,
			},
		},
		{
			str: "nil",
			exp: Exp{
				Type: ExpNil,
			},
		},
		{

			str: "0",
			exp: Exp{
				Type: ExpInt,
				Int:  0,
			},
		},
		{
			str: "333",
			exp: Exp{
				Type: ExpInt,
				Int:  333,
			},
		},
		{
			str: "333.333",
			exp: Exp{
				Type:  ExpFloat,
				Float: 333.333,
			},
		},
		{
			str: `""`,
			exp: Exp{
				Type: ExpStr,
				Str:  "",
			},
		},
		{
			str: `''`,
			exp: Exp{
				Type: ExpStr,
				Str:  "",
			},
		},
		{
			str: `"hello, \"world\""`,
			exp: Exp{
				Type: ExpStr,
				Str:  `hello, "world"`,
			},
		},
		{
			str: `'hello, \'world\''`,
			exp: Exp{
				Type: ExpStr,
				Str:  `hello, 'world'`,
			},
		},
		{
			str: "v",
			exp: Exp{
				Type:     ExpVar,
				Variable: "v",
			},
		},
		{
			str: "var1",
			exp: Exp{
				Type:     ExpVar,
				Variable: "var1",
			},
		},
		{
			str: "obj.key",
			exp: Exp{
				Type:     ExpCalc,
				Operator: ".",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "obj",
				},
				Right: &Exp{
					Type: ExpStr,
					Str:  "key",
				},
			},
		},
		{
			str: "obj.key.subkey",
			exp: Exp{
				Type:     ExpCalc,
				Operator: ".",
				Left: &Exp{
					Type:     ExpCalc,
					Operator: ".",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "obj",
					},
					Right: &Exp{
						Type: ExpStr,
						Str:  "key",
					},
				},
				Right: &Exp{
					Type: ExpStr,
					Str:  "subkey",
				},
			},
		},
		{
			str: "var1[attr1]",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "[",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type:     ExpVar,
					Variable: "attr1",
				},
			},
		},
		{
			str: "var1['attr1']",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "[",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpStr,
					Str:  "attr1",
				},
			},
		},
		{
			str: "var1[3]",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "[",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1[3]['attr1']",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "[",
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "[",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "var1",
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  3,
					},
				},
				Right: &Exp{
					Type: ExpStr,
					Str:  "attr1",
				},
			},
		},
		{
			str: "func1()",
			exp: Exp{
				Type:       ExpFunc,
				FuncName:   "func1",
				FuncParams: make([]*Exp, 0),
			},
		},
		{
			str: "func1_2()",
			exp: Exp{
				Type:       ExpFunc,
				FuncName:   "func1_2",
				FuncParams: make([]*Exp, 0),
			},
		},
		{
			str: "func1(1)",
			exp: Exp{
				Type:     ExpFunc,
				FuncName: "func1",
				FuncParams: []*Exp{
					{
						Type: ExpInt,
						Int:  1,
					},
				},
			},
		},
		{
			str: `func1("param1", param2, 333, 555.5, true, false, nil, var1.attr)`,
			exp: Exp{
				Type:     ExpFunc,
				FuncName: "func1",
				FuncParams: []*Exp{
					{
						Type: ExpStr,
						Str:  "param1",
					},
					{
						Type:     ExpVar,
						Variable: "param2",
					},
					{
						Type: ExpInt,
						Int:  333,
					},
					{
						Type:  ExpFloat,
						Float: 555.5,
					},
					{
						Type: ExpBool,
						Bool: true,
					},
					{
						Type: ExpBool,
						Bool: false,
					},
					{
						Type: ExpNil,
					},
					{
						Type:     ExpCalc,
						Operator: ".",
						Left: &Exp{
							Type:     ExpVar,
							Variable: "var1",
						},
						Right: &Exp{
							Type: ExpStr,
							Str:  "attr",
						},
					},
				},
			},
		},
		{
			str: `func1(func2())`,
			exp: Exp{
				Type:     ExpFunc,
				FuncName: "func1",
				FuncParams: []*Exp{
					{
						Type:       ExpFunc,
						FuncName:   "func2",
						FuncParams: []*Exp{},
					},
				},
			},
		},
		{
			str: `func1(func2(a, b), func3(c, d))`,
			exp: Exp{
				Type:     ExpFunc,
				FuncName: "func1",
				FuncParams: []*Exp{
					{
						Type:     ExpFunc,
						FuncName: "func2",
						FuncParams: []*Exp{
							{
								Type:     ExpVar,
								Variable: "a",
							},
							{
								Type:     ExpVar,
								Variable: "b",
							},
						},
					},
					{
						Type:     ExpFunc,
						FuncName: "func3",
						FuncParams: []*Exp{
							{
								Type:     ExpVar,
								Variable: "c",
							},
							{
								Type:     ExpVar,
								Variable: "d",
							},
						},
					},
				},
			},
		},
		{
			str: `func1(func2(a, func3(3)), (func4(c, d) + 1) / 2)`,
			exp: Exp{
				Type:     ExpFunc,
				FuncName: "func1",
				FuncParams: []*Exp{
					{
						Type:     ExpFunc,
						FuncName: "func2",
						FuncParams: []*Exp{
							{
								Type:     ExpVar,
								Variable: "a",
							},
							{
								Type:     ExpFunc,
								FuncName: "func3",
								FuncParams: []*Exp{
									{
										Type: ExpInt,
										Int:  3,
									},
								},
							},
						},
					},
					{
						Type:     ExpCalc,
						Operator: "/",
						Left: &Exp{
							Type:     ExpCalc,
							Operator: "+",
							Left: &Exp{
								Type:     ExpFunc,
								FuncName: "func4",
								FuncParams: []*Exp{
									{
										Type:     ExpVar,
										Variable: "c",
									},
									{
										Type:     ExpVar,
										Variable: "d",
									},
								},
							},
							Right: &Exp{
								Type: ExpInt,
								Int:  1,
							},
						},
						Right: &Exp{
							Type: ExpInt,
							Int:  2,
						},
					},
				},
			},
		},
		{
			str: "-1",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "-",
				Right: &Exp{
					Type: ExpInt,
					Int:  1,
				},
			},
		},
		{
			str: "-var1",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "-",
				Right: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
			},
		},
		{
			str: "-func1()",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "-",
				Right: &Exp{
					Type:       ExpFunc,
					FuncName:   "func1",
					FuncParams: make([]*Exp, 0),
				},
			},
		},
		{
			str: "!var1",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "!",
				Right: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
			},
		},
		{
			str: "var1 >3",
			exp: Exp{
				Type:     ExpCalc,
				Operator: ">",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1 >=   3",
			exp: Exp{
				Type:     ExpCalc,
				Operator: ">=",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1 <3",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "<",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1<=3",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "<=",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type: ExpInt,
					Int:  3,
				},
			},
		},
		{
			str: "var1 == 3.14159",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "==",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type:  ExpFloat,
					Float: 3.14159,
				},
			},
		},
		{
			str: "var1 != 3.14159",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "!=",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type:  ExpFloat,
					Float: 3.14159,
				},
			},
		},
		{
			str: "var1 && !var2",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "&&",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "!",
					Left:     nil,
					Right: &Exp{
						Type:     ExpVar,
						Variable: "var2",
					},
				},
			},
		},
		{
			str: "var1 || !var2",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "||",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "!",
					Left:     nil,
					Right: &Exp{
						Type:     ExpVar,
						Variable: "var2",
					},
				},
			},
		},
		{
			str: "var1 == 3.14159 && var2 >= var1",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "&&",
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "==",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "var1",
					},
					Right: &Exp{
						Type:  ExpFloat,
						Float: 3.14159,
					},
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: ">=",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "var2",
					},
					Right: &Exp{
						Type:     ExpVar,
						Variable: "var1",
					},
				},
			},
		},
		{
			str: "var1+ 111 + (var2 -var3)*5",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "+",
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "+",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "var1",
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  111,
					},
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "*",
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "-",
						Left: &Exp{
							Type:     ExpVar,
							Variable: "var2",
						},
						Right: &Exp{
							Type:     ExpVar,
							Variable: "var3",
						},
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  5,
					},
				},
			},
		},
		{
			str: "a + b * 3",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "+",
				Left: &Exp{
					Type:     ExpVar,
					Variable: "a",
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "*",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "b",
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  3,
					},
				},
			},
		},
		{
			str: "(((a + b)*c)+d)+e",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "+",
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "+",
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "*",
						Left: &Exp{
							Type:     ExpCalc,
							Operator: "+",
							Left: &Exp{
								Type:     ExpVar,
								Variable: "a",
							},
							Right: &Exp{
								Type:     ExpVar,
								Variable: "b",
							},
						},
						Right: &Exp{
							Type:     ExpVar,
							Variable: "c",
						},
					},
					Right: &Exp{
						Type:     ExpVar,
						Variable: "d",
					},
				},
				Right: &Exp{
					Type:     ExpVar,
					Variable: "e",
				},
			},
		},
		{
			str: "{}",
			exp: Exp{
				Type: ExpMap,
				Map:  map[string]*Exp{},
			},
		},
		{
			str: "{a:1}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
				},
			},
		},
		{
			str: "{a:1;}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
				},
			},
		},
		{
			str: "{a:1;\nb:2}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
					"b": {
						Type: ExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;\nb:2;}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
					"b": {
						Type: ExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;\n\rb:2;}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
					"b": {
						Type: ExpInt,
						Int:  2,
					},
				},
			},
		},
		{
			str: "{a:1;   b:1.1 ;  c:\"str\";d:'str'; e: true;\n f: false;\n\r g:nil; h:var1; i:a+b+3;  j:  func1();  k:func2(a,b);}",
			exp: Exp{
				Type: ExpMap,
				Map: map[string]*Exp{
					"a": {
						Type: ExpInt,
						Int:  1,
					},
					"b": {
						Type:  ExpFloat,
						Float: 1.1,
					},
					"c": {
						Type: ExpStr,
						Str:  "str",
					},
					"d": {
						Type: ExpStr,
						Str:  "str",
					},
					"e": {
						Type: ExpBool,
						Bool: true,
					},
					"f": {
						Type: ExpBool,
						Bool: false,
					},
					"g": {
						Type: ExpNil,
					},
					"h": {
						Type:     ExpVar,
						Variable: "var1",
					},
					"i": {
						Type:     ExpCalc,
						Operator: "+",
						Left: &Exp{
							Type:     ExpCalc,
							Operator: "+",
							Left: &Exp{
								Type:     ExpVar,
								Variable: "a",
							},
							Right: &Exp{
								Type:     ExpVar,
								Variable: "b",
							},
						},
						Right: &Exp{
							Type: ExpInt,
							Int:  3,
						},
					},
					"j": {
						Type:       ExpFunc,
						FuncName:   "func1",
						FuncParams: []*Exp{},
					},
					"k": {
						Type:     ExpFunc,
						FuncName: "func2",
						FuncParams: []*Exp{
							{
								Type:     ExpVar,
								Variable: "a",
							},
							{
								Type:     ExpVar,
								Variable: "b",
							},
						},
					},
				},
			},
		},
		{
			str: "var1 ? a : b",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "?",
				TenaryCondition: &Exp{
					Type:     ExpVar,
					Variable: "var1",
				},
				Left: &Exp{
					Type:     ExpVar,
					Variable: "a",
				},
				Right: &Exp{
					Type:     ExpVar,
					Variable: "b",
				},
			},
		},
		{
			str: "a + b > 3 ? c + d : e * f",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "?",
				TenaryCondition: &Exp{
					Type:     ExpCalc,
					Operator: ">",
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "+",
						Left: &Exp{
							Type:     ExpVar,
							Variable: "a",
						},
						Right: &Exp{
							Type:     ExpVar,
							Variable: "b",
						},
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  3,
					},
				},
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "+",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "c",
					},
					Right: &Exp{
						Type:     ExpVar,
						Variable: "d",
					},
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "*",
					Left: &Exp{
						Type:     ExpVar,
						Variable: "e",
					},
					Right: &Exp{
						Type:     ExpVar,
						Variable: "f",
					},
				},
			},
		},
		{
			str: "a + b > 3 ? (c + d < 5 ? -1 : -2) : (e * f < 10 ? -3 : -4)",
			exp: Exp{
				Type:     ExpCalc,
				Operator: "?",
				TenaryCondition: &Exp{
					Type:     ExpCalc,
					Operator: ">",
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "+",
						Left: &Exp{
							Type:     ExpVar,
							Variable: "a",
						},
						Right: &Exp{
							Type:     ExpVar,
							Variable: "b",
						},
					},
					Right: &Exp{
						Type: ExpInt,
						Int:  3,
					},
				},
				Left: &Exp{
					Type:     ExpCalc,
					Operator: "?",
					TenaryCondition: &Exp{
						Type:     ExpCalc,
						Operator: "<",
						Left: &Exp{
							Type:     ExpCalc,
							Operator: "+",
							Left: &Exp{
								Type:     ExpVar,
								Variable: "c",
							},
							Right: &Exp{
								Type:     ExpVar,
								Variable: "d",
							},
						},
						Right: &Exp{
							Type: ExpInt,
							Int:  5,
						},
					},
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "-",
						Left:     nil,
						Right: &Exp{
							Type: ExpInt,
							Int:  1,
						},
					},
					Right: &Exp{
						Type:     ExpCalc,
						Operator: "-",
						Left:     nil,
						Right: &Exp{
							Type: ExpInt,
							Int:  2,
						},
					},
				},
				Right: &Exp{
					Type:     ExpCalc,
					Operator: "?",
					TenaryCondition: &Exp{
						Type:     ExpCalc,
						Operator: "<",
						Left: &Exp{
							Type:     ExpCalc,
							Operator: "*",
							Left: &Exp{
								Type:     ExpVar,
								Variable: "e",
							},
							Right: &Exp{
								Type:     ExpVar,
								Variable: "f",
							},
						},
						Right: &Exp{
							Type: ExpInt,
							Int:  10,
						},
					},
					Left: &Exp{
						Type:     ExpCalc,
						Operator: "-",
						Left:     nil,
						Right: &Exp{
							Type: ExpInt,
							Int:  3,
						},
					},
					Right: &Exp{
						Type:     ExpCalc,
						Operator: "-",
						Left:     nil,
						Right: &Exp{
							Type: ExpInt,
							Int:  4,
						},
					},
				},
			},
		},
		{
			str: "a+b >",
			err: "exp.incompleteExpression",
		},
		{
			str: "a+b<",
			err: "exp.incompleteExpression",
		},
		{
			str: "'hello",
			err: "exp.mismatchedSingleQuotationMark",
		},
		{
			str: "\"hello",
			err: "exp.mismatchedDoubleQuotationMark",
		},
		{
			str: "hello'",
			err: "exp.mismatchedSingleQuotationMark",
		},
		{
			str: "hello\"",
			err: "exp.mismatchedDoubleQuotationMark",
		},
		{
			str: "(a+b",
			err: "exp.mismatchedParenthesis",
		},
		{
			str: "a+b)",
			err: "exp.mismatchedParenthesis",
		},
		{
			str: "(a+b))+c",
			err: "exp.mismatchedParenthesis",
		},
		{
			str: "a+ b>3 ? c:",
			err: "exp.incompleteExpression",
		},
		{
			str: "func1(,b)",
			err: "exp.expectingParameter",
		},
		{
			str: "func1(a,)",
			err: "exp.expectingParameter",
		},
		{
			str: "func1(a",
			err: "exp.mismatchedParenthesis",
		},
		{
			str: "func1(a, func2(a)",
			err: "exp.mismatchedParenthesis",
		},
	}
)

func TestParseExp(t *testing.T) {
	var realExp *Exp
	var err error
	var str string
	for _, testCase := range parseExpTestCases {
		str = testCase.str
		t.Logf("expression: %s\n", str)
		realExp, err = ParseExp(str)
		if err != nil {
			if tpe, ok := err.(*DdlError); ok {
				if testCase.err != "" && tpe.etype == testCase.err {
					t.Log(tpe.Error())
					continue
				}
			}
			t.Fatalf("error: %s", err)
		} else if !realExp.Equal(&testCase.exp) {
			spew.Dump(*realExp)
			spew.Dump(testCase.exp)
			t.Fatalf("the expression is not parsed as expected\n")
		}
	}
}
