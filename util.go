package rview

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/TinyWisp/rview/tperr"
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
		return tperr.NewTypedError("util.SetStructField.fieldNotExist", field)
	}

	if !IsStructFieldExported(field) {
		return tperr.NewTypedError("util.SetStructField.unexportedField", field)
	}

	if !fieldVal.CanSet() {
		return tperr.NewTypedError("util.SetStructField.cannotSetFieldValue", field)
	}

	if isRef(fieldVal.Interface()) {
		if typable, ok := fieldVal.Interface().(interface{ Type() reflect.Type }); ok {
			if reflect.TypeOf(val) != typable.Type() && !reflect.ValueOf(val).CanConvert(typable.Type()) {
				return tperr.NewTypedError("util.SetStructField.typeMismatch", reflect.TypeOf(fieldVal), typable.Type())
			}
		}

		setMethod := fieldVal.MethodByName("Set")
		setMethod.Call([]reflect.Value{reflect.ValueOf(val)})
		return nil
	}

	if fieldVal.Type() != reflect.TypeOf(val) {
		if !reflect.ValueOf(val).CanConvert(fieldVal.Type()) {
			return tperr.NewTypedError("util.SetStructField.typeMismatch", reflect.TypeOf(fieldVal), reflect.TypeOf(val))
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
		return nil, tperr.NewTypedError("util.GetStructField.fieldNotExist", field)
	}

	if !fieldVal.CanInterface() {
		return nil, tperr.NewTypedError("util.GetStructField.unavailableField", field)
	}

	if isRef(fieldVal.Interface()) {
		getMethod := fieldVal.MethodByName("Get")
		vals := getMethod.Call([]reflect.Value{})
		return vals[0].Interface(), nil
	}

	return fieldVal.Interface(), nil
}

func sprintComponentNode(node *ComponentNode, level int) string {
	spaces := strings.Repeat("    ", level)
	nspaces := strings.Repeat("    ", level+1)

	str := ""
	if node.Comp == nil {
		str += fmt.Sprintln(spaces, "Comp: nil")
	} else {
		str += fmt.Sprintln(spaces, "Comp: ", node.Comp.GetName())
	}
	str += fmt.Sprintln(spaces, "HasIf: ", node.HasIf)
	str += fmt.Sprintln(spaces, "If: ", node.If)
	str += fmt.Sprintln(spaces, "HasElseIf: ", node.HasElseIf)
	str += fmt.Sprintln(spaces, "ElseIf: ", node.ElseIf)
	str += fmt.Sprintln(spaces, "HasElse: ", node.HasElse)
	str += fmt.Sprintln(spaces, "HasFor: ", node.HasFor)
	str += fmt.Sprintln(spaces, "Ignore: ", node.Ignore)
	str += fmt.Sprintln(spaces, "InheritVars: ", node.InheritVars)

	str += fmt.Sprintln(spaces, "Vars: ", len(node.Vars))
	for key, val := range node.Vars {
		str += fmt.Sprintln(nspaces, key, ":", val)
	}

	str += fmt.Sprintln(spaces, "Children:", len(node.Children))
	for idx, child := range node.Children {
		str += fmt.Sprintln(nspaces, idx)
		str += sprintComponentNode(child, level+1)
	}

	return str
}
