package rview

import (
	"testing"

	"github.com/TinyWisp/rview/ddl"
	"github.com/davecgh/go-spew/spew"
)

type CalcExpCase struct {
	exp    string
	expect string
	err    string
}

type AnimalStruct struct {
	Name string
}

var variableMap = map[string]interface{}{
	"chickenStruct":  AnimalStruct{Name: "chicken"},
	"duckStruct":     AnimalStruct{Name: "duck"},
	"int8var2":       int8(2),
	"int16var2":      int16(2),
	"int32var2":      int32(2),
	"int64var2":      int64(2),
	"uint8var2":      int8(2),
	"uint16var2":     int16(2),
	"uint32var2":     int32(2),
	"uint64var2":     int64(2),
	"float32var3":    float32(3.0),
	"float64var3":    float64(3.0),
	"stringvarhello": "hello",
	"stringvarworld": "world",
	"stringvarName":  "Name",
	"boolvarfalse":   false,
	"boolvartrue":    true,
	"nilvar":         nil,
	"mapStrStr": map[string]string{
		"hello": "hello",
		"world": "world",
	},
	"arrStr": []string{"hello", "world"},
	"arrInt": []int{10, 11, 12, 13},
}

func getVariable(name string) (interface{}, error) {
	if val, ok := variableMap[name]; ok {
		return val, nil
	}

	return nil, nil
}

