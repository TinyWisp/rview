package comp

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/rivo/tview"
)

func SetProp(prim tview.Primitive, prop string, val interface{}) error {
	vprim := reflect.ValueOf(prim)
	funcName := fmt.Sprintf("Set%s", strcase.ToCamel(prop))
	setter := vprim.MethodByName(funcName)

	if setter.IsZero() {
		return fmt.Errorf("'%s' is not allowed", prop)
	}

	fmt.Printf("%s setter.Call\n", funcName)
	setter.Call([]reflect.Value{reflect.ValueOf(val)})
	return nil
}

func getFromPropMap(props map[string]interface{}, name string, required bool, defaultValue interface{}) (interface{}, error) {
	prop, ok := props[name]
	if !ok {
		if required {
			return nil, fmt.Errorf("property %s is required", name)
		} else {
			return defaultValue, nil
		}
	}

	return prop, nil
}
