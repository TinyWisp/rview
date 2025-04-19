package rview

import (
	"testing"

	"github.com/TinyWisp/rview/ddl"
	"github.com/davecgh/go-spew/spew"
)

type CalcExpCase struct {
	exp    *ddl.Exp
	expect *ddl.Exp
	err    string
}

type AnimalStruct struct {
	Name string
}

var chickenStruct1 = AnimalStruct{Name: "chicken"}
var chickenStruct2 = AnimalStruct{Name: "chicken"}
var duckStruct1 = AnimalStruct{Name: "duck"}
var intvar1 int = 2
var floatvar1 float64 = 2.0
var stringvar1 string = "world"

var calcExpCases = []CalcExpCase{
	{
		exp: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  1,
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  1,
		},
	},
	{
		exp: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 123.5,
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 123.5,
		},
	},
	{
		exp: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type: ddl.ExpNil,
		},
		expect: &ddl.Exp{
			Type: ddl.ExpNil,
		},
	},
	{
		exp: &ddl.Exp{
			Type: ddl.ExpStr,
			Str:  "hello",
		},
		expect: &ddl.Exp{
			Type: ddl.ExpStr,
			Str:  "hello",
		},
	},
	{
		exp: &ddl.Exp{
			Type:      ddl.ExpInterface,
			Interface: "hello",
		},
		expect: &ddl.Exp{
			Type:      ddl.ExpInterface,
			Interface: "hello",
		},
	},
	{
		exp: &ddl.Exp{
			Type:      ddl.ExpInterface,
			Interface: chickenStruct1,
		},
		expect: &ddl.Exp{
			Type:      ddl.ExpInterface,
			Interface: chickenStruct1,
		},
	},

	// -----------------   +   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  -1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 2.2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 2.2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.5,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 3.7,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "+",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   -   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  1,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 1.8,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 0.2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.5,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: -1.3,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "-",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   *   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  -3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  -6,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 3.6,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 2.4,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.5,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 3,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "*",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   /   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  -2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  -2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  5,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  2,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 2.5,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 0.6,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.2,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.5,
			},
		},
		expect: &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: 0.48,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "/",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   >   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  5,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.99999,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.00000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.00001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.00001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000002,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   >=   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  5,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.99999,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.00000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.00001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.00001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000002,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   <   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  5,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.00000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.99999,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.000001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   <=   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  5,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.00001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.9999999,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.999999,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.00001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.999999,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000002,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "<=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: ">=",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   ==   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.999999,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpNil,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpNil,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello1",
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &floatvar1,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "==",
			Left: &ddl.Exp{
				Type: ddl.ExpNil,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   !=   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  4,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  3,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 3.000001,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.000001,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 1.999999,
			},
			Right: &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: 2.0,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpNil,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpNil,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello1",
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &chickenStruct2,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &intvar1,
			},
			Right: &ddl.Exp{
				Type:      ddl.ExpInterface,
				Interface: &floatvar1,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "!=",
			Left: &ddl.Exp{
				Type: ddl.ExpNil,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   &&   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "&&",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpNil,
			},
		},
		err: "calc.operandTypeMismatch",
	},

	// -----------------   ||   ----------------------
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: true,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: false,
			},
		},
		expect: &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		},
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  2,
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpStr,
				Str:  "hello",
			},
		},
		err: "calc.operandTypeMismatch",
	},
	{
		exp: &ddl.Exp{
			Type:     ddl.ExpCalc,
			Operator: "||",
			Left: &ddl.Exp{
				Type: ddl.ExpInt,
				Int:  1,
			},
			Right: &ddl.Exp{
				Type: ddl.ExpNil,
			},
		},
		err: "calc.operandTypeMismatch",
	},
}

func TestCalcExp(t *testing.T) {
	varGetter := func(name string) (interface{}, error) {
		return nil, nil
	}

	t.Log(len(calcExpCases))
	for idx, testCase := range calcExpCases {
		t.Log(idx)
		spew.Dump(testCase)
		res, err := CalcExp(testCase.exp, varGetter)
		if err == nil && !testCase.expect.Equal(res) {
			t.Log("---------exp:")
			spew.Dump(testCase.exp)
			t.Log("---------expect:")
			spew.Dump(testCase.expect)
			t.Log("---------result:")
			spew.Dump(res)
			t.Fatal("this expression is not calculated as expected")
		} else if err != nil {
			if fmtErr, ok := err.(*FmtError); ok {
				if fmtErr.etype != testCase.err {
					t.Log("---------exp:")
					spew.Dump(testCase.exp)
					t.Log("---------err:")
					t.Log(testCase.err)
					t.Log("---------real err:")
					spew.Dump(fmtErr)
					t.Fatal("this expression is not calculated as expected")
				}
			} else {
				t.Log(err)
				t.Fatal("this expression is not calculated as expected")
			}
		}
	}
}
