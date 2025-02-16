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
	"exp.incompleteExpression":          "incomplete expression",
	"exp.invalidTenaryExpression":       "invalid tenary expression",
	"exp.expectingParameter":            "expecting a parameter",

	"tpl.missingOpeningTag":             "missing opening tag",
	"tpl.missingClosingTag":             "missing closing tag",
	"tpl.incompleteTag":                 "incomplete tag",
	"tpl.mismatchedTag":                 "mismatched tag",
	"tpl.mismatchedSingleQuotationMark": "mismatched single quotation mark",
	"tpl.mismatchedDoubleQuotationMark": "mismatched double quotation mark",
	"tpl.duplicateAttribute":            "duplicate attribute",
	"tpl.duplicateDirective":            "duplicate directive",
	"tpl.conflictedDirective":           "conflicted directives",
	"tpl.duplicateEventHandler":         "duplicate event handler",
	"tpl.invalidForDirective":           "invalid v-for directive",
	"tpl.invalidDefAttr":                "invalid def attribute",
}

type DdlParseError struct {
	pos int
	ddl string
	err string
}

func NewDdlParseError(ddl string, err string, pos int) *DdlParseError {
	return &DdlParseError{
		ddl: ddl,
		err: err,
		pos: pos,
	}
}

func (dpe *DdlParseError) Error() string {
	msg := ""

	if err, ok := ErrorMap[dpe.err]; ok {
		msg = err + "\n"
	} else {
		msg = dpe.err + "\n"
	}

	if dpe.pos == -1 {
		dpe.pos = len(dpe.ddl) - 1
	}

	lines := strings.Split(dpe.ddl, "\n")
	lastLineEndPos := 0
	errRow := 0
	errCol := dpe.pos
	for idx := 0; idx < len(lines); idx++ {
		if lastLineEndPos+len(lines[idx]) >= dpe.pos {
			errRow = idx
			errCol = dpe.pos - lastLineEndPos
			break
		}
		lastLineEndPos += len(lines[idx]) + 1
	}

	beginLine := errRow - 2
	if beginLine < 0 {
		beginLine = 0
	}
	for idx := beginLine; idx <= errRow; idx++ {
		msg += lines[idx] + "\n"
	}
	for idx := 0; idx < errCol; idx++ {
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

func (dpe *DdlParseError) AddOffset(offset int) {
	dpe.pos += offset
}

func (dpe *DdlParseError) SetDdl(ddl string) {
	dpe.ddl = ddl
}

func (dpe *DdlParseError) SetPos(pos int) {
	dpe.pos = pos
}

func (dpe *DdlParseError) IsExpError() bool {
	return strings.HasPrefix(dpe.err, "exp")
}

func (dpe *DdlParseError) IsCssError() bool {
	return strings.HasPrefix(dpe.err, "css")
}
