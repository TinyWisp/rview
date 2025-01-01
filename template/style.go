package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CSSTokenType int

const (
	CSSTokenNum CSSTokenType = iota
	CSSTokenStr
	CSSTokenOperator
	CSSTokenFunc
)

type CSSUnit int

const (
	ch  CSSUnit = iota
	vw          // viewport's width
	vh          // viewport's height
	pfw         // parent's full width
	pfh         // parent's full height
	paw         // parent's available width
	pah         // parent's available height
	noUnit
)

var cssUnitMap = map[string]CSSUnit{
	"ch":  ch,
	"vw":  vw,
	"vh":  vh,
	"pfw": pfw,
	"pfh": pfh,
	"paw": paw,
	"pah": pah,
	"":    noUnit,
}

type CSSToken struct {
	Type      CSSTokenType
	Num       float64
	Unit      CSSUnit
	Str       string
	Variable  string
	Operator  rune
	FuncName  string
	Arguments [][]CSSToken
}

var (
	cssPattern = struct {
		num      *regexp.Regexp
		funct    *regexp.Regexp
		operator *regexp.Regexp
		variable *regexp.Regexp
		str      *regexp.Regexp
	}{
		num:      regexp.MustCompile("^([0-9]+[.]?[0-9]*)(ch|vw|vh|pfw|pfh|paw|pah)?"),
		funct:    regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(`),
		operator: regexp.MustCompile(`^[+\-*/()]`),
		str:      regexp.MustCompile(`^\S+`),
	}
)

func ParseCssRule(rule string) ([]CSSToken, error) {
	tokens := make([]CSSToken, 0, 20)

	byteArr := []byte(rule)
	blen := len(byteArr)
	begin := 0

	for {
		ch := byteArr[begin]
		left := string(byteArr[begin:])

		// string literal
		if ch == '\'' {
			for end := begin; end < blen; end++ {
				if byteArr[end] == '\'' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						Type: CSSTokenStr,
						Str:  strings.ReplaceAll(string(byteArr[begin:end+1]), "\\'", "'"),
					})
					begin = end + 1
					break
				}
			}

			// func
		} else if matches := cssPattern.funct.FindStringSubmatch(left); len(matches) > 0 {
			funcName := matches[1]
			bracketNum := 1
			for end := begin + len(funcName) + 1; end < blen; end++ {
				if byteArr[end] == '(' {
					bracketNum += 1
				} else if byteArr[end] == ')' {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					argBegin := begin + len(funcName) + 1
					argEnd := end - 1
					argTokens, err := ParseCssRule(string(byteArr[argBegin : argEnd+1]))
					if err != nil {
						return tokens, err
					}
					args := getFuncArguments(argTokens)

					tokens = append(tokens, CSSToken{
						Type:      CSSTokenFunc,
						FuncName:  funcName,
						Arguments: args,
					})
					begin = end + 1
					break
				}
			}

			if bracketNum > 0 {
				return tokens, fmt.Errorf("unmatched bracket")
			}

			// num
		} else if matches := cssPattern.num.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[1], 64)
			unitStr := matches[2]
			unit, ok := cssUnitMap[unitStr]
			if !ok {
				return tokens, fmt.Errorf("unexpected token: %s", left)
			}

			tokens = append(tokens, CSSToken{
				Type: CSSTokenNum,
				Num:  num,
				Unit: unit,
			})
			begin += len(matches[0])

			// str
		} else if matches := cssPattern.str.FindStringSubmatch(left); len(matches) > 0 {
			tokens = append(tokens, CSSToken{
				Type: CSSTokenStr,
				Str:  matches[0],
			})
			begin += len(matches[0])

			// space
		} else if ch == ' ' {
			begin += 1

			// others
		} else {
			return tokens, fmt.Errorf("unexpected token: %s", left)
		}

		if begin > blen-1 {
			break
		}
	}

	return tokens, nil
}

func getFuncArguments(tokens []CSSToken) [][]CSSToken {
	args := make([][]CSSToken, 0, 10)

	arg := make([]CSSToken, 0, 10)
	for _, token := range tokens {
		if token.Type == CSSTokenOperator && token.Operator == ',' {
			args = append(args, arg)
			arg = arg[:0]
		} else {
			arg = append(arg, token)
		}
	}
	args = append(args, arg)

	return args
}
