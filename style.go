package rview

import (
	"regexp"
	"strconv"
	"strings"
	"errors"
	"fmt"
)

type CSSTokenType int

const (
  CSSTokenNum CSSTokenType = iota
	CSSTokenStr
	CSSTokenOperator
	CSSTokenFunc
	CSSTokenVar
)

type CSSUnit int

const (
	ch CSSUnit = iota
	vw
	vh
	pfw
	pfh
	pcw
	pch
)

type CSSToken struct {
	org string
	kind CSSTokenType
	num float64
	unit string
	str string
	variable string
	operator rune
	funcName string
	arguments [][]CSSToken
}

var (
	numPattern *regexp.Regexp
	funcPattern *regexp.Regexp
	operatorPattern *regexp.Regexp
	variablePattern *regexp.Regexp
	strPattern *regexp.Regexp
)

func init() {
	numPattern = regexp.MustCompile("^([0-9]+[.]?[0-9]*)(ch|vw|vh|pfw|pfh|pcw|pch)?")
	funcPattern = regexp.MustCompile("^([a-zA-Z_][a-zA-Z0-9_]*)\\(")
	operatorPattern = regexp.MustCompile("^[+\\-*/()]")
	variablePattern = regexp.MustCompile("^--[a-zA-Z0-9]+")
	strPattern = regexp.MustCompile("^\\S+")
}

func ParseRule(rule string) ([]CSSToken, error) {
	tokens := make([]CSSToken, 0, 20)

	byteArr := []byte(rule)
	blen := len(byteArr)
	begin := 0

	for {
		ch := byteArr[begin]
		left := string(byteArr[begin:])

		// string literal
		if ch == '\'' {
			for end:=begin; end<blen; end++ {
				if byteArr[end] == '\'' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						kind: CSSTokenStr,
						str: strings.ReplaceAll(string(byteArr[begin:end+1]), "\\'", "'"),
					})
					begin = end + 1
					break
				}
			}
		
		// func
		} else if matches := funcPattern.FindStringSubmatch(left); len(matches) > 0 {
			funcName := matches[1]
			bracketNum := 1
			for end:=begin + len(funcName); end<blen; end++ {
				if byteArr[end] == '(' {
					bracketNum += 1
				} else if byteArr[end] == ')' {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					argBegin := begin + len(funcName)
					argEnd := end - 1
					argTokens, err := ParseRule(string(byteArr[argBegin: argEnd+1]))
					if err != nil {
						return tokens, err
					}
					args := getFuncArguments(argTokens)

					tokens = append(tokens, CSSToken{
						kind: CSSTokenFunc,
						funcName: funcName,
						arguments: args,
					})
					begin = end + 1
					break
				}
			}

		// variable
		} else if matches := variablePattern.FindStringSubmatch(left); len(matches) > 0 {
			variable := matches[1]
			tokens = append(tokens, CSSToken{
				kind: CSSTokenVar,
				variable: variable,
			})
			begin += len(matches[0])

		// num
		} else if matches := numPattern.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[1], 64)
			unit := matches[2]

			tokens = append(tokens, CSSToken{
				kind: CSSTokenNum,
				num: num,
				unit: unit,
			})
			begin += len(matches[0])

		// str
		} else if matches := strPattern.FindStringSubmatch(left); len(matches) > 0 {
			tokens = append(tokens, CSSToken{
				kind: CSSTokenStr,
				str: matches[0],
			})
			begin += len(matches[0])

		// space
		} else if ch == ' ' {
			begin += 1

		// others
		} else {
			return tokens, errors.New(fmt.Sprintf("unexpected token: %s", left))
		}

		if begin > blen - 1 {
			break
		}
	}

	return tokens, nil
}

func getFuncArguments(tokens []CSSToken) [][]CSSToken {
	args := make([][]CSSToken, 0, 10)

	arg := make([]CSSToken, 0, 10)
	for _, token := range tokens {
		if token.kind == CSSTokenOperator && token.operator == ',' {
			args = append(args, )
		} else {
			arg = append(arg, token)
		}
	}
	args = append(args, arg)

	return args
}