package ddl

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	expPattern = struct {
		reservedWord *regexp.Regexp
		variable     *regexp.Regexp
		intNum       *regexp.Regexp
		floatNum     *regexp.Regexp
		operator     *regexp.Regexp
		function     *regexp.Regexp
		whitespace   *regexp.Regexp
	}{
		reservedWord: regexp.MustCompile("^(true|false|nil)([^a-zA-Z0-9_]+|$)"),
		variable:     regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*"),
		intNum:       regexp.MustCompile("^[0-9]+"),
		floatNum:     regexp.MustCompile(`^[0-9]+\.[0-9]+`),
		operator:     regexp.MustCompile(`^(\{|\}|\(|\)|\[|\]|\+|-|\*|/|%|==|!=|>=|<=|>|<|&&|\|\||!|,|\.|:|;|\?)`),
		function:     regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(`),
		whitespace:   regexp.MustCompile(`^\s+`),
	}

	operatorPriority = map[string]int{
		".": 6,

		"!":        5,
		"negative": 5,

		"*": 4,
		"/": 4,
		"%": 4,

		"+": 3,
		"-": 3,

		">":  2,
		">=": 2,
		"<":  2,
		"<=": 2,
		"==": 2,
		"!=": 2,

		"&&": 1,
		"||": 1,

		"?": -1,
		":": 0,
	}
)

type Exp struct {
	Type            ExpType
	Int             int64
	Float           float64
	Str             string
	Bool            bool
	FuncName        string
	FuncParams      []*Exp
	Variable        string
	Operator        string
	Map             map[string]*Exp
	Left            *Exp
	Right           *Exp
	TenaryCondition *Exp
	Pos             int
	Interface       interface{}
}

type ExpType int

const (
	ExpStr ExpType = iota
	ExpInt
	ExpFloat
	ExpBool
	ExpNil
	ExpVar
	ExpOperator
	ExpFunc
	ExpCalc
	ExpMap
	ExpInterface
)

var ExpTypeName = map[ExpType]string{
	ExpStr:       "string",
	ExpInt:       "int",
	ExpFloat:     "float",
	ExpBool:      "bool",
	ExpNil:       "nil",
	ExpVar:       "variable",
	ExpOperator:  "operator",
	ExpFunc:      "function",
	ExpCalc:      "calculation",
	ExpMap:       "map",
	ExpInterface: "interface",
}

type ExpOperatorDirection int

const (
	LTR ExpOperatorDirection = iota
	RTL
)

func readExp(str string) ([]Exp, error) {
	exps := make([]Exp, 0)

	byteArr := []byte(str)
	blen := len(byteArr)
	pos := 0

	for {
		ch := str[pos]
		left := str[pos:]

		// string literal
		if ch == '\'' {
			match := false
			for end := pos + 1; end < blen; end++ {
				if str[end] == '\'' && str[end-1] != '\\' {
					exps = append(exps, Exp{
						Type: ExpStr,
						Str:  strings.ReplaceAll(str[pos+1:end], "\\'", "'"),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return exps, NewDdlParseError(str, "exp.mismatchedSingleQuotationMark", pos)
			}

			// string literal
		} else if ch == '"' {
			match := false
			for end := pos + 1; end < blen; end++ {
				if str[end] == '"' && str[end-1] != '\\' {
					exps = append(exps, Exp{
						Type: ExpStr,
						Str:  strings.ReplaceAll(str[pos+1:end], "\\\"", "\""),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return exps, NewDdlParseError(str, "exp.mismatchedDoubleQuotationMark", pos)
			}

			// true, false, nil
		} else if matches := expPattern.reservedWord.FindStringSubmatch(left); len(matches) > 0 {
			word := matches[1]
			switch word {
			case "true":
				exps = append(exps, Exp{
					Type: ExpBool,
					Bool: true,
					Pos:  pos,
				})
			case "false":
				exps = append(exps, Exp{
					Type: ExpBool,
					Bool: false,
					Pos:  pos,
				})
			case "nil":
				exps = append(exps, Exp{
					Type: ExpNil,
					Pos:  pos,
				})
			}
			pos += len(word)

			// func
		} else if matches := expPattern.function.FindStringSubmatch(left); len(matches) > 0 {
			exps = append(exps, Exp{
				Type:     ExpFunc,
				FuncName: matches[1],
				Pos:      pos,
			})
			pos += len(matches[0])

			// variable
		} else if matches := expPattern.variable.FindStringSubmatch(left); len(matches) > 0 {
			variable := matches[0]
			count := len(exps)
			if count > 0 && exps[count-1].Type == ExpOperator && exps[count-1].Operator == "." {
				exps = append(exps, Exp{
					Type: ExpStr,
					Str:  variable,
					Pos:  pos,
				})
			} else {
				exps = append(exps, Exp{
					Type:     ExpVar,
					Variable: variable,
					Pos:      pos,
				})
			}
			pos += len(matches[0])

			// float
		} else if matches := expPattern.floatNum.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[0], 64)
			exps = append(exps, Exp{
				Type:  ExpFloat,
				Float: num,
				Pos:   pos,
			})
			pos += len(matches[0])

			// int
		} else if matches := expPattern.intNum.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseInt(matches[0], 10, 64)
			exps = append(exps, Exp{
				Type: ExpInt,
				Int:  num,
				Pos:  pos,
			})
			pos += len(matches[0])

			// operator
		} else if matches := expPattern.operator.FindStringSubmatch(left); len(matches) > 0 {
			optr := matches[0]
			if optr == "-" {
				if len(exps) == 0 {
					optr = "negative"
				} else if exps[len(exps)-1].Type == ExpOperator {
					lastOptr := exps[len(exps)-1].Operator
					if lastOptr != ")" && lastOptr != "]" {
						optr = "negative"
					}
				}
			}
			exps = append(exps, Exp{
				Type:     ExpOperator,
				Operator: optr,
				Pos:      pos,
			})
			pos += len(matches[0])

			// space
		} else if matches := expPattern.whitespace.FindStringSubmatch(left); len(matches) > 0 {
			pos += len(matches[0])

			// others
		} else {
			return exps, NewDdlParseError(str, "exp.unexpectedToken", pos)
		}

		if pos > blen-1 {
			break
		}
	}

	return exps, nil
}

func generateExpTree(exps []Exp) (*Exp, error) {
	optrStack := make([]*Exp, 0)
	opndStack := make([]*Exp, 0)

	pos := 0
	for {
		if pos > len(exps)-1 {
			break
		}

		exp := exps[pos]

		// {...}
		if exp.Type == ExpOperator && exp.Operator == "{" {
			bracketEnd := pos
			bracketNum := 1
			key := ""
			valBegin := 0
			valEnd := 0
			expMap := make(map[string]*Exp)
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == ExpOperator && exps[idx].Operator == "{" {
					bracketNum += 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == "}" {
					bracketNum -= 1
					valEnd = idx - 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == ":" {
					if exps[idx-1].Type != ExpVar {
						return nil, errors.New("there must be a key before ':'")
					}
					if exps[idx-2].Type != ExpOperator || (exps[idx-2].Operator != "{" && exps[idx-2].Operator != ";") {
						return nil, errors.New("invalid key before ':'")
					}
					key = exps[idx-1].Variable
					valBegin = idx + 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == ";" {
					valEnd = idx - 1
				}

				if bracketNum <= 1 && valEnd-valBegin >= 0 && key != "" {
					val, err := generateExpTree(exps[valBegin : valEnd+1])
					if err != nil {
						return nil, err
					}
					expMap[key] = val
					key = ""
				}
				if bracketNum == 0 {
					bracketEnd = idx
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewDdlParseError("", "exp.mismatchedCurlyBrace", exp.Pos)
			}

			opndStack = append(opndStack, &Exp{
				Type: ExpMap,
				Map:  expMap,
			})
			pos = bracketEnd + 1

			// (...)
		} else if exp.Type == ExpOperator && exp.Operator == "(" {
			bracketBegin := pos
			bracketEnd := pos
			bracketNum := 1
			for idx := pos + 1; idx < len(exps); idx++ {
				if (exps[idx].Type == ExpOperator && exps[idx].Operator == "(") || exps[idx].Type == ExpFunc {
					bracketNum += 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == ")" {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					bracketEnd = idx
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewDdlParseError("", "exp.mismatchedParenthesis", exp.Pos)
			}
			parsedExp, err := generateExpTree(exps[bracketBegin+1 : bracketEnd])
			if err != nil {
				return nil, err
			}
			opndStack = append(opndStack, parsedExp)
			pos = bracketEnd + 1

			// [...]
		} else if exp.Type == ExpOperator && exp.Operator == "[" {
			bracketBegin := pos
			bracketEnd := pos
			bracketNum := 1
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == ExpOperator && exps[idx].Operator == "[" {
					bracketNum += 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == "]" {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					bracketEnd = idx
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewDdlParseError("", "exp.mismatchedSquareBracket", exp.Pos)
			}
			parsedExp, err := generateExpTree(exps[bracketBegin+1 : bracketEnd])
			if err != nil {
				return nil, err
			}
			lastOpnd := opndStack[len(opndStack)-1]
			opndStack[len(opndStack)-1] = &Exp{
				Type:     ExpCalc,
				Operator: "[",
				Left:     lastOpnd,
				Right:    parsedExp,
			}
			/*
				opndStack = append(opndStack, parsedExp)
				optrStack = append(optrStack, &exp)
			*/
			pos = bracketEnd + 1

			// func
		} else if exp.Type == ExpFunc {
			argBegin := pos + 1
			bracketNum := 1
			args := make([]*Exp, 0)
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == ExpOperator && exps[idx].Operator == "(" {
					bracketNum += 1
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == ")" {
					bracketNum -= 1
					if idx > argBegin {
						arg, err := generateExpTree(exps[argBegin:idx])
						if err != nil {
							return nil, err
						}
						args = append(args, arg)
					}
				} else if exps[idx].Type == ExpOperator && exps[idx].Operator == "," && bracketNum == 1 {
					if idx == argBegin {
						return nil, NewDdlParseError("", "exp.expectingParameter", exps[idx].Pos)
					}
					if len(exps) > idx+1 && exps[idx+1].Operator == ")" {
						return nil, NewDdlParseError("", "exp.expectingParameter", exps[idx+1].Pos)
					}
					arg, err := generateExpTree(exps[argBegin:idx])
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
					argBegin = idx + 1
				}
				if bracketNum == 0 {
					pos = idx + 1
					exp.FuncParams = args
					opndStack = append(opndStack, &exp)
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewDdlParseError("", "exp.mismatchedParentheses", exp.Pos+len(exp.FuncName))
			}

			// operator
		} else if exp.Type == ExpOperator {
			if len(optrStack) == 0 {
				optrStack = append(optrStack, &exp)
				pos += 1
			} else {
				lastOptr := optrStack[len(optrStack)-1].Operator
				curOptr := exp.Operator
				if operatorPriority[curOptr] <= operatorPriority[lastOptr] {
					var err error
					opndStack, optrStack, err = calculate(opndStack, optrStack)
					if err != nil {
						return nil, err
					}
				} else {
					optrStack = append(optrStack, &exp)
					pos += 1
				}
			}

			// others
		} else {
			opndStack = append(opndStack, &exp)
			pos += 1
		}
	}

	for {
		if len(optrStack) == 0 {
			break
		}

		var err error
		opndStack, optrStack, err = calculate(opndStack, optrStack)
		if err != nil {
			return nil, err
		}
	}

	return opndStack[0], nil
}

func calculate(opndStack []*Exp, optrStack []*Exp) ([]*Exp, []*Exp, error) {
	if len(optrStack) == 0 {
		return opndStack, optrStack, nil
	}

	optr := optrStack[len(optrStack)-1]
	switch optr.Operator {
	case ")":
		return opndStack, optrStack, NewDdlParseError("", "exp.mismatchedParenthesis", optr.Pos)

	case "]":
		return opndStack, optrStack, NewDdlParseError("", "exp.mismatchedSquareBracket", optr.Pos)

	case "}":
		return opndStack, optrStack, NewDdlParseError("", "exp.mismatchedCurlyBracket", optr.Pos)

	case "negative":
		if len(opndStack) == 0 {
			return opndStack, optrStack, NewDdlParseError("", "exp.incompleteExpression", optr.Pos)
		}
		oexp := opndStack[len(opndStack)-1]
		nexp := &Exp{
			Type:     ExpCalc,
			Operator: "-",
			Left:     nil,
			Right:    oexp,
		}
		opndStack[len(opndStack)-1] = nexp
		optrStack = optrStack[:len(optrStack)-1]

	case "!":
		if len(opndStack) == 0 {
			return opndStack, optrStack, NewDdlParseError("", "exp.incompleteExpression", optr.Pos)
		}
		oexp := opndStack[len(opndStack)-1]
		nexp := &Exp{
			Type:     ExpCalc,
			Operator: "!",
			Left:     nil,
			Right:    oexp,
		}
		opndStack[len(opndStack)-1] = nexp
		optrStack = optrStack[:len(optrStack)-1]

	case "?":
		if len(opndStack) < 2 {
			return opndStack, optrStack, NewDdlParseError("", "exp.incompleteExpression", optr.Pos)
		}
		exp1 := opndStack[len(opndStack)-1]
		exp2 := opndStack[len(opndStack)-2]
		if exp1.Type != ExpCalc || exp1.Operator != ":" {
			return opndStack, optrStack, NewDdlParseError("", "exp.invalidTenaryExpression", optr.Pos)
		}
		nexp := &Exp{
			Type:            ExpCalc,
			Operator:        "?",
			TenaryCondition: exp2,
			Left:            exp1.Left,
			Right:           exp1.Right,
		}
		opndStack = opndStack[:len(opndStack)-2]
		opndStack = append(opndStack, nexp)
		optrStack = optrStack[:len(optrStack)-1]

	default:
		if len(opndStack) < 2 {
			return opndStack, optrStack, NewDdlParseError("", "exp.incompleteExpression", optr.Pos)
		}
		exp1 := opndStack[len(opndStack)-1]
		exp2 := opndStack[len(opndStack)-2]
		nexp := &Exp{
			Type:     ExpCalc,
			Operator: optr.Operator,
			Left:     exp2,
			Right:    exp1,
		}
		opndStack = opndStack[:len(opndStack)-2]
		opndStack = append(opndStack, nexp)
		optrStack = optrStack[:len(optrStack)-1]
	}

	return opndStack, optrStack, nil
}

func ParseExp(str string) (*Exp, error) {
	exps, err := readExp(str)
	if err != nil {
		return nil, err
	}
	tree, err2 := generateExpTree(exps)
	if err2 != nil {
		if tpe, ok := err2.(*DdlParseError); ok {
			tpe.SetDdl(str)
		}
		return nil, err2
	}

	return tree, nil
}

func (a *Exp) Equal(b *Exp) bool {
	if a == nil && b == nil {
		return true
	}

	if (a != nil && b == nil) ||
		(a == nil && b != nil) {
		return false
	}

	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case ExpInt:
		if a.Int != b.Int {
			return false
		}

	case ExpFloat:
		if math.Abs(a.Float-b.Float) > 1e-9 {
			return false
		}

	case ExpBool:
		if a.Bool != b.Bool {
			return false
		}

	case ExpStr:
		if a.Str != b.Str {
			return false
		}

	case ExpVar:
		if a.Variable != b.Variable {
			return false
		}

	case ExpOperator:
		if a.Operator != b.Operator {
			return false
		}

	case ExpFunc:
		if a.FuncName != b.FuncName || len(a.FuncParams) != len(b.FuncParams) {
			return false
		}
		for i := 0; i < len(a.FuncParams); i++ {
			if !a.FuncParams[i].Equal(b.FuncParams[i]) {
				return false
			}
		}

	case ExpMap:
		if len(a.Map) != len(b.Map) {
			return false
		}
		for k, v := range a.Map {
			if b.Map[k] == nil {
				return false
			}
			if !v.Equal(b.Map[k]) {
				return false
			}
		}

	case ExpCalc:
		if a.Operator != b.Operator ||
			(a.Left == nil && b.Left != nil) ||
			(a.Left != nil && b.Left == nil) ||
			(a.Right == nil && b.Right != nil) ||
			(a.Right != nil && b.Right == nil) ||
			(a.Left != nil && b.Left != nil && !a.Left.Equal(b.Left)) ||
			(a.Right != nil && b.Right != nil && !a.Right.Equal(b.Right)) {
			return false
		}
		if a.Operator == "?" && !a.TenaryCondition.Equal(b.TenaryCondition) {
			return false
		}

	case ExpInterface:
		aval := reflect.ValueOf(a.Interface)
		bval := reflect.ValueOf(b.Interface)
		if aval.Comparable() {
			return aval.Equal(bval)
		}
		if aval.Pointer() != bval.Pointer() {
			return false
		}
	}

	return true
}
