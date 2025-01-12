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
	CSSTokenColor
	CSSTokenOperator
	CSSTokenFunc
	CSSTokenVar
	CSSTokenClass
)

type CSSUnit int

const (
	noUnit CSSUnit = iota
	ch             // 1 character
	vw             // relative to 1% of the width of the viewport
	vh             // relative to 1% of the height of the viewport
	pfw            // relative to 1% of the full width of the parent
	pfh            // relative to 1% of the full height of the parent
	pcw            // relative to 1% of the width of the parent's content
	pch            // relative to 1% of the height of the parent's content
)

var cssUnitMap = map[string]CSSUnit{
	"ch":  ch,
	"vw":  vw,
	"vh":  vh,
	"pfw": pfw,
	"pfh": pfh,
	"pcw": pcw,
	"pch": pch,
	"":    noUnit,
}

type CSSToken struct {
	Type      CSSTokenType
	Num       float64
	Unit      CSSUnit
	Str       string
	Variable  string
	Operator  string
	FuncName  string
	Class     string
	Color     string
	Arguments [][]CSSToken
	Pos       int
}

type CSSPropVal []CSSToken

type CSSClass map[string][]CSSToken

type CSSClassMap map[string]CSSClass

var (
	cssPattern = struct {
		num      *regexp.Regexp
		funct    *regexp.Regexp
		operator *regexp.Regexp
		variable *regexp.Regexp
		color    *regexp.Regexp
		str      *regexp.Regexp
		class    *regexp.Regexp
		prop     *regexp.Regexp
	}{
		num:      regexp.MustCompile("^([0-9]+[.]?[0-9]*)(ch|vw|vh|pfw|pfh|pcw|pch)?"),
		color:    regexp.MustCompile("^#[0-9a-fA-F]{6}"),
		funct:    regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(`),
		operator: regexp.MustCompile(`^[+\-*/(){}]`),
		str:      regexp.MustCompile(`^[^\s:;.{}()+\-*/]+`),
		class:    regexp.MustCompile(`^\.([0-9a-zA-Z_\-]+)`),
		variable: regexp.MustCompile(`^var\(.*?\)`),
	}

	propRuleMap = map[string][][]CSSTokenType{
		"margin":         {{CSSTokenNum}, {CSSTokenNum, CSSTokenNum}, {CSSTokenNum, CSSTokenNum, CSSTokenNum, CSSTokenNum}},
		"margin-top":     {{CSSTokenNum}},
		"margin-bottom":  {{CSSTokenNum}},
		"margin-left":    {{CSSTokenNum}},
		"margin-right":   {{CSSTokenNum}},
		"padding":        {{CSSTokenNum}, {CSSTokenNum, CSSTokenNum}, {CSSTokenNum, CSSTokenNum, CSSTokenNum, CSSTokenNum}},
		"padding-top":    {{CSSTokenNum}},
		"padding-bottom": {{CSSTokenNum}},
		"padding-left":   {{CSSTokenNum}},
		"padding-right":  {{CSSTokenNum}},
		"border":         {{CSSTokenNum}, {CSSTokenNum, CSSTokenColor}, {CSSTokenNum, CSSTokenColor, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr}},
		"border-color":   {{CSSTokenColor}},
		"border-ch":      {{CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr}},
		"border-ch-lt":   {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
		"border-ch-rt":   {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
		"border-ch-lb":   {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
		"border-ch-rb":   {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
		"border-ch-v":    {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
		"border-ch-h":    {{CSSTokenStr}, {CSSTokenStr, CSSTokenColor}},
	}
)

func parseCss(css string) (CSSClassMap, error) {
	tokens, err := tokenizeCss(css)
	if err != nil {
		return nil, err
	}

	classMap, err2 := genCssClassMap(tokens)
	if err2 != nil {
		return nil, err2
	}

	return classMap, nil
}

func genCssClassMap(tokens []CSSToken) (CSSClassMap, error) {
	classMap := make(CSSClassMap)

	className := ""
	expect := "class"
	key := ""
	val := []CSSToken{}
	for _, token := range tokens {
		if expect == "class" && token.Type == CSSTokenClass {
			className = token.Class
			classMap[className] = make(CSSClass)
			expect = "{"

		} else if expect == "{" && token.Type == CSSTokenOperator && token.Operator == "{" {
			expect = "key"

		} else if expect == "key" && token.Type == CSSTokenStr {
			key = token.Str
			expect = ":"

		} else if expect == ":" && token.Type == CSSTokenOperator && token.Operator == ":" {
			expect = "val"

		} else if expect == "val" && (token.Type != CSSTokenOperator || (token.Operator != "}" && token.Operator != ";")) {
			val = append(val, token)
			expect = "val|;|}"

		} else if expect == "val|;" && (token.Type != CSSTokenOperator || (token.Operator != "}" && token.Operator != ";")) {
			val = append(val, token)
			expect = "val|;|}"

		} else if expect == "val|;" && token.Type == CSSTokenOperator && token.Operator == ";" {
			classMap[className][key] = val
			expect = "key|}"

		} else if expect == "key|}" && token.Type == CSSTokenStr {
			key = token.Str
			expect = ":"

		} else if expect == "key|}" && token.Type == CSSTokenOperator && token.Operator == "}" {
			expect = "class"

		} else {
			return classMap, fmt.Errorf("unexpected")
		}
	}

	return classMap, nil
}

func checkCssPropRule(classMap CSSClassMap) error {
	for _, cpropMap := range classMap {
		for pkey, pval := range cpropMap {
			if strings.HasPrefix(pkey, "--") {
				continue
			}
			rules, ok := propRuleMap[pkey]
			if !ok {
				return fmt.Errorf("invalid prop: %s", pkey)
			}
			valid := false
			for _, rule := range rules {
				if len(pval) != len(rule) {
					continue
				}
				for tidx, ruleTokenType := range rule {
					if pval[tidx].Type != ruleTokenType {
						break
					}
					if tidx == len(rule)-1 {
						valid = true
						break
					}
				}
				if valid {
					break
				}
			}
		}
	}

	return nil
}

func transformCssProp(classMap CSSClassMap) {
	for cname, cpropMap := range classMap {
		// margin
		if marginTokens, ok := cpropMap["margin"]; ok {
			if len(marginTokens) == 1 {
				mtoken := marginTokens[0]
				newMarginTokens := CSSPropVal{
					mtoken, mtoken, mtoken, mtoken,
				}
				classMap[cname]["margin"] = newMarginTokens
			} else if len(marginTokens) == 2 {
				mtokenTB := marginTokens[0]
				mtokenLR := marginTokens[1]
				newMarginTokens := CSSPropVal{
					mtokenTB, mtokenLR, mtokenTB, mtokenLR,
				}
				classMap[cname]["margin"] = newMarginTokens
			}

			// padding
		} else if paddingTokens, ok := cpropMap["padding"]; ok {
			if len(paddingTokens) == 1 {
				ptoken := paddingTokens[0]
				newPaddingTokens := CSSPropVal{
					ptoken, ptoken, ptoken, ptoken,
				}
				classMap[cname]["padding"] = newPaddingTokens
			} else if len(paddingTokens) == 2 {
				ptokenTB := paddingTokens[0]
				ptokenLR := paddingTokens[1]
				newPaddingTokens := CSSPropVal{
					ptokenTB, ptokenLR, ptokenTB, ptokenLR,
				}
				classMap[cname]["padding"] = newPaddingTokens
			}

		}
	}
}

func tokenizeCss(css string) ([]CSSToken, error) {
	tokens := make([]CSSToken, 0, 20)

	byteArr := []byte(css)
	blen := len(byteArr)
	pos := 0

	for {
		ch := byteArr[pos]
		left := string(byteArr[pos:])

		// string literal
		if ch == '\'' {
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '\'' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						Type: CSSTokenStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\'", "'"),
						Pos:  pos,
					})
					pos = end + 1
					break
				}
			}

			// string literal
		} else if ch == '"' {
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '"' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						Type: CSSTokenStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\\"", "\""),
						Pos:  pos,
					})
					pos = end + 1
					break
				}
			}

			// color
		} else if matches := cssPattern.color.FindStringSubmatch(left); len(matches) > 0 {
			color := matches[0]
			tokens = append(tokens, CSSToken{
				Type:  CSSTokenColor,
				Color: color,
				Pos:   pos,
			})
			pos += len(color)

			// class name
		} else if matches := cssPattern.class.FindStringSubmatch(left); len(matches) > 0 {
			class := matches[1]
			tokens = append(tokens, CSSToken{
				Type:  CSSTokenClass,
				Class: class,
				Pos:   pos,
			})
			pos += len(class) + 1

			// variable
		} else if matches := cssPattern.variable.FindStringSubmatch(left); len(matches) > 0 {
			variable := matches[1]
			tokens = append(tokens, CSSToken{
				Type:     CSSTokenVar,
				Variable: variable,
				Pos:      pos,
			})
			pos += len(variable) + 5

			// func
		} else if matches := cssPattern.funct.FindStringSubmatch(left); len(matches) > 0 {
			funcName := matches[1]
			bracketNum := 1
			for end := pos + len(funcName) + 1; end < blen; end++ {
				if byteArr[end] == '(' {
					bracketNum += 1
				} else if byteArr[end] == ')' {
					bracketNum -= 1
				}
				if bracketNum == 0 {
					argBegin := pos + len(funcName) + 1
					argEnd := end - 1
					argTokens, err := tokenizeCss(string(byteArr[argBegin : argEnd+1]))
					if err != nil {
						return tokens, err
					}
					args := getFuncArguments(argTokens)

					tokens = append(tokens, CSSToken{
						Type:      CSSTokenFunc,
						FuncName:  funcName,
						Arguments: args,
						Pos:       pos,
					})
					pos = end + 1
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
				Pos:  pos,
			})
			pos += len(matches[0])

			// operator
		} else if matches := cssPattern.operator.FindStringSubmatch(left); len(matches) > 0 {
			operator := matches[0]
			tokens = append(tokens, CSSToken{
				Type:     CSSTokenOperator,
				Operator: operator,
				Pos:      pos,
			})
			pos += len(operator)

			// str
		} else if matches := cssPattern.str.FindStringSubmatch(left); len(matches) > 0 {
			tokens = append(tokens, CSSToken{
				Type: CSSTokenStr,
				Str:  matches[0],
				Pos:  pos,
			})
			pos += len(matches[0])

			// space
		} else if ch == ' ' {
			pos += 1

			// others
		} else {
			return tokens, fmt.Errorf("unexpected token: %s", left)
		}

		if pos > blen-1 {
			break
		}
	}

	return tokens, nil
}

func getFuncArguments(tokens []CSSToken) [][]CSSToken {
	args := make([][]CSSToken, 0, 10)

	arg := make([]CSSToken, 0, 10)
	for _, token := range tokens {
		if token.Type == CSSTokenOperator && token.Operator == "," {
			args = append(args, arg)
			arg = arg[:0]
		} else {
			arg = append(arg, token)
		}
	}
	args = append(args, arg)

	return args
}
