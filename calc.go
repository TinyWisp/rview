package rview

import (
	"math"
	"reflect"

	"github.com/TinyWisp/rview/ddl"
)

type VarGetter func(name string) (interface{}, error)

func CalcExp(exp *ddl.Exp, varGetter VarGetter) (*ddl.Exp, error) {
	if exp.Type == ddl.ExpCalc {
		left, lerr := CalcExp(exp.Left, varGetter)
		if lerr != nil {
			return nil, lerr
		}

		right, rerr := CalcExp(exp.Right, varGetter)
		if rerr != nil {
			return nil, rerr
		}

		switch exp.Operator {
		case "+":
			return calcPlus(left, right)

		case "-":
			return calcMinus(left, right)

		case "*":
			return calcTimes(left, right)

		case "/":
			return calcDivision(left, right)

		case ">":
			return calcGreater(left, right)

		case ">=":
			return calcGreaterOrEqual(left, right)

		case "<":
			return calcLess(left, right)

		case "<=":
			return calcLessOrEqual(left, right)

		case "==":
			return calcEqual(left, right)

		case "!=":
			return calcNotEqual(left, right)

		case "&&":
			return calcLogicalAnd(left, right)

		case "||":
			return calcLogicalOr(left, right)

		case "!":
			return calcLogicalNot(left, right)

		case "?":
			condition, terr := CalcExp(exp.TenaryCondition, varGetter)
			if terr != nil {
				return nil, terr
			}
			return calcTenaryCondition(condition, left, right)

		case ".":
			return calcDot(left, right)

		case "[":
			return calcSquareBracket(left, right)

		default:
			return nil, NewError("calc.unsupportedOperator", exp.Operator)
		}
	}

	if exp.Type == ddl.ExpVar {
		return calcVar(exp, varGetter)
	}

	return exp, nil
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

	return nil, NewError("calc.operandTypeMismatch", "+", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", "-", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", "*", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", "/", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalAnd(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool && right.Bool,
		}, nil
	}

	return nil, NewError("calc.operandTypeMismatch", "&&", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalOr(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpBool && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: left.Bool || right.Bool,
		}, nil
	}

	return nil, NewError("calc.operandTypeMismatch", "||", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcLogicalNot(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if left.Type == ddl.ExpNil && right.Type == ddl.ExpBool {
		return &ddl.Exp{
			Type: ddl.ExpBool,
			Bool: !right.Bool,
		}, nil
	}

	return nil, NewError("calc.operandTypeMismatch", "!", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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
		return nil, NewError("calc.operandTypeMismatch", "==", reflect.TypeOf(left.Interface).String(), reflect.TypeOf(right.Interface).String())
	}

	return nil, NewError("calc.operandTypeMismatch", "==", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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
		return nil, NewError("calc.operandTypeMismatch", "!=", reflect.TypeOf(left.Interface).String(), reflect.TypeOf(right.Interface).String())
	}

	return nil, NewError("calc.operandTypeMismatch", "!=", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", ">", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", ">", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", "<", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
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

	return nil, NewError("calc.operandTypeMismatch", "<=", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcTenaryCondition(condition *ddl.Exp, left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
	if condition.Type == ddl.ExpBool {
		if condition.Bool {
			return left, nil
		} else {
			return right, nil
		}
	}

	return nil, NewError("calc.operandTypeMismatch", "?", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
}

func calcVar(exp *ddl.Exp, varGetter VarGetter) (*ddl.Exp, error) {
	if exp.Type != ddl.ExpVar {
		return nil, NewError("calc.mustBeVarType")
	}

	if exp.Variable == "" {
		return nil, NewError("calc.emptyVariableName")
	}

	val, err := varGetter(exp.Variable)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(val).Kind() == reflect.Pointer {
		val = reflect.ValueOf(val).Elem().Interface()
	}

	switch v := val.(type) {
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

func ConvertVariableToExp(variable interface{}) (*ddl.Exp, error) {
	switch v := variable.(type) {
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

func calcDot(left *ddl.Exp, right *ddl.Exp) (*ddl.Exp, error) {
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
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if (mapKeyKind == reflect.Float32 || mapKeyKind == reflect.Float64) && right.Type == ddl.ExpFloat {
				key := reflect.ValueOf(right.Float).Convert(mapKeyType)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.String && right.Type == ddl.ExpStr {
				key := reflect.ValueOf(right.Str)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.Bool && right.Type == ddl.ExpBool {
				key := reflect.ValueOf(right.Bool)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if right.Type == ddl.ExpInterface && reflect.ValueOf(right.Interface).CanConvert(mapKeyType) {
				key := reflect.ValueOf(right.Interface).Convert(mapKeyType)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}
		}

		if leftVal.Kind() == reflect.Array || leftVal.Kind() == reflect.Slice {
			if right.Type == ddl.ExpInt {
				idx := right.Int
				res := leftVal.Index(int(idx))
				return ConvertVariableToExp(res)
			}
		}

		if right.Type != ddl.ExpInterface {
			return nil, NewError("calc.operandTypeMismatch", ".", leftVal.Type().String(), ddl.ExpTypeName[right.Type])
		}
		return nil, NewError("calc.operandTypeMismatch", ".", leftVal.Type().String(), reflect.TypeOf(right.Interface).String())
	}

	if right.Type != ddl.ExpInterface {
		return nil, NewError("calc.operandTypeMismatch", ".", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
	}
	return nil, NewError("calc.operandTypeMismatch", ".", ddl.ExpTypeName[left.Type], reflect.TypeOf(right.Interface).Name())
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
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if (mapKeyKind == reflect.Float32 || mapKeyKind == reflect.Float64) && right.Type == ddl.ExpFloat {
				key := reflect.ValueOf(right.Float).Convert(mapKeyType)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.String && right.Type == ddl.ExpStr {
				key := reflect.ValueOf(right.Str)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if mapKeyKind == reflect.Bool && right.Type == ddl.ExpBool {
				key := reflect.ValueOf(right.Bool)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}

			if right.Type == ddl.ExpInterface && reflect.ValueOf(right.Interface).CanConvert(mapKeyType) {
				key := reflect.ValueOf(right.Interface).Convert(mapKeyType)
				res := leftVal.MapIndex(key)
				return ConvertVariableToExp(res)
			}
		}

		if leftVal.Kind() == reflect.Array || leftVal.Kind() == reflect.Slice {
			if right.Type == ddl.ExpInt {
				idx := right.Int
				res := leftVal.Index(int(idx))
				return ConvertVariableToExp(res)
			}
		}

		if right.Type != ddl.ExpInterface {
			return nil, NewError("calc.operandTypeMismatch", "[]", leftVal.Type().Name(), ddl.ExpTypeName[right.Type])
		}
		return nil, NewError("calc.operandTypeMismatch", "[]", leftVal.Type().Name(), reflect.TypeOf(right.Interface).Name())
	}

	if right.Type != ddl.ExpInterface {
		return nil, NewError("calc.operandTypeMismatch", "[]", ddl.ExpTypeName[left.Type], ddl.ExpTypeName[right.Type])
	}
	return nil, NewError("calc.operandTypeMismatch", "[]", ddl.ExpTypeName[left.Type], reflect.TypeOf(right.Interface).Name())
}