var calcExpCases = []CalcExpCase{
	{
		exp:    `1`,
		expect: `1`,
	},
	{
		exp:    `123.5`,
		expect: `123.5`,
	},
	{
		exp:    `true`,
		expect: `true`,
	},
	{
		exp:    `nil`,
		expect: `nil`,
	},
	{
		exp:    `"hello"`,
		expect: `"hello"`,
	},

	// -----------------   +   ----------------------
	{
		exp:    `1+3`,
		expect: "4",
	},
	{
		exp:    `1+1.2`,
		expect: `2.2`,
	},
	{
		exp:    `1.2+1`,
		expect: `2.2`,
	},
	{
		exp:    `1.2+2.5`,
		expect: `3.7`,
	},
	{
		exp: `3+"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true+3`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   -   ----------------------
	{
		exp:    `2-1`,
		expect: `1`,
	},
	{
		exp:    `3-1.2`,
		expect: `1.8`,
	},
	{
		exp:    `1.2-1`,
		expect: `0.2`,
	},
	{
		exp:    `2.5-1.2`,
		expect: `1.3`,
	},
	{
		exp: `3-"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true-1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   *   ----------------------
	{
		exp:    `2*3`,
		expect: `6`,
	},
	{
		exp:    `3*1.2`,
		expect: `3.6`,
	},
	{
		exp:    `1.2*2`,
		expect: `2.4`,
	},
	{
		exp:    `1.2*2.5`,
		expect: `3.0`,
	},
	{
		exp: `3*"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true*1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   /   ----------------------
	{
		exp:    `4/2`,
		expect: `2`,
	},
	{
		exp:    `5/2`,
		expect: `2`,
	},
	{
		exp:    `3/1.2`,
		expect: `2.5`,
	},
	{
		exp:    `1.2/2`,
		expect: `0.6`,
	},
	{
		exp:    `1.2/2.5`,
		expect: `0.48`,
	},
	{
		exp: `3/"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true/1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   >   ----------------------
	{
		exp:    `4>3`,
		expect: `true`,
	},
	{
		exp:    `4>4`,
		expect: `false`,
	},
	{
		exp:    `4>5`,
		expect: `false`,
	},
	{
		exp:    `3>2.99999`,
		expect: `true`,
	},
	{
		exp:    `3>3.0`,
		expect: `false`,
	},
	{
		exp:    `3>3.00000001`,
		expect: `false`,
	},
	{
		exp:    `2.0000001>2`,
		expect: `true`,
	},
	{
		exp:    `2.000001>3`,
		expect: `false`,
	},
	{
		exp:    `2.0>2.0`,
		expect: `false`,
	},
	{
		exp:    `2.000001>2.0`,
		expect: `true`,
	},
	{
		exp:    `2.0000001>2.0000002`,
		expect: `false`,
	},
	{
		exp: `3>"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true>1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   >=   ----------------------
	{
		exp:    `4>=3`,
		expect: `true`,
	},
	{
		exp:    `4>=4`,
		expect: `true`,
	},
	{
		exp:    `4>=5`,
		expect: `false`,
	},
	{
		exp:    `3>=2.999999`,
		expect: `true`,
	},
	{
		exp:    `3>=3.0`,
		expect: `true`,
	},
	{
		exp:    `3>=3.000000001`,
		expect: `false`,
	},
	{
		exp:    `2.00000001>=2`,
		expect: `true`,
	},
	{
		exp:    `2.0>=2`,
		expect: `true`,
	},
	{
		exp:    `2.000001>=3`,
		expect: `false`,
	},
	{
		exp:    `2.0>=2.0`,
		expect: `true`,
	},
	{
		exp:    `2.0000001>=2.0`,
		expect: `true`,
	},
	{
		exp:    `2.00000001>=2.00000002`,
		expect: `false`,
	},
	{
		exp: `3>="hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true>=1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   <   ----------------------
	{
		exp:    `3<4`,
		expect: `true`,
	},
	{
		exp:    `4<4`,
		expect: `false`,
	},
	{
		exp:    `5<4`,
		expect: `false`,
	},
	{
		exp:    `2<2.0000001`,
		expect: `true`,
	},
	{
		exp:    `3<3.0`,
		expect: `false`,
	},
	{
		exp:    `3<3.00000001`,
		expect: `true`,
	},
	{
		exp:    `1.9999999<2`,
		expect: `true`,
	},
	{
		exp:    `3.0000001<3`,
		expect: `false`,
	},
	{
		exp:    `2.0<2.0`,
		expect: `false`,
	},
	{
		exp:    `2.0<2.0000001`,
		expect: `true`,
	},
	{
		exp:    `2.0000001<2.0`,
		expect: `false`,
	},
	{
		exp: `3<"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true<1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   <=   ----------------------
	{
		exp:    `3<=4`,
		expect: `true`,
	},
	{
		exp:    `4<=4`,
		expect: `true`,
	},
	{
		exp:    `5<=4`,
		expect: `false`,
	},
	{
		exp:    `3<=3.0000001`,
		expect: `true`,
	},
	{
		exp:    `3<=3.0`,
		expect: `true`,
	},
	{
		exp:    `3<=2.9999999`,
		expect: `false`,
	},
	{
		exp:    `1.99999999<=2`,
		expect: `true`,
	},
	{
		exp:    `2.0<=2`,
		expect: `true`,
	},
	{
		exp:    `2.00000001<=2`,
		expect: `false`,
	},
	{
		exp:    `2.0<=2.0`,
		expect: `true`,
	},
	{
		exp:    `1.99999999<=2.0`,
		expect: `true`,
	},
	{
		exp:    `2.0000002<=2.0000001`,
		expect: `false`,
	},
	{
		exp: `3<="hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `true>=1`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   ==   ----------------------
	{
		exp:    `3==3`,
		expect: `true`,
	},
	{
		exp:    `3==4`,
		expect: `false`,
	},
	{
		exp:    `3==3.0`,
		expect: `true`,
	},
	{
		exp:    `3==3.0000001`,
		expect: `false`,
	},
	{
		exp:    `2.0==2`,
		expect: `true`,
	},
	{
		exp:    `2.00000001==2`,
		expect: `false`,
	},
	{
		exp:    `2.0==2.0`,
		expect: `true`,
	},
	{
		exp:    `1.9999999==2.0`,
		expect: `false`,
	},
	{
		exp:    `nil==nil`,
		expect: `true`,
	},
	{
		exp:    `false==false`,
		expect: `true`,
	},
	{
		exp:    `true==true`,
		expect: `true`,
	},
	{
		exp:    `false==true`,
		expect: `false`,
	},
	{
		exp:    `"hello"=="hello"`,
		expect: `true`,
	},
	{
		exp:    `"hello1"=="hello"`,
		expect: `false`,
	},
	{
		exp: `1=="hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1==true`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `nil==false`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   !=   ----------------------
	{
		exp:    `3!=3`,
		expect: `false`,
	},
	{
		exp:    `3!=4`,
		expect: `true`,
	},
	{
		exp:    `3!=3.0`,
		expect: `false`,
	},
	{
		exp:    `3!=3.0000001`,
		expect: `true`,
	},
	{
		exp:    `2.0!=2`,
		expect: `false`,
	},
	{
		exp:    `2.0000001!=2`,
		expect: `true`,
	},
	{
		exp:    `2.0!=2.0`,
		expect: `false`,
	},
	{
		exp:    `1.999999!=2.0`,
		expect: `true`,
	},
	{
		exp:    `nil!=nil`,
		expect: `false`,
	},
	{
		exp:    `false!=false`,
		expect: `false`,
	},
	{
		exp:    `true!=true`,
		expect: `false`,
	},
	{
		exp:    `false!=true`,
		expect: `true`,
	},
	{
		exp:    `"hello"!="hello"`,
		expect: `false`,
	},
	{
		exp:    `"hello1"!="hello"`,
		expect: `true`,
	},
	{
		exp: `1!="hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1!=true`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `nil!=true`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   &&   ----------------------
	{
		exp:    `false && true`,
		expect: `false`,
	},
	{
		exp:    `true && false`,
		expect: `false`,
	},
	{
		exp:    `true && true`,
		expect: `true`,
	},
	{
		exp:    `false && false`,
		expect: `false`,
	},
	{
		exp: `1 && 2`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1 && "hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1 && nil`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   ||   ----------------------
	{
		exp:    `false || true`,
		expect: `true`,
	},
	{
		exp:    `true || false`,
		expect: `true`,
	},
	{
		exp:    `true || true`,
		expect: `true`,
	},
	{
		exp:    `false || false`,
		expect: `false`,
	},
	{
		exp: `1 || 2`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1 || "hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `1 || nil`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------   !   ----------------------
	{
		exp:    `!false`,
		expect: `true`,
	},
	{
		exp:    `!true`,
		expect: `false`,
	},
	{
		exp: `!1`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `!"hello"`,
		err: "calc.operandTypeMismatch",
	},
	{
		exp: `!nil`,
		err: "calc.operandTypeMismatch",
	},

	// -----------------------  ? : -------------------------
	{
		exp:    `true ? "hello" : "world"`,
		expect: `"hello"`,
	},
	{
		exp:    `false ? "hello" : "world"`,
		expect: `"world"`,
	},
	{
		exp:    `true ? 2 : 1`,
		expect: `2`,
	},
	{
		exp:    `false ? 2 : 1`,
		expect: `1`,
	},
	{
		exp: `false ? "hello" : 1`,
		err: `calc.ternaryDataNotSameType`,
	},
	{
		exp: `3 ? "hello" : "world"`,
		err: `calc.invalidTernaryCondition`,
	},
	{
		exp: `nil ? "hello" : "world"`,
		err: `calc.invalidTernaryCondition`,
	},
	{
		exp: `"abc" ? "hello" : "world"`,
		err: `calc.invalidTernaryCondition`,
	},

	// --------------------- variable------------------------
	{
		exp:    `int8var2`,
		expect: `2`,
	},
	{
		exp:    `int16var2`,
		expect: `2`,
	},
	{
		exp:    `int32var2`,
		expect: `2`,
	},
	{
		exp:    `int64var2`,
		expect: `2`,
	},
	{
		exp:    `uint8var2`,
		expect: `2`,
	},
	{
		exp:    `uint16var2`,
		expect: `2`,
	},
	{
		exp:    `uint32var2`,
		expect: `2`,
	},
	{
		exp:    `uint64var2`,
		expect: `2`,
	},
	{
		exp:    `float32var3`,
		expect: `3.0`,
	},
	{
		exp:    `float64var3`,
		expect: `3.0`,
	},
	{
		exp:    `boolvarfalse`,
		expect: `false`,
	},
	{
		exp:    `boolvartrue`,
		expect: `true`,
	},
	{
		exp:    `nilvar`,
		expect: `nil`,
	},

	// ----------------------- . ---------------------------
	{
		exp:    `chickenStruct.Name`,
		expect: `"chicken"`,
	},
	{
		exp:    `mapStrStr.hello`,
		expect: `"hello"`,
	},
	{
		exp:    `mapStrStr.world`,
		expect: `"world"`,
	},

	// ----------------------- [] ---------------------------
	{
		exp:    `chickenStruct["Name"]`,
		expect: `"chicken"`,
	},
	{
		exp:    `chickenStruct[stringvarName]`,
		expect: `"chicken"`,
	},
	{
		exp:    `mapStrStr["world"]`,
		expect: `"world"`,
	},
	{
		exp:    `arrStr[0]`,
		expect: `"hello"`,
	},
	{
		exp:    `arrStr[1]`,
		expect: `"world"`,
	},
	{
		exp:    `arrInt[1]`,
		expect: "11",
	},

	// --------  + - * / ( ) > >= < <= == != && ! || . [] ------------
	{
		exp:    `3+2-1`,
		expect: `4`,
	},
	{
		exp:    `3+2*(5-1*3)/(1+1)`,
		expect: `5`,
	},
	{
		exp:    `3*5/3 + 2/2 - 3/3`,
		expect: `5`,
	},
	{
		exp:    `3*(5+2+1) - 6/(2+1)`,
		expect: `22`,
	},
	{
		exp:    `1 + 0.1 * 3 - 0.1 * 3 + 0.00001 > 1`,
		expect: `true`,
	},
	{
		exp:    `1 + 0.1 * 3 - 0.1 * 3 + 0.000001 >= 1`,
		expect: `true`,
	},
	{
		exp:    `1 + 0.00001 < 1 + 0.00001 * 2`,
		expect: `true`,
	},
	{
		exp:    `1 + 0.00001 <= 1 + 0.00001 * 2`,
		expect: `true`,
	},
	{
		exp:    `3*5/3 == 5`,
		expect: `true`,
	},
	{
		exp:    `3*(5+3) == 24`,
		expect: `true`,
	},
	{
		exp:    `3*5/3 != 5`,
		expect: `false`,
	},
	{
		exp:    `3*(5+3) != 24`,
		expect: `false`,
	},
	{
		exp:    `arrInt[0] + arrInt[1]`,
		expect: `21`,
	},
	{
		exp:    `!(3 + 5 > 2)`,
		expect: `false`,
	},
	{
		exp:    `!(3 + 5 == 2)`,
		expect: `true`,
	},
	{
		exp:    `stringvarhello == "hello" && stringvarworld == "world" && int32var2 > 0 && 5 + 3*2 > 10`,
		expect: `true`,
	},
	{
		exp:    `stringvarhello == "hello" && stringvarworld == "world" && int32var2 > 0 && 5 + 3*2 > 10`,
		expect: `true`,
	},
	{
		exp:    `arrInt[1] == 11 && arrStr[0] == "hello" && 3 > 2 && int32var2 > 0 && 9 + 8 * 7 > 8 + 8 * 7`,
		expect: `true`,
	},
}

func TestCalcExp(t *testing.T) {
	for _, testCase := range calcExpCases {
		t.Log(testCase.exp)

		exp, err1 := ddl.ParseExp(testCase.exp)
		if err1 != nil {
			t.Fatal(err1)
		}

		res, err2 := CalcExp(exp, getVariable)
		if err2 != nil {
			if fmtErr, ok := err2.(*FmtError); ok {
				if fmtErr.etype != testCase.err {
					t.Log(err2)
					t.Fatal("the error occured during the executing of this expression is not as expected. ")
				}
				continue
			}
		}

		expect, err3 := ddl.ParseExp(testCase.expect)
		if err3 != nil {
			t.Fatal(err3)
		}

		if !expect.Equal(res) {
			spew.Dump(expect)
			spew.Dump(res)
			t.Fatal("this expression is not calculated as expected")
		}
	}
}
