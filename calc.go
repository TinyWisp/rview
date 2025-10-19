package rview

import (
	"math"
	"reflect"

	"github.com/TinyWisp/rview/ddl"
	"github.com/TinyWisp/rview/tperr"
)

// a function that get a variable according to its name
type VarGetter func(name string) (interface{}, error)

func CalcExp(exp *ddl.Exp, varGetter VarGetter) (*ddl.Exp, error) {
	if exp == nil {
		return &ddl.Exp{
			Type: ddl.ExpNil,
		}, nil
	}

	res := exp
	err := error(nil)

	switch exp.Type {
	case ddl.ExpVar:
		res, err = calcVar(exp, varGetter)

	case ddl.ExpFunc:
		res, err = calcFunc(exp, varGetter)

	case ddl.ExpCalc:
		left, lerr := CalcExp(exp.Left, varGetter)
		if lerr != nil {
			res = nil
			err = lerr
			break
		}

		right, rerr := CalcExp(exp.Right, varGetter)
		if rerr != nil {
			res = nil
			err = rerr
			break
		}

		switch exp.Operator {
		case "+":
			res, err = calcPlus(left, right)

		case "-":
			res, err = calcMinus(left, right)

		case "*":
			res, err = calcTimes(left, right)

		case "/":
			res, err = calcDivision(left, right)

		case ">":
			res, err = calcGreater(left, right)

		case ">=":
			res, err = calcGreaterOrEqual(left, right)

		case "<":
			res, err = calcLess(left, right)

		case "<=":
			res, err = calcLessOrEqual(left, right)

		case "==":
			res, err = calcEqual(left, right)

		case "!=":
			res, err = calcNotEqual(left, right)

		case "&&":
			res, err = calcLogicalAnd(left, right)

		case "||":
			res, err = calcLogicalOr(left, right)

		case "!":
			res, err = calcLogicalNot(left, right)

		case "?":
			condition, terr := CalcExp(exp.TenaryCondition, varGetter)
			if terr != nil {
				res = nil
				err = terr
			} else {
				res, err = calcTernaryCondition(condition, left, right)
			}

		case ".":
			res, err = calcDot(left, right)

		case "[":
			res, err = calcSquareBracket(left, right)

		default:
			res = nil
			err = tperr.NewTypedError("calc.unsupportedOperator", exp.Operator)
		}
	}

	if err != nil {
		if terr, ok := err.(*tperr.TypedError); ok {
			err = ddl.NewDdlError("", exp.Pos, terr.GetEtype(), terr.GetVars()...)
		} else if _, ok := err.(*ddl.DdlError); !ok {
			err = ddl.NewDdlError("", exp.Pos, err.Error())
		}
	}

	return res, err
}

