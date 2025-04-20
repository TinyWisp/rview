package rview

import (
	"reflect"
	"unicode"
)

func IsStructFieldExported(field string) bool {
	if len(field) == 0 {
		return false
	}
	return unicode.IsUpper(rune(field[0]))
}

func SetStructField(structVar interface{}, field string, val interface{}) error {
	comp := reflect.ValueOf(structVar)
	if comp.Kind() == reflect.Pointer {
		comp = comp.Elem()
	}
	fieldVal := comp.FieldByName(field)

	if !fieldVal.IsValid() {
		return NewError("util.SetStructField.fieldNotExist", field)
	}

	if !IsStructFieldExported(field) {
		return NewError("util.SetStructField.unexportedField", field)
	}

	if !fieldVal.CanSet() {
		return NewError("util.SetStructField.cannotSetFieldValue", field)
	}

	if isRef(fieldVal.Interface()) {
		if typable, ok := fieldVal.Interface().(interface{ Type() reflect.Type }); ok {
			if reflect.TypeOf(val) != typable.Type() && !reflect.ValueOf(val).CanConvert(typable.Type()) {
				return NewError("util.SetStructField.typeMismatch", reflect.TypeOf(fieldVal), typable.Type())
			}
		}

		setMethod := fieldVal.MethodByName("Set")
		setMethod.Call([]reflect.Value{reflect.ValueOf(val)})
		return nil
	}

	if fieldVal.Type() != reflect.TypeOf(val) {
		if !reflect.ValueOf(val).CanConvert(fieldVal.Type()) {
			return NewError("util.SetStructField.typeMismatch", reflect.TypeOf(fieldVal), reflect.TypeOf(val))
		}
		fieldVal.Set(reflect.ValueOf(val).Convert(fieldVal.Type()))
		return nil
	}

	fieldVal.Set(reflect.ValueOf(val))
	return nil
}

func GetStructField(structVar interface{}, field string) (interface{}, error) {
	comp := reflect.ValueOf(structVar)
	if comp.Kind() == reflect.Pointer {
		comp = comp.Elem()
	}
	fieldVal := comp.FieldByName(field)

	if !fieldVal.IsValid() {
		return nil, NewError("util.GetStructField.fieldNotExist", field)
	}

	if !IsStructFieldExported(field) {
		return nil, NewError("util.GetStructField.unexportedField", field)
	}

	if isRef(fieldVal.Interface()) {
		getMethod := fieldVal.MethodByName("Get")
		vals := getMethod.Call([]reflect.Value{})
		return vals[0].Interface(), nil
	}

	return fieldVal.Interface(), nil
}
