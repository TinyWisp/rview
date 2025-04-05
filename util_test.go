package rview

import (
	"fmt"
	"testing"
)

func TestIsStructFieldExported(t *testing.T) {
	if IsStructFieldExported("lowercase") {
		t.Fatalf("IsStructFieldExported('lowercase')")
	}

	if !IsStructFieldExported("Uppercase") {
		t.Fatalf("IsStructFieldExported('Uppercase')")
	}
}

func TestSetStructField(t *testing.T) {
	type TinyStruct1 struct {
		a int
	}

	type TinyStruct2 struct {
		b int
	}

	type TestStruct struct {
		IntVar     int
		PIntVar    *int
		RIntVar    *Ref[int]
		FloatVar   float64
		PFloatVar  *float64
		RFloatVar  *Ref[float64]
		StrVar     string
		PStrVar    *string
		RStrVar    *Ref[string]
		BoolVar    bool
		PBoolVar   *bool
		RBoolVar   *Ref[bool]
		StructVar  TinyStruct1
		PStructVar *TinyStruct1
		RStructVar *Ref[TinyStruct1]
		unexported int
	}

	test := &TestStruct{}

	err := SetStructField(test, "IntVar", 333)
	if err != nil || test.IntVar != 333 {
		t.Fatalf("util.SetStructField, int")
	}

	tmpIntVar := 1
	err = SetStructField(test, "PIntVar", &tmpIntVar)
	if err != nil || test.PIntVar != &tmpIntVar {
		t.Fatalf("util.SetStructField, *int")
	}

	test.RIntVar = NewRef(1)
	err = SetStructField(test, "RIntVar", 2)
	if err != nil || test.RIntVar.Get() != 2 {
		t.Fatalf("util.SetStructField, Ref[int]")
	}

	err = SetStructField(test, "IntVar", 666.12)
	if err != nil || test.IntVar != 666 {
		t.Fatalf("util.SetStructField, float -> int")
	}

	tmpInt8 := int8(9)
	err = SetStructField(test, "IntVar", tmpInt8)
	if err != nil || test.IntVar != 9 {
		t.Fatalf("util.SetStructField, int8 -> int")
	}

	err = SetStructField(test, "IntVar", "hello")
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		fmt.Println(err)
		t.Fatalf("util.SetStructField, invalid assignment, str -> int")
	}

	err = SetStructField(test, "IntVar", true)
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, bool -> int")
	}

	err = SetStructField(test, "IntVar", &tmpIntVar)
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, *int -> int")
	}

	tmpStruct := TinyStruct1{a: 1}
	err = SetStructField(test, "IntVar", tmpStruct)
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, struct -> int")
	}

	err = SetStructField(test, "PIntVar", &tmpStruct)
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, *struct -> *int")
	}

	err = SetStructField(test, "StructVar", TinyStruct2{})
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, struct1 -> struct2")
	}

	err = SetStructField(test, "PStructVar", &TinyStruct2{})
	if err == nil || !IsErrorType(err, "util.SetStructField.typeMismatch") {
		t.Fatalf("util.SetStructField, invalid assignment, *struct1 -> *struct2")
	}

	err = SetStructField(test, "Abcd", "123")
	if err == nil || !IsErrorType(err, "util.SetStructField.fieldNotExist") {
		t.Fatalf("util.SetStructField, not exist field")
	}

	err = SetStructField(test, "Abcd", "123")
	if err == nil || !IsErrorType(err, "util.SetStructField.fieldNotExist") {
		t.Fatalf("util.SetStructField, not exist field")
	}

	err = SetStructField(test, "unexported", 123)
	if err == nil || !IsErrorType(err, "util.SetStructField.unexportedField") {
		fmt.Println(err)
		t.Fatalf("util.SetStructField, unexported field")
	}
}

func TestGetStructField(t *testing.T) {
	type TinyStruct1 struct {
		a int
	}

	type TinyStruct2 struct {
		b int
	}

	type TestStruct struct {
		IntVar     int
		PIntVar    *int
		RIntVar    *Ref[int]
		FloatVar   float64
		PFloatVar  *float64
		RFloatVar  *Ref[float64]
		StrVar     string
		PStrVar    *string
		RStrVar    *Ref[string]
		BoolVar    bool
		PBoolVar   *bool
		RBoolVar   *Ref[bool]
		StructVar  TinyStruct1
		PStructVar *TinyStruct1
		RStructVar *Ref[TinyStruct1]
		unexported int
	}

	test := &TestStruct{
		IntVar: 3,
	}

	intVar, err := GetStructField(test, "IntVar")
	if err != nil || intVar.(int) != 3 {
		t.Fatalf("util.GetStructField, int")
	}

	_, err = GetStructField(test, "unexported")
	if err == nil || !IsErrorType(err, "util.GetStructField.unexportedField") {
		t.Fatalf("util.GetStructField, unexported field")
	}
}