func calcPlus(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  left.Int + right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(left.Int) + right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float + float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float + right.Float,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "+", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcMinus(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  left.Int - right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(left.Int) - right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float - float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float - right.Float,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "-", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcTimes(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  left.Int * right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(left.Int) * right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float * float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float * right.Float,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "*", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcDivision(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  left.Int / right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(left.Int) / right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float / float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: left.Float / right.Float,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "/", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalAnd(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool && right.Bool,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "&&", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalOr(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool || right.Bool,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "||", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalNot(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if (left == nil || left.Type == ddl.ExpNil) && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: !right.Bool,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "!", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcEqual(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	const minDiff = 1e-9

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int == right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(float64(left.Int)-right.Float) < minDiff,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(left.Float-float64(right.Int)) < minDiff,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(left.Float-right.Float) < minDiff,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str == right.Str,
		}, nil
	}

	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool == right.Bool,
		}, nil
	}

	if left.Type == ddl.ExpNil && right.Type == ddl.ExpNil {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: true,
		}, nil
	}

	if left.Type == ddl.ExpInterface && right.Type == ddl.ExpInterface {
		if reflect.TypeOf(left.Interface) == reflect.TypeOf(right.Interface) {
			if reflect.TypeOf(left.Interface).Comparable() {
				return &ddl.Exp{
					Type: ddl.ExpBool,
					Bool: reflect.ValueOf(left.Interface).Equal(reflect.ValueOf(right.Interface)),
				}, nil
			}
			return &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: reflect.ValueOf(left.Interface).Pointer() == reflect.ValueOf(right.Interface).Pointer(),
			}, nil
		}
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", "==", reflect.TypeOf(left.Interface).String(), reflect.TypeOf(right.Interface).String())
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "==", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcNotEqual(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	const minDiff = 1e-9

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int != right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(float64(left.Int)-right.Float) > minDiff,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(left.Float-float64(right.Int)) > minDiff,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: math.Abs(left.Float-right.Float) > minDiff,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str != right.Str,
		}, nil
	}

	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool != right.Bool,
		}, nil
	}

	if left.Type == ddl.ExpNil && right.Type == ddl.ExpNil {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: false,
		}, nil
	}

	if left.Type == ddl.ExpInterface && right.Type == ddl.ExpInterface {
		if reflect.TypeOf(left.Interface) == reflect.TypeOf(right.Interface) {
			if reflect.TypeOf(left.Interface).Comparable() {
				return &ddl.Exp{
					Type: ddl.ExpBool,
					Bool: !reflect.ValueOf(left.Interface).Equal(reflect.ValueOf(right.Interface)),
				}, nil
			}
			return &ddl.Exp{
				Type: ddl.ExpBool,
				Bool: reflect.ValueOf(left.Interface).Pointer() != reflect.ValueOf(right.Interface).Pointer(),
			}, nil
		}
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", "!=", reflect.TypeOf(left.Interface).String(), reflect.TypeOf(right.Interface).String())
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "!=", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcGreater(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int > right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: float64(left.Int) > right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float > float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float > right.Float,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str > right.Str,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", ">", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcGreaterOrEqual(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int >= right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: float64(left.Int) >= right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float >= float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float >= right.Float,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str >= right.Str,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", ">", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLess(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int < right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: float64(left.Int) < right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float < float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float < right.Float,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str < right.Str,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "<", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLessOrEqual(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInt && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Int <= right.Int,
		}, nil
	}

	if left.Type == ddl.ExpInt && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: float64(left.Int) <= right.Float,
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpInt {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float <= float64(right.Int),
		}, nil
	}

	if left.Type == ddl.ExpFloat && right.Type == ddl.ExpFloat {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Float <= right.Float,
		}, nil
	}

	if left.Type == ddl.ExpStr && right.Type == ddl.ExpStr {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Str <= right.Str,
		}, nil
	}

	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "<=", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcTernaryCondition(condition *ddl.Exp, left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if condition == nil || condition.Type != ddl.ExpBool {
		return nil, tperr.NewTypedError("calc.invalidTernaryCondition", ddl.ExpTypeName[condition.Type])
	}

	if left.Type != right.Type {
		return nil, tperr.NewTypedError("calc.ternaryDataNotSameType")
	}

	if condition.Bool {
		return left, nil
	} else {
		return right, nil
	}
}

func calcVar(exp *ddl.Exp, varGetter VarGetter) (*ddl.Exp, error) {
	if exp.Type != ddl.ExpVar {
		return nil, tperr.NewTypedError("calc.expMustBeVarType", `calcVar.exp`)
	}

	if exp.Variable == "" {
		return nil, tperr.NewTypedError("calc.emptyVariableName")
	}

	val, err := varGetter(exp.Variable)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return &ddl.Exp{
			Type: ddl.ExpNil,
		}, nil
	}

	if reflect.ValueOf(val).Kind() == reflect.Pointer {
		val = reflect.ValueOf(val).Elem().Interface()
	}

	return ConvertVariableToExp(val)
}

