package template

import (
	"fmt"
	"strings"
)

var ErrorMap = map[string]string{
	"css.mismatchedCurlyBrace":          "mismatched brace",
	"css.unexpectedToken":               "unexpected token",
	"css.unexpectedEnd":                 "unexpected end",
	"css.unsupportedProp":               "unsupported property",
	"css.invalidPropVal":                "invalid property value",
	"css.mismatchedSingleQuotationMark": "mismatched single quotation mark",
	"css.mismatchedDoubleQuotationMark": "mismatched double quotation mark",
	"css.mismatchedParenthesis":         "mismatched parenthesis",

	"exp.mismatchedCurlyBrace":          "mismatched brace",
	"exp.mismatchedSingleQuotationMark": "mismatched single quotation mark",
	"exp.mismatchedDoubleQuotationMark": "mismatched double quotation mark",
	"exp.mismatchedParenthesis":         "mismatched parenthesis",
	"exp.unexpectedToken":               "unexpected token",
}

type TplParseError struct {
	pos int
	tpl string
	err string
}

func NewTplParseError(tpl string, err string, pos int) *TplParseError {
	return &TplParseError{
		tpl: tpl,
		err: err,
		pos: pos,
	}
}

func (tpe *TplParseError) Error() string {
	msg := ""

	if err, ok := ErrorMap[tpe.err]; ok {
		msg = err + "\n"
	} else {
		msg = tpe.err + "\n"
	}

	if tpe.pos == -1 {
		tpe.pos = len(tpe.tpl) - 1
	}
	fmt.Println(tpe.pos)

	lines := strings.Split(tpe.tpl, "\n")
	lineEndPos := 0
	errRow := 0
	errCol := 0
	for idx := 0; idx < len(lines); idx++ {
		lineEndPos += len(lines[idx]) + 1
		if lineEndPos >= tpe.pos {
			errRow = idx
			errCol = tpe.pos - lineEndPos
			break
		}
	}

	beginLine := errRow - 2
	if beginLine < 0 {
		beginLine = 0
	}
	for idx := beginLine; idx <= errRow; idx++ {
		msg += lines[idx] + "\n"
	}
	for idx := 0; idx < errCol-1; idx++ {
		msg += " "
	}
	msg += "^\n"
	for idx := errRow + 1; idx < errRow+2 && idx < len(lines); idx++ {
		msg += lines[idx] + "\n"
	}

	fmt.Printf("errRow:%d, errCol:%d, beginLine:%d, endLine:%d\n", errRow, errCol, beginLine, len(lines)-1)
	fmt.Println(lines)

	return msg
}

func (tpe *TplParseError) AddOffset(offset int) {
	tpe.pos += offset
}

func (tpe *TplParseError) SetTpl(tpl string) {
	tpe.tpl = tpl
}
