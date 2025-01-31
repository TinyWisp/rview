package template

import (
	"fmt"

	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

type TplExp struct {
	Type            TplExpType
	Int             int64
	Float           float64
	Str             string
	Bool            bool
	FuncName        string
	FuncParams      []*TplExp
	Variable        string
	Operator        string
	Map             map[string]*TplExp
	Left            *TplExp
	Right           *TplExp
	TenaryCondition *TplExp
	Pos             int
}

type TplExpType int

const (
	TplExpStr TplExpType = iota
	TplExpInt
	TplExpFloat
	TplExpBool
	TplExpNil
	TplExpVar
	TplExpOperator
	TplExpFunc
	TplExpCalc
	TplExpMap
)

func readTplExp(str string) ([]TplExp, error) {
	exps := make([]TplExp, 0)

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
					exps = append(exps, TplExp{
						Type: TplExpStr,
						Str:  strings.ReplaceAll(str[pos+1:end], "\\'", "'"),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return exps, NewTplParseError(str, "exp.mismatchedSingleQuotationMark", pos)
			}

			// string literal
		} else if ch == '"' {
			match := false
			for end := pos + 1; end < blen; end++ {
				if str[end] == '"' && str[end-1] != '\\' {
					exps = append(exps, TplExp{
						Type: TplExpStr,
						Str:  strings.ReplaceAll(str[pos+1:end], "\\\"", "\""),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return exps, NewTplParseError(str, "exp.mismatchedDoubleQuotationMark", pos)
			}

			// true, false, nil
		} else if matches := expPattern.reservedWord.FindStringSubmatch(left); len(matches) > 0 {
			word := matches[1]
			switch word {
			case "true":
				exps = append(exps, TplExp{
					Type: TplExpBool,
					Bool: true,
					Pos:  pos,
				})
			case "false":
				exps = append(exps, TplExp{
					Type: TplExpBool,
					Bool: false,
					Pos:  pos,
				})
			case "nil":
				exps = append(exps, TplExp{
					Type: TplExpNil,
					Pos:  pos,
				})
			}
			pos += len(word)

			// func
		} else if matches := expPattern.function.FindStringSubmatch(left); len(matches) > 0 {
			exps = append(exps, TplExp{
				Type:     TplExpFunc,
				FuncName: matches[1],
				Pos:      pos,
			})
			pos += len(matches[0])

			// variable
		} else if matches := expPattern.variable.FindStringSubmatch(left); len(matches) > 0 {
			variable := matches[0]
			count := len(exps)
			if count > 0 && exps[count-1].Type == TplExpOperator && exps[count-1].Operator == "." {
				exps = append(exps, TplExp{
					Type: TplExpStr,
					Str:  variable,
					Pos:  pos,
				})
			} else {
				exps = append(exps, TplExp{
					Type:     TplExpVar,
					Variable: variable,
					Pos:      pos,
				})
			}
			pos += len(matches[0])

			// float
		} else if matches := expPattern.floatNum.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[0], 64)
			exps = append(exps, TplExp{
				Type:  TplExpFloat,
				Float: num,
				Pos:   pos,
			})
			pos += len(matches[0])

			// int
		} else if matches := expPattern.intNum.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseInt(matches[0], 10, 64)
			exps = append(exps, TplExp{
				Type: TplExpInt,
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
				} else if exps[len(exps)-1].Type == TplExpOperator {
					lastOptr := exps[len(exps)-1].Operator
					if lastOptr != ")" && lastOptr != "]" {
						optr = "negative"
					}
				}
			}
			exps = append(exps, TplExp{
				Type:     TplExpOperator,
				Operator: optr,
				Pos:      pos,
			})
			pos += len(matches[0])

			// space
		} else if matches := expPattern.whitespace.FindStringSubmatch(left); len(matches) > 0 {
			pos += len(matches[0])

			// others
		} else {
			return exps, NewTplParseError(str, "exp.unexpectedToken", pos)
		}

		if pos > blen-1 {
			break
		}
	}

	return exps, nil
}