func calcFunc(exp *ddl.Exp, varGetter VarGetter) (*ddl.Exp, error) {
	if exp.Type != ddl.ExpFunc {
		return nil, tperr.NewTypedError("calc.expMustBeFuncType", `calcFunc.exp`)
	}

	if exp.FuncName == "" {
		return nil, tperr.NewTypedError("calc.emptyFuncName")
	}

	funcVar, err := varGetter(exp.FuncName)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(funcVar).Kind() != reflect.Func {
		return nil, tperr.NewTypedError("calc.varIsNotFunc", exp.FuncName)
	}
	theFunc := reflect.ValueOf(funcVar)
	theFuncType := reflect.TypeOf(funcVar)

	if !theFuncType.IsVariadic() && len(exp.FuncParams) != theFuncType.NumIn() {
		return nil, tperr.NewTypedError("calc.argumentNumberMismatch", exp.FuncName, theFuncType.NumIn(), len(exp.FuncParams))
	}

	if theFuncType.IsVariadic() && len(exp.FuncParams) != theFuncType.NumIn() {
		return nil, tperr.NewTypedError("calc.argumentNumberNotEnough", exp.FuncName, theFuncType.NumIn()-1, len(exp.FuncParams))
	}

	theParams := []reflect.Value{}
	for idx, paramExp := range exp.FuncParams {
		res, rerr := CalcExp(paramExp, varGetter)
		if rerr != nil {
			return nil, rerr
		}

		param := reflect.Value{}
		switch res.Type {
		case ddl.ExpBool:
			param = reflect.ValueOf(res.Bool)

		case ddl.ExpStr:
			param = reflect.ValueOf(res.Str)

		case ddl.ExpInt:
			param = reflect.ValueOf(res.Int).Convert(theFuncType.In(idx))

		case ddl.ExpFloat:
			param = reflect.ValueOf(res.Float).Convert(theFuncType.In(idx))

		case ddl.ExpNil:
			param = reflect.ValueOf(nil)

		case ddl.ExpInterface:
			param = reflect.ValueOf(exp.Interface)
		}
		theParams = append(theParams, param)
	}

	res := theFunc.Call(theParams)
	if len(res) == 0 {
		return &ddl.Exp{
			Type: ddl.ExpNil,
		}, nil
	}

	val := res[0].Interface()

	return ConvertVariableToExp(val)
}

