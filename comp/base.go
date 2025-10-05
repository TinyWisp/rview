package comp

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview/tperr"
	"github.com/iancoleman/strcase"
)

type Base struct {
	name string
	inst interface{}
}

func (b *Base) GetName() string {
	return b.name
}

func (b *Base) Primitive() Primitive {
	return b.inst.(Primitive)
}

func (b *Base) SetProp(prop string, val interface{}) error {
	primVal := reflect.ValueOf(b.inst)
	funcName := fmt.Sprintf("Set%s", strcase.ToCamel(prop))
	setter := primVal.MethodByName(funcName)

	if !setter.IsValid() {
		return tperr.NewTypedError("comp.SetProp.propNotAllowed", prop, b.GetName())
	}

	setterType := setter.Type()
	if setterType.In(0) == reflect.TypeOf(val) {
		setter.Call([]reflect.Value{reflect.ValueOf(val)})
		return nil
	}

	if reflect.ValueOf(val).CanConvert(setterType.In(0)) {
		setter.Call([]reflect.Value{
			reflect.ValueOf(val).Convert(setterType.In(0)),
		})
		return nil
	}

	return tperr.NewTypedError("comp.SetProp.propTypeMismatch", reflect.TypeOf(val).Name(), prop, b.GetName(), setterType.In(0).Kind())
}

func (b *Base) GetProp(prop string) (interface{}, error) {
	primVal := reflect.ValueOf(b.inst)
	funcName := fmt.Sprintf("Get%s", strcase.ToCamel(prop))
	fmt.Println(funcName)
	setter := primVal.MethodByName(funcName)

	if !setter.IsValid() {
		return nil, tperr.NewTypedError("comp.GetProp.propNotExist", prop, b.GetName())
	}

	res := setter.Call([]reflect.Value{})
	return res[0].Interface(), nil
}

func (b *Base) CanAddItem() bool {
	return false
}

func (b *Base) AddItem(comp *Component, props map[string]interface{}) error {
	return nil
}
