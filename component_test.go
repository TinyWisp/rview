package rview

import (
	"fmt"
	"testing"
)

func TestSetCompProp(t *testing.T) {
	type Component1 struct {
		Inst      *ComponentInstance
		IntVar    int
		PIntVar   *int
		FloatVar  float64
		PFloatVar *float64
		StrVar    string
		PStrVar   *string
		BoolVar   bool
		PBoolVar  *bool
	}

	comp1 := &Component1{}
	compInst1 := &ComponentInstance{
		Comp: comp1,
	}

	fmt.Println(compInst1.SetCompProp("IntVar", 333))
	if comp1.IntVar != 333 {
		t.Fatalf("ComponentInstance.SetCompProp, int")
	}

	tmpIntVar := 1
	compInst1.SetCompProp("PIntVar", &tmpIntVar)
	if comp1.PIntVar != &tmpIntVar {
		t.Fatalf("ComponentInstance.SetCompProp, *int")
	}

	compInst1.SetCompProp("FloatVar", 3.33333)
	if comp1.FloatVar-3.33333 < -0.0001 || comp1.FloatVar-3.3333 > 0.0001 {
		t.Fatalf("ComponentInstance.SetCompProp, float")
	}

	tmpFloatVar := 1.2222
	compInst1.SetCompProp("PFloatVar", &tmpFloatVar)
	if comp1.PFloatVar != &tmpFloatVar {
		t.Fatalf("ComponentInstance.SetCompProp, *float")
	}

	compInst1.SetCompProp("StrVar", "hello")
	if comp1.StrVar != "hello" {
		t.Fatalf("ComponentInstance.SetCompProp, string")
	}

	tmpStrVar := "hello"
	compInst1.SetCompProp("PStrVar", &tmpStrVar)
	if comp1.PStrVar != &tmpStrVar {
		t.Fatalf("ComponentInstance.SetCompProp, *string")
	}

	compInst1.SetCompProp("BoolVar", false)
	if comp1.BoolVar != false {
		t.Fatalf("ComponentInstance.SetCompProp, bool")
	}

	tmpBoolVar := false
	compInst1.SetCompProp("PBoolVar", &tmpBoolVar)
	if comp1.PBoolVar != &tmpBoolVar {
		t.Fatalf("ComponentInstance.SetCompProp, *bool")
	}
}