func calcDot(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInterface && right.Type == ddl.ExpStr {
		leftVal := reflect.ValueOf(left.Interface)

		if leftVal.Kind() == reflect.Struct && right.Type == ddl.ExpStr {
			res, err := GetStructField(left.Interface, right.Str)
			if err != nil {
				return nil, err
			}
			return ConvertVariableToExp(res)
		}

		if leftVal.Kind() == reflect.Map {
			mapKeyType := reflect.TypeOf(left.Interface).Key()
			mapKeyKind := mapKeyType.Kind()

			if (mapKeyKind == reflect.Uint8 || mapKeyKind == reflect.Uint16 || mapKeyKind == reflect.Uint32 ||
				mapKeyKind == reflect.Uint64 || mapKeyKind == reflect.Int8 || mapKeyKind == reflect.Int16 ||
				mapKeyKind == reflect.Int32 || mapKeyKind == reflect.Int64) && right.Type == ddl.ExpInt {
				key := reflect.ValueOf(right.Int).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if (mapKeyKind == reflect.Float32 || mapKeyKind == reflect.Float64) && right.Type == ddl.ExpFloat {
				key := reflect.ValueOf(right.Float).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.String && right.Type == ddl.ExpStr {
				key := reflect.ValueOf(right.Str)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.Bool && right.Type == ddl.ExpBool {
				key := reflect.ValueOf(right.Bool)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if right.Type == ddl.ExpInterface && reflect.ValueOf(right.Interface).CanConvert(mapKeyType) {
				key := reflect.ValueOf(right.Interface).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}
		}

		if right.Type != ddl.ExpInterface {
			return nil, tperr.NewTypedError("calc.operandTypeMismatch", ".", leftVal.Type().String(), ddl.ExpTypeName[right.Type])
		}
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", ".", leftVal.Type().String(), reflect.TypeOf(right.Interface).String())
	}

	if right.Type != ddl.ExpInterface {
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", ".", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
	}
	return nil, tperr.NewTypedError("calc.operandTypeMismatch", ".", ddl.ExpTypeName[left.Type], reflect.TypeOf(right.Interface).Name())
}

func calcSquareBracket(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpInterface {
		leftVal := reflect.ValueOf(left.Interface)

		if leftVal.Kind() == reflect.Struct && right.Type == ddl.ExpStr {
			res, err := GetStructField(left.Interface, right.Str)
			if err != nil {
				return nil, err
			}
			return ConvertVariableToExp(res)
		}

		if leftVal.Kind() == reflect.Map {
			mapKeyType := reflect.TypeOf(left.Interface).Key()
			mapKeyKind := mapKeyType.Kind()

			if (mapKeyKind == reflect.Uint8 || mapKeyKind == reflect.Uint16 || mapKeyKind == reflect.Uint32 ||
				mapKeyKind == reflect.Uint64 || mapKeyKind == reflect.Int8 || mapKeyKind == reflect.Int16 ||
				mapKeyKind == reflect.Int32 || mapKeyKind == reflect.Int64) && right.Type == ddl.ExpInt {
				key := reflect.ValueOf(right.Int).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if (mapKeyKind == reflect.Float32 || mapKeyKind == reflect.Float64) && right.Type == ddl.ExpFloat {
				key := reflect.ValueOf(right.Float).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.String && right.Type == ddl.ExpStr {
				key := reflect.ValueOf(right.Str)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.Bool && right.Type == ddl.ExpBool {
				key := reflect.ValueOf(right.Bool)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}

			if right.Type == ddl.ExpInterface && reflect.ValueOf(right.Interface).CanConvert(mapKeyType) {
				key := reflect.ValueOf(right.Interface).Convert(mapKeyType)
				res := leftVal.MapIndex(key).Interface()
				return ConvertVariableToExp(res)
			}
		}

		if leftVal.Kind() == reflect.Array || leftVal.Kind() == reflect.Slice {
			if right.Type == ddl.ExpInt {
				idx := right.Int
				res := leftVal.Index(int(idx)).Interface()
				return ConvertVariableToExp(res)
			}
		}

		if right.Type != ddl.ExpInterface {
			return nil, tperr.NewTypedError("calc.operandTypeMismatch", "[]", leftVal.Type().Name(), ddl.ExpTypeName[right.Type])
		}
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", "[]", leftVal.Type().Name(), reflect.TypeOf(right.Interface).Name())
	}

	if right.Type != ddl.ExpInterface {
		return nil, tperr.NewTypedError("calc.operandTypeMismatch", "[]", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
	}
	return nil, tperr.NewTypedError("calc.operandTypeMismatch", "[]", ddl.ExpTypeName[left.Type], reflect.TypeOf(right.Interface).Name())
}

func ConvertVariableToExp(variable interface{}) (*ddl.Exp, error) {
	if variable == nil {
		return &ddl.Exp{
			Type: ddl.ExpNil,
		}, nil
	}

	switch v := variable.(type) {
	case int:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case int8:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case int16:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case int32:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case int64:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  v,
		}, nil

	case uint8:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case uint16:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case uint32:
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case uint64:
		if v > math.MaxInt64 {
			return &ddl.Exp{
				Type:  ddl.ExpFloat,
				Float: float64(v),
			}, nil
		}
		return &ddl.Exp{
			Type: ddl.ExpInt,
			Int:  int64(v),
		}, nil

	case float32:
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(v),
		}, nil

	case float64:
		return &ddl.Exp{
			Type:  ddl.ExpFloat,
			Float: float64(v),
		}, nil

	case string:
		return &ddl.Exp{
			Type: ddl.ExpStr,
			Str:  v,
		}, nil

	case bool:
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: v,
		}, nil

	default:
		return &ddl.Exp{
			Type:      ddl.ExpInterface,
			Interface: v,
		}, nil
	}
}
