package ddl

import (
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
	CSSTokenProp

	CSSTokenAuto
	CSSTokenHidden
	CSSTokenVisible
	CSSTokenScroll
)

type CSSUnit int

const (
	NoUnit CSSUnit = iota
	CH             // 1 character
	VW             // relative to 1% of the width of the viewport
	VH             // relative to 1% of the height of the viewport
	PFW            // relative to 1% of the full width of the parent
	PFH            // relative to 1% of the full height of the parent
	PCW            // relative to 1% of the width of the parent's content
	PCH            // relative to 1% of the height of the parent's content
)

var cssUnitMap = map[string]CSSUnit{
	"ch":  CH,
	"vw":  VW,
	"vh":  VH,
	"pfw": PFW,
	"pfh": PFH,
	"pcw": PCW,
	"pch": PCH,
	"":    NoUnit,
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
	IntColor  int32
	Prop      string
	Arguments [][]CSSToken
	Pos       int
}

type CSSPropVal []CSSToken

type CSSClass map[string][]CSSToken

type CSSClassMap map[string]CSSClass

var (
	cssPattern = struct {
		num        *regexp.Regexp
		funct      *regexp.Regexp
		operator   *regexp.Regexp
		variable   *regexp.Regexp
		color      *regexp.Regexp
		str        *regexp.Regexp
		class      *regexp.Regexp
		prop       *regexp.Regexp
		whitespace *regexp.Regexp
		keyword    *regexp.Regexp
	}{
		num:        regexp.MustCompile("^([0-9]+[.]?[0-9]*)(ch|vw|vh|pfw|pfh|pcw|pch)?"),
		color:      regexp.MustCompile("^#[0-9a-fA-F]{6}"),
		funct:      regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(`),
		operator:   regexp.MustCompile(`^[+\-*/(){}:;]`),
		prop:       regexp.MustCompile(`^(;|\{)(\s*)([a-zA-Z0-9_\-]+)`),
		str:        regexp.MustCompile(`^[^\s:;.{}()+\-*/]+`),
		class:      regexp.MustCompile(`^\.([0-9a-zA-Z_\-]+)`),
		variable:   regexp.MustCompile(`^var\(.*?\)`),
		whitespace: regexp.MustCompile(`^\s+`),
		keyword:    regexp.MustCompile(`^(auto|scroll|hidden)`),
	}

	propRuleMap = map[string][][]CSSTokenType{
		"margin":              {{CSSTokenNum}, {CSSTokenNum, CSSTokenNum}, {CSSTokenNum, CSSTokenNum, CSSTokenNum, CSSTokenNum}},
		"margin-top":          {{CSSTokenNum}},
		"margin-bottom":       {{CSSTokenNum}},
		"margin-left":         {{CSSTokenNum}},
		"margin-right":        {{CSSTokenNum}},
		"padding":             {{CSSTokenNum}, {CSSTokenNum, CSSTokenNum}, {CSSTokenNum, CSSTokenNum, CSSTokenNum, CSSTokenNum}},
		"padding-top":         {{CSSTokenNum}},
		"padding-bottom":      {{CSSTokenNum}},
		"padding-left":        {{CSSTokenNum}},
		"padding-right":       {{CSSTokenNum}},
		"border-width":        {{CSSTokenNum}, {CSSTokenNum, CSSTokenNum}, {CSSTokenNum, CSSTokenNum, CSSTokenNum, CSSTokenNum}},
		"border-left-width":   {{CSSTokenNum}},
		"border-right-width":  {{CSSTokenNum}},
		"border-top-width":    {{CSSTokenNum}},
		"border-bottom-width": {{CSSTokenNum}},
		"border-color":        {{CSSTokenColor}, {CSSTokenColor, CSSTokenColor}, {CSSTokenColor, CSSTokenColor, CSSTokenColor, CSSTokenColor}},
		"border-left-color":   {{CSSTokenColor}},
		"border-right-color":  {{CSSTokenColor}},
		"border-top-color":    {{CSSTokenColor}},
		"border-bottom-color": {{CSSTokenColor}},
		"border-char":         {{CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr, CSSTokenStr}},
		"background-color":    {{CSSTokenColor}},
		"overflow":            {{CSSTokenAuto}, {CSSTokenHidden}, {CSSTokenScroll}},
		"overflow-x":          {{CSSTokenAuto}, {CSSTokenHidden}, {CSSTokenScroll}},
		"overflow-y":          {{CSSTokenAuto}, {CSSTokenHidden}, {CSSTokenScroll}},
	}

	keywordMap = map[string]CSSTokenType{
		"auto":   CSSTokenAuto,
		"scroll": CSSTokenScroll,
		"hidden": CSSTokenHidden,
	}
)

func parseCss(css string) (CSSClassMap, error) {
	tokens, err := tokenizeCss(css)
	if err != nil {
		return nil, err
	}

	classMap, err2 := genCssClassMap(tokens)
	if err2 != nil {
		if tpe, ok := err2.(*DdlParseError); ok {
			tpe.SetDdl(css)
		}
		return nil, err2
	}

	err3 := checkCssPropRule(classMap)
	if err3 != nil {
		if tpe, ok := err3.(*DdlParseError); ok {
			tpe.SetDdl(css)
		}
		return nil, err3
	}

	coverCssProp(classMap)

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

		} else if expect == "key" && token.Type == CSSTokenProp {
			key = token.Prop
			if _, ok := propRuleMap[key]; !ok && !strings.HasPrefix(key, "--") {
				return classMap, NewDdlParseError("", "css.unsupportedProp", token.Pos)
			}
			expect = ":"

		} else if expect == ":" && token.Type == CSSTokenOperator && token.Operator == ":" {
			expect = "val"

		} else if expect == "val" && (token.Type != CSSTokenOperator || (token.Operator != "}" && token.Operator != ";")) {
			val = append(val, token)
			expect = "val|;"

		} else if expect == "val|;" && (token.Type != CSSTokenOperator || (token.Operator != "}" && token.Operator != ";")) {
			val = append(val, token)
			expect = "val|;"

		} else if expect == "val|;" && token.Type == CSSTokenOperator && token.Operator == ";" {
			classMap[className][key] = val
			expect = "key|}"
			val = []CSSToken{}

		} else if expect == "key|}" && token.Type == CSSTokenProp {
			key = token.Prop
			if _, ok := propRuleMap[key]; !ok && !strings.HasPrefix(key, "--") {
				return classMap, NewDdlParseError("", "css.unsupportedProp", token.Pos)
			}
			expect = ":"

		} else if expect == "key|}" && token.Type == CSSTokenOperator && token.Operator == "}" {
			expect = "class"

		} else {
			return classMap, NewDdlParseError("", "css.unexpectedToken", token.Pos)
		}
	}

	if expect != "class" {
		return classMap, NewDdlParseError("", "css.unexpectedEnd", -1)
	}

	return classMap, nil
}

func checkCssPropRule(classMap CSSClassMap) error {
	for _, cpropMap := range classMap {
		for pkey, pval := range cpropMap {
			if strings.HasPrefix(pkey, "--") {
				continue
			}
			rules := propRuleMap[pkey]
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
			if !valid {
				return NewDdlParseError("", "css.invalidPropVal", pval[0].Pos)
			}
		}
	}

	return nil
}

func coverCssProp(classMap CSSClassMap) {
	// ----- margin, padding, border-width, border-color -------
	var tprops = []string{"margin", "padding", "border-width", "border-color"}
	var pattern = regexp.MustCompile("^([a-z]+)")

	for _, cpropMap := range classMap {
		for _, tprop := range tprops {
			if tokens, ok := cpropMap[tprop]; ok {
				var left, right, top, bottom CSSToken

				if len(tokens) == 1 {
					left = tokens[0]
					right = left
					top = left
					bottom = left
				} else if len(tokens) == 2 {
					top = tokens[0]
					bottom = tokens[0]
					left = tokens[1]
					right = left
				} else if len(tokens) == 4 {
					top = tokens[0]
					right = tokens[1]
					bottom = tokens[2]
					left = tokens[3]
				}

				tleftProp := pattern.ReplaceAllString(tprop, "$1-left")
				if _, ok := cpropMap[tleftProp]; !ok {
					cpropMap[tleftProp] = []CSSToken{left}
				}
				trightProp := pattern.ReplaceAllString(tprop, "$1-right")
				if _, ok := cpropMap[trightProp]; !ok {
					cpropMap[trightProp] = []CSSToken{right}
				}
				ttopProp := pattern.ReplaceAllString(tprop, "$1-top")
				if _, ok := cpropMap[ttopProp]; !ok {
					cpropMap[ttopProp] = []CSSToken{top}
				}
				tbottomProp := pattern.ReplaceAllString(tprop, "$1-bottom")
				if _, ok := cpropMap[tbottomProp]; !ok {
					cpropMap[tbottomProp] = []CSSToken{bottom}
				}
			}
		}
	}

	// --------------------- ovreflow ---------------------------
	for _, cpropMap := range classMap {
		if overflow, ok := cpropMap["overflow"]; ok {
			if _, ok2 := cpropMap["overflow-x"]; !ok2 {
				cpropMap["overflow-x"] = overflow
			}
			if _, ok2 := cpropMap["overflow-y"]; !ok2 {
				cpropMap["overflow-y"] = overflow
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
			match := false
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '\'' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						Type: CSSTokenStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\'", "'"),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return tokens, NewDdlParseError(css, "css.mismatchedSingleQuotationMark", pos)
			}

			// string literal
		} else if ch == '"' {
			match := false
			for end := pos + 1; end < blen; end++ {
				if byteArr[end] == '"' && byteArr[end-1] != '\\' {
					tokens = append(tokens, CSSToken{
						Type: CSSTokenStr,
						Str:  strings.ReplaceAll(string(byteArr[pos+1:end]), "\\\"", "\""),
						Pos:  pos,
					})
					pos = end + 1
					match = true
					break
				}
			}
			if !match {
				return tokens, NewDdlParseError(css, "css.mismatchedDoubleQuotationMark", pos)
			}

			// color
		} else if matches := cssPattern.color.FindStringSubmatch(left); len(matches) > 0 {
			color := matches[0]
			hexColor := strings.TrimPrefix(color, "#")
			intColor, _ := strconv.ParseInt(hexColor, 16, 32)
			tokens = append(tokens, CSSToken{
				Type:     CSSTokenColor,
				Color:    color,
				IntColor: int32(intColor),
				Pos:      pos,
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

			// prop
		} else if matches := cssPattern.prop.FindStringSubmatch(left); len(matches) > 0 {
			operator := matches[1]
			white := matches[2]
			prop := matches[3]
			tokens = append(tokens, CSSToken{
				Type:     CSSTokenOperator,
				Operator: operator,
				Pos:      pos,
			}, CSSToken{
				Type: CSSTokenProp,
				Prop: prop,
				Pos:  pos + len(operator) + len(white),
			})
			pos += len(matches[0])

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
				return tokens, NewDdlParseError(css, "css.mismatchedParenthesis", pos+len(funcName))
			}

			// num
		} else if matches := cssPattern.num.FindStringSubmatch(left); len(matches) > 0 {
			num, _ := strconv.ParseFloat(matches[1], 64)
			unitStr := matches[2]
			unit, ok := cssUnitMap[unitStr]
			if !ok {
				return tokens, NewDdlParseError(css, "css.unexpectedToken", pos)
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

			// keyword: auto, scroll, hidden
		} else if matches := cssPattern.keyword.FindStringSubmatch(left); len(matches) > 0 {
			word := matches[0]
			tokens = append(tokens, CSSToken{
				Type: keywordMap[word],
			})
			pos += len(word)

			// str
		} else if matches := cssPattern.str.FindStringSubmatch(left); len(matches) > 0 {
			tokens = append(tokens, CSSToken{
				Type: CSSTokenStr,
				Str:  matches[0],
				Pos:  pos,
			})
			pos += len(matches[0])

			// whitespace characters
		} else if matches := cssPattern.whitespace.FindStringSubmatch(left); len(matches) > 0 {
			pos += len(matches[0])

			// others
		} else {
			return tokens, NewDdlParseError(css, "css.unexpectedToken", pos)
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
