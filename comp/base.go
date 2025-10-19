package comp

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview/tperr"
	"github.com/iancoleman/strcase"
)

type Base[T any] struct {
	name      string
	outerInst interface{}
	tviewInst T
}

func (b *Base[T]) GetName() string {
	return b.name
}

func (b *Base[T]) SetProp(prop string, val interface{}) error {
	outerInstVal := reflect.ValueOf(b.outerInst)
	tviewInstVal := reflect.ValueOf(b.tviewInst)
	funcName := fmt.Sprintf("Set%s", strcase.ToCamel(prop))

	outerSetter := outerInstVal.MethodByName(funcName)
	tviewSetter := tviewInstVal.MethodByName(funcName)
	if !outerSetter.IsValid() && !tviewSetter.IsValid() {
		return tperr.NewTypedError("comp.SetProp.propNotAllowed", prop, b.GetName())
	}

	setter := outerSetter
	if !setter.IsValid() {
		setter = tviewSetter
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

func (b *Base[T]) GetProp(prop string) (interface{}, error) {
	outerInstVal := reflect.ValueOf(b.outerInst)
	tviewInstVal := reflect.ValueOf(b.tviewInst)
	funcName := fmt.Sprintf("Get%s", strcase.ToCamel(prop))

	outerGetter := outerInstVal.MethodByName(funcName)
	tviewGetter := tviewInstVal.MethodByName(funcName)
	if !outerGetter.IsValid() && !tviewGetter.IsValid() {
		return nil, tperr.NewTypedError("comp.GetProp.propNotExist", prop, b.GetName())
	}

	getter := outerGetter
	if !getter.IsValid() {
		getter = tviewGetter
	}
	res := getter.Call([]reflect.Value{})
	return res[0].Interface(), nil
}

func (b *Base[T]) CanAddItem() bool {
	return false
}

func (b *Base[T]) AddItem(comp *Component, props map[string]interface{}) error {
	return nil
}
