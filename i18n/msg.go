package i18n

var msgMap = map[string]string{
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

	"comp.unregisteredComponent":    "unregistered component: %s",
	"comp.cannotResolveComponent":   "failed to resolve component: %s",
	"comp.SetCompProp.propNotExist": "invalid property: %s",
	"comp.SetCompProp.typeMismatch": "cannot assign %s to %s",
}

func T(msg string) string {
	if cnt, ok := msgMap[msg]; ok {
		return cnt
	}

	return msg
}
