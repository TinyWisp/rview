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

	"util.SetStructField.fieldNotExist":       "invalid field: %s",
	"util.SetStructField.typeMismatch":        "cannot assign %s to %s",
	"util.SetStructField.unexportedField":     "unexported field: %s",
	"util.SetStructField.cannotSetFieldValue": "cannot set the value of the field: %s",
	"util.GetStructField.fieldNotExist":       "invalid field: %s",
	"util.GetStructField.unexportedField":     "unexported field: %s",

	"calc.emptyVariableName":       "empty variable name",
	"calc.operandTypeMismatch":     "Type mismatch - cannot perform '%s' operation between '%s' and '%s'",
	"calc.invalidTernaryCondition": "The ternary condition must evaluate to a boolean. Received type '%s'",
	"calc.ternaryDataNotSameType":  `In a ternary expression "exp ? a : b", a and b must have the same data type.`,
	"calc.variableIsNotFunc":       `'%s' is not a function`,
	"calc.expMustBeVarType":        `the '%s' must be a *ddl.Exp with the ddl.ExpVar type`,
	"calc.expMustBeFuncType":       `the '%s' must be a *ddl.Exp with the ddl.ExpFunc type`,
	"calc.varIsNotFuncType":        `the '%s' is not a function`,
	"calc.argumentNumberMismatch":  `function "%s" expect %d arguments, but got %d`,
	"calc.argumentNumberNotEnough": `function "%s" expect %d or more arguments, but got %d`,

	"comp.compNotFound":               "component not found: '%s' is not recognized. check if it is registered or spelled correctly.",
	"comp.propNotAllowed":             "invalid property: '%s' is not allowed on '%s",
	"comp.propTypeMismatch":           "invalid property: cannot assign %s to %s on %s",
	"comp.unregisteredutilonent":      "unregistered utilonent: %s",
	"comp.cannotResolveComponent":     "failed to resolve component: %s",
	"comp.SetCompProp.propNotExist":   "invalid property: %s",
	"comp.SetCompProp.typeMismatch":   "cannot assign %s to %s",
	"comp.vifDirectiveMustBeBool":     "Invalid expression in v-if: expected a boolean, got %s instead",
	"comp.velseifDirectiveMustBeBool": "Invalid expression in v-else-if: expected a boolean, got %s instead",
	"comp.velseDirectiveMustBeBool":   "Invalid expression in v-else: expected a boolean, got %s instead",
}

func T(msg string) string {
	if cnt, ok := msgMap[msg]; ok {
		return cnt
	}

	return msg
}
