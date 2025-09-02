package comp

import (
	"fmt"
	"reflect"

	"github.com/TinyWisp/rview"
	"github.com/iancoleman/strcase"
	"github.com/rivo/tview"
)

type Base[T tview.Primitive] struct {
	inst T
}

func (b *Base[T]) Primitive() tview.Primitive {
	return b.inst
}

func (b *Base[T]) SetProp(prop string, val interface{}) error {
	vprim := reflect.ValueOf(b.inst)
	funcName := fmt.Sprintf("Set%s", strcase.ToCamel(prop))
	setter := vprim.MethodByName(funcName)

	if setter.IsValid() {
		return rview.NewTypedError("comp.propNotAllowed")
	}

	setterType := setter.Type()
	if setterType.NumIn() != 1 {
		return rview.NewTypedError("comp.propSetterNotOneParameter")
	}

	if setterType.In(0) != reflect.TypeOf(val) {
		return rview.NewTypedError("comp.propTypeMismatch")
	}

	setter.Call([]reflect.Value{reflect.ValueOf(val)})
	return nil
}

func (b *Base[T]) CanAddItem() bool {
	return false
}

func (b *Base[T]) AddItem(comp *Component, props map[string]interface{}) error {
	return nil
}
