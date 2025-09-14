package comp

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview/tperr"
	"github.com/iancoleman/strcase"
	"github.com/rivo/tview"
)

type Base[T tview.Primitive] struct {
	name string
	inst T
}

func (b *Base[T]) GetName() string {
	return b.name
}

func (b *Base[T]) Primitive() tview.Primitive {
	return b.inst
}

func (b *Base[T]) SetProp(prop string, val interface{}) error {
	vprim := reflect.ValueOf(b.inst)
	funcName := fmt.Sprintf("Set%s", strcase.ToCamel(prop))
	setter := vprim.MethodByName(funcName)

	if setter.IsValid() {
		return tperr.NewTypedError("comp.SetProp.propNotAllowed", prop, b.GetName())
	}

	setterType := setter.Type()
	if setterType.In(0) != reflect.TypeOf(val) {
		return tperr.NewTypedError("comp.SetProp.propTypeMismatch", reflect.TypeOf(val).Name(), setterType.In(0).Name(), b.GetName())
	}

	setter.Call([]reflect.Value{reflect.ValueOf(val)})
	return nil
}

func (b *Base[T]) GetProp(prop string) (interface{}, error) {
	vprim := reflect.ValueOf(b.inst)
	funcName := fmt.Sprintf("Get%s", strcase.ToCamel(prop))
	setter := vprim.MethodByName(funcName)

	if setter.IsValid() {
		return nil, tperr.NewTypedError("comp.GetProp.propNotExist", prop, b.GetName())
	}

	res := setter.Call([]reflect.Value{reflect.ValueOf(prop)})
	return res[0].Interface(), nil
}

func (b *Base[T]) CanAddItem() bool {
	return false
}

func (b *Base[T]) AddItem(comp *Component, props map[string]interface{}) error {
	return nil
}
