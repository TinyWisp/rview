package template

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	varPattern      *regexp.Regexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*")
	intPattern      *regexp.Regexp = regexp.MustCompile("^[0-9]+")
	floatPattern    *regexp.Regexp = regexp.MustCompile(`^[0-9]+\.[0-9]+`)
	operatorPattern *regexp.Regexp = regexp.MustCompile(`^(\(|\)|\+|-|\*|/|%|==|!=|>=|<=|>|<|&&|\|\||!|,)`)
	funcPattern     *regexp.Regexp = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(`)

	operatorPriority = map[string]int{
		"!": 5,

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
	}
)

type TplExp struct {
	Type       TplExpType
	Int        int64
	Float      float64
	Str        string
	Bool       bool
	FuncName   string
	FuncParams []*TplExp
	Variable   string
	Operator   string
	Left       *TplExp
	Right      *TplExp
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
)

func readTplExp(str string) ([]TplExp, error) {
	exps := make([]TplExp, 0)

	byteArr := []byte(str)
	blen := len(byteArr)
	pos := 0

	for {
		ch := byteArr[pos]
		left := string(byteArr[pos:])

		// true
		if strings.HasPrefix(left, "true") {
			exps = append(exps, TplExp{
				Type: TplExpBool,
				Bool: true,
			})
			pos += 4

			// false
		} else if strings.HasPrefix(left, "false") {
			exps = append(exps, TplExp{
				Type: TplExpBool,
				Bool: false,
			})
			pos += 5

			// nil
		} else if strings.HasPrefix(left, "nil") {
			exps = append(exps, TplExp{
				Type: TplExpNil,
			})
			pos += 3

			// string literal
		} else if ch == '\'' {
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '\'' && byteArr[end-1] != '\\' {
					exps = append(exps, TplExp{
						Type: TplExpStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\'", "'"),
					})
					pos = end + 1
					break
				}
			}

			// string literal
		} else if ch == '"' {
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '"' && byteArr[end-1] != '\\' {
					exps = append(exps, TplExp{
						Type: TplExpStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\\"", "\""),
					})
					pos = end + 1
					break
				}
			}

			// func
		} else if matches := funcPattern.FindStringSubmatch(left); len(matches) > 0 {
			exps = append(exps, TplExp{
				Type:     TplExpFunc,
				FuncName: matches[1],
			})
			pos += len(matches[0])

			// variable
		} else if matches := varPattern.FindStringSubmatch(left); len(matches) > 0 {
			variable := matches[0]
			exps = append(exps, TplExp{
				Type:     TplExpVar,
				Variable: variable,
			})
			pos += len(matches[0])

			// float
		} else if matches := floatPattern.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[0], 64)
			exps = append(exps, TplExp{
				Type:  TplExpFloat,
				Float: num,
			})
			pos += len(matches[0])

			// int
		} else if matches := intPattern.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseInt(matches[0], 10, 64)
			exps = append(exps, TplExp{
				Type: TplExpInt,
				Int:  num,
			})
			pos += len(matches[0])

			// operator
		} else if matches := operatorPattern.FindStringSubmatch(left); len(matches) > 0 {
			exps = append(exps, TplExp{
				Type:     TplExpOperator,
				Operator: matches[0],
			})
			pos += len(matches[0])

			// space
		} else if ch == ' ' {
			pos += 1

			// others
		} else {
			return exps, fmt.Errorf("unexpected token: %s", left)
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

		// bracket
		if exp.Type == TplExpOperator && exp.Operator == "(" {
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
				return nil, errors.New("mismatched brackets")
			}
			parsedExp, err := generateTplExpTree(exps[bracketBegin+1 : bracketEnd])
			if err != nil {
				return nil, err
			}
			opndStack = append(opndStack, parsedExp)
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
						return nil, fmt.Errorf("invalid function parameters")
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
				return nil, fmt.Errorf(`no matching parentheses for "%s("`, exp.Left.FuncName)
			}

			// operator
		} else if exp.Type == TplExpOperator {
			if len(optrStack) == 0 {
				optrStack = append(optrStack, &exp)
				pos += 1
			} else {
				lastOptr := optrStack[len(optrStack)-1].Operator
				curOptr := exp.Operator
				if operatorPriority[curOptr] <= operatorPriority[lastOptr] && lastOptr != "!" {
					var exp1 *TplExp = nil
					exp2 := opndStack[len(opndStack)-1]
					if len(opndStack) >= 2 {
						exp1 = opndStack[len(opndStack)-2]
						opndStack = opndStack[:len(opndStack)-2]
					} else {
						opndStack = opndStack[:len(opndStack)-1]
					}
					optrStack = optrStack[:len(optrStack)-1]
					opndStack = append(opndStack, &TplExp{
						Type:     TplExpCalc,
						Operator: lastOptr,
						Left:     exp1,
						Right:    exp2,
					})
				} else if operatorPriority[curOptr] <= operatorPriority[lastOptr] && lastOptr == "!" {
					exp1 := opndStack[len(opndStack)-1]
					opndStack = opndStack[:len(opndStack)-1]
					optrStack = optrStack[:len(optrStack)-1]
					opndStack = append(opndStack, &TplExp{
						Type:     TplExpCalc,
						Operator: lastOptr,
						Left:     nil,
						Right:    exp1,
					})
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

	// the priority of the second operator is greater than the first one's
	// eg: var1 || var2 > 3
	if len(optrStack) == 2 {
		firstOptr := optrStack[0].Operator
		secondOptr := optrStack[1].Operator
		if operatorPriority[secondOptr] >= operatorPriority[firstOptr] && secondOptr != "!" {
			exp1 := opndStack[len(opndStack)-1]
			exp2 := opndStack[len(opndStack)-2]
			operator := secondOptr
			opndStack = opndStack[:len(opndStack)-2]
			opndStack = append(opndStack, &TplExp{
				Type:     TplExpCalc,
				Operator: operator,
				Left:     exp2,
				Right:    exp1,
			})
			optrStack = optrStack[:1]
		} else if operatorPriority[secondOptr] >= operatorPriority[firstOptr] && secondOptr == "!" {
			exp1 := opndStack[len(opndStack)-1]
			operator := secondOptr
			opndStack = opndStack[:len(opndStack)-1]
			opndStack = append(opndStack, &TplExp{
				Type:     TplExpCalc,
				Operator: operator,
				Left:     nil,
				Right:    exp1,
			})
			optrStack = optrStack[:1]
		}
	}

	if len(optrStack) == 1 && len(opndStack) == 2 {
		exp1 := opndStack[0]
		exp2 := opndStack[1]
		operator := optrStack[0].Operator
		opndStack = append(opndStack[:0], &TplExp{
			Type:     TplExpCalc,
			Operator: operator,
			Left:     exp1,
			Right:    exp2,
		})
	} else if len(optrStack) == 1 && len(opndStack) == 1 {
		var exp1 *TplExp = nil
		exp2 := opndStack[0]
		operator := optrStack[0].Operator
		opndStack = append(opndStack[:0], &TplExp{
			Type:     TplExpCalc,
			Operator: operator,
			Left:     exp1,
			Right:    exp2,
		})
	}

	return opndStack[0], nil
}

func ParseTplExp(str string) (*TplExp, error) {
	exps, err := readTplExp(str)
	if err != nil {
		return nil, err
	}
	tree, err2 := generateTplExpTree(exps)
	if err2 != nil {
		return nil, err2
	}

	return tree, nil
}