func generateTplExpTree(exps []TplExp) (*TplExp, error) {
	optrStack := make([]*TplExp, 0)
	opndStack := make([]*TplExp, 0)

	pos := 0
	for {
		if pos > len(exps)-1 {
			break
		}

		exp := exps[pos]

		fmt.Printf("################loop:%d###############\n", pos)
		spew.Dump(optrStack)
		spew.Dump(opndStack)
		spew.Dump(exp)

		// {...}
		if exp.Type == TplExpOperator && exp.Operator == "{" {
			bracketEnd := pos
			bracketNum := 1
			key := ""
			valBegin := 0
			valEnd := 0
			expMap := make(map[string]*TplExp)
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == TplExpOperator && exps[idx].Operator == "{" {
					bracketNum += 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == "}" {
					bracketNum -= 1
					valEnd = idx - 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == ":" {
					if exps[idx-1].Type != TplExpVar {
						return nil, errors.New("there must be a key before ':'")
					}
					if exps[idx-2].Type != TplExpOperator || (exps[idx-2].Operator != "{" && exps[idx-2].Operator != ";") {
						return nil, errors.New("invalid key before ':'")
					}
					key = exps[idx-1].Variable
					valBegin = idx + 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == ";" {
					valEnd = idx - 1
				}

				if bracketNum <= 1 && valEnd-valBegin >= 0 && key != "" {
					val, err := generateTplExpTree(exps[valBegin : valEnd+1])
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
				return nil, NewTplParseError("", "exp.mismatchedCurlyBrace", exp.Pos)
			}

			opndStack = append(opndStack, &TplExp{
				Type: TplExpMap,
				Map:  expMap,
			})
			pos = bracketEnd + 1

			// (...)
		} else if exp.Type == TplExpOperator && exp.Operator == "(" {
			bracketBegin := pos
			bracketEnd := pos
			bracketNum := 1
			for idx := pos + 1; idx < len(exps); idx++ {
				if (exps[idx].Type == TplExpOperator && exps[idx].Operator == "(") || exps[idx].Type == TplExpFunc {
					bracketNum += 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == ")" {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					bracketEnd = idx
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewTplParseError("", "exp.mismatchedParenthesis", exp.Pos)
			}
			parsedExp, err := generateTplExpTree(exps[bracketBegin+1 : bracketEnd])
			if err != nil {
				return nil, err
			}
			opndStack = append(opndStack, parsedExp)
			pos = bracketEnd + 1

			// [...]
		} else if exp.Type == TplExpOperator && exp.Operator == "[" {
			bracketBegin := pos
			bracketEnd := pos
			bracketNum := 1
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == TplExpOperator && exps[idx].Operator == "[" {
					bracketNum += 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == "]" {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					bracketEnd = idx
					break
				}
			}
			if bracketNum > 0 {
				return nil, NewTplParseError("", "exp.mismatchedSquareBracket", exp.Pos)
			}
			parsedExp, err := generateTplExpTree(exps[bracketBegin+1 : bracketEnd])
			if err != nil {
				return nil, err
			}
			lastOpnd := opndStack[len(opndStack)-1]
			opndStack[len(opndStack)-1] = &TplExp{
				Type:     TplExpCalc,
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
		} else if exp.Type == TplExpFunc {
			argBegin := pos + 1
			bracketNum := 1
			args := make([]*TplExp, 0)
			for idx := pos + 1; idx < len(exps); idx++ {
				if exps[idx].Type == TplExpOperator && exps[idx].Operator == "(" {
					bracketNum += 1
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == ")" {
					bracketNum -= 1
					if idx > argBegin {
						arg, err := generateTplExpTree(exps[argBegin:idx])
						if err != nil {
							return nil, err
						}
						args = append(args, arg)
					}
				} else if exps[idx].Type == TplExpOperator && exps[idx].Operator == "," && bracketNum == 1 {
					if idx == argBegin {
						return nil, NewTplParseError("", "exp.expectingParameter", exps[idx].Pos)
					}
					arg, err := generateTplExpTree(exps[argBegin:idx])
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
				return nil, NewTplParseError("", "exp.mismatchedParentheses", exp.Pos+len(exp.FuncName))
			}

			// operator
		} else if exp.Type == TplExpOperator {
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

		fmt.Printf("################reduce###############\n")
		spew.Dump(optrStack)
		spew.Dump(opndStack)

		var err error
		opndStack, optrStack, err = calculate(opndStack, optrStack)
		if err != nil {
			return nil, err
		}
	}

	return opndStack[0], nil
}

func calculate(opndStack []*TplExp, optrStack []*TplExp) ([]*TplExp, []*TplExp, error) {
	if len(optrStack) == 0 {
		return opndStack, optrStack, nil
	}

	optr := optrStack[len(optrStack)-1]
	switch optr.Operator {
	case "negative":
		if len(opndStack) == 0 {
			return opndStack, optrStack, NewTplParseError("", "exp.unexpectedToken", optr.Pos)
		}
		oexp := opndStack[len(opndStack)-1]
		nexp := &TplExp{
			Type:     TplExpCalc,
			Operator: "-",
			Left:     nil,
			Right:    oexp,
		}
		opndStack[len(opndStack)-1] = nexp
		optrStack = optrStack[:len(optrStack)-1]

	case "!":
		if len(opndStack) == 0 {
			return opndStack, optrStack, NewTplParseError("", "exp.unexpectedToken", optr.Pos)
		}
		oexp := opndStack[len(opndStack)-1]
		nexp := &TplExp{
			Type:     TplExpCalc,
			Operator: "!",
			Left:     nil,
			Right:    oexp,
		}
		opndStack[len(opndStack)-1] = nexp
		optrStack = optrStack[:len(optrStack)-1]

	case "?":
		if len(opndStack) < 2 {
			return opndStack, optrStack, NewTplParseError("", "exp.unexpectedToken", optr.Pos)
		}
		exp1 := opndStack[len(opndStack)-1]
		exp2 := opndStack[len(opndStack)-2]
		if exp1.Type != TplExpCalc || exp1.Operator != ":" {
			fmt.Println("-------------------------------")
			spew.Dump(optrStack)
			spew.Dump(opndStack)
			return opndStack, optrStack, NewTplParseError("", "exp.invalidTenaryExp", optr.Pos)
		}
		nexp := &TplExp{
			Type:            TplExpCalc,
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
			return opndStack, optrStack, NewTplParseError("", "exp.unexpectedToken", optr.Pos)
		}
		exp1 := opndStack[len(opndStack)-1]
		exp2 := opndStack[len(opndStack)-2]
		nexp := &TplExp{
			Type:     TplExpCalc,
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

func transformTenaryExp(exp *TplExp) error {
	if exp.Type != TplExpCalc || exp.Operator != "?" {
		return nil
	}
	if exp.Left == nil || exp.Right == nil || exp.Right.Operator != ":" {
		return NewTplParseError("", "exp.invalidTenaryExp", exp.Pos)
	}
	exp.TenaryCondition = exp.Left
	exp.Left = exp.Right.Left
	exp.Right = exp.Right.Right

	return nil
}

func ParseTplExp(str string) (*TplExp, error) {
	exps, err := readTplExp(str)
	if err != nil {
		return nil, err
	}
	tree, err2 := generateTplExpTree(exps)
	if err2 != nil {
		if tpe, ok := err2.(*TplParseError); ok {
			tpe.SetTpl(str)
		}
		return nil, err2
	}

	return tree, nil
}
