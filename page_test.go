package rview

import (
	"testing"

	"github.com/TinyWisp/rview/comp"
	"github.com/TinyWisp/rview/ddl"
	"github.com/TinyWisp/rview/tperr"
)

type TestDef struct {
	Tpl            string
	Components     map[string]func() comp.Component
	ChickenStruct  AnimalStruct
	DuckStruct     AnimalStruct
	Int8Var2       int8
	Int16Var2      int16
	Int32Var2      int32
	Int64Var2      int64
	Uint8Var2      int8
	Uint16Var2     int16
	Uint32Var2     int32
	Uint64Var2     int64
	Float32Var3    float32
	Float64Var3    float64
	StringVarHello string
	StringVarWorld string
	StringVarName  string
	BoolVarFalse   bool
	BoolVarTrue    bool
	NilVar         *int
	MapStrStr      map[string]string
	ArrStr         []string
	ArrInt         []int

	FuncWithoutParamsNorReturn func()
	FuncPlus                   func(int, int) int
	FuncMinus                  func(int, int) int
	FuncTimes                  func(int, int) int
	FuncDivision               func(int, int) int
	FuncPlusReturnMulti        func(int, int) (int, error)
}

var testDef = TestDef{
	ChickenStruct:  AnimalStruct{Name: "chicken"},
	DuckStruct:     AnimalStruct{Name: "duck"},
	Int8Var2:       int8(2),
	Int16Var2:      int16(2),
	Int32Var2:      int32(2),
	Int64Var2:      int64(2),
	Uint8Var2:      int8(2),
	Uint16Var2:     int16(2),
	Uint32Var2:     int32(2),
	Uint64Var2:     int64(2),
	Float32Var3:    float32(3.0),
	Float64Var3:    float64(3.0),
	StringVarHello: "hello",
	StringVarWorld: "world",
	StringVarName:  "Name",
	BoolVarFalse:   false,
	BoolVarTrue:    true,
	NilVar:         nil,
	MapStrStr: map[string]string{
		"hello": "hello",
		"world": "world",
	},
	ArrStr: []string{"hello", "world"},
	ArrInt: []int{10, 11, 12, 13},

	FuncWithoutParamsNorReturn: func() {},
	FuncPlus:                   func(a int, b int) int { return a + b },
	FuncMinus:                  func(a int, b int) int { return a - b },
	FuncTimes:                  func(a int, b int) int { return a * b },
	FuncDivision:               func(a int, b int) int { return a / b },
	FuncPlusReturnMulti:        func(a int, b int) (int, error) { return a + b, nil },
}

type CreateNodeTestCase struct {
	tpl    string
	expect *ComponentNode
	err    string
}

var createNodeTestCases []CreateNodeTestCase = []CreateNodeTestCase{
	/*
		{
			tpl: "<template></template>",
			err: "page.tplMustContainOneNode",
		},
		{
			tpl: "<template><flex></flex></template>",
			expect: &ComponentNode{
				Comp:        comp.CreateTemplate(),
				Parent:      nil,
				InheritVars: true,
				Children: []*ComponentNode{
					{
						Comp:        comp.CreateFlex(),
						InheritVars: true,
					},
				},
			},
		},
		{
			tpl: "<template><flex/></template>",
			expect: &ComponentNode{
				Comp:        comp.CreateTemplate(),
				Parent:      nil,
				InheritVars: true,
				Children: []*ComponentNode{
					{
						Comp:        comp.CreateFlex(),
						InheritVars: true,
					},
				},
			},
		},
		{
			tpl: `<template>
						<flex>
							<box v-if="3>1" />
						</flex>
					</template>`,
			expect: &ComponentNode{
				Comp:        comp.CreateTemplate(),
				Parent:      nil,
				InheritVars: true,
				Children: []*ComponentNode{
					{
						Comp:        comp.CreateFlex(),
						InheritVars: true,
						Children: []*ComponentNode{
							{
								Comp:        comp.CreateBox(),
								InheritVars: true,
								HasIf:       true,
								If:          true,
							},
						},
					},
				},
			},
		},
		{
			tpl: `<template>
					<flex>
						<box v-if="false" />
					</flex>
				</template>`,
			expect: &ComponentNode{
				Comp:        comp.CreateTemplate(),
				Parent:      nil,
				InheritVars: true,
				Children: []*ComponentNode{
					{
						Comp:        comp.CreateFlex(),
						InheritVars: true,
						Children: []*ComponentNode{
							{
								Comp:        nil,
								InheritVars: true,
								Ignore:      true,
								HasIf:       true,
								If:          false,
							},
						},
					},
				},
			},
		},
	*/
	{
		tpl: `<template>
			<flex>
				<box v-if="int8var2 > 0" />
			</flex>
		</template>`,
		err: "page.undefinedVariable",
	},
	{
		tpl: `<template>
			<flex>
				<box v-if="Int8Var2 > 0" />
			</flex>
		</template>`,
		expect: &ComponentNode{
			Comp:        comp.CreateTemplate(),
			Parent:      nil,
			InheritVars: true,
			Children: []*ComponentNode{
				{
					Comp:        comp.CreateFlex(),
					InheritVars: true,
					Children: []*ComponentNode{
						{
							Comp:        comp.CreateBox(),
							InheritVars: true,
							Ignore:      false,
							HasIf:       true,
							If:          true,
						},
					},
				},
			},
		},
	},
	{
		tpl: `<template>
			<flex>
				<box v-if="Int8Var2 > 0" />
			</flex>
		</template>`,
		expect: &ComponentNode{
			Comp:        comp.CreateTemplate(),
			Parent:      nil,
			InheritVars: true,
			Children: []*ComponentNode{
				{
					Comp:        comp.CreateFlex(),
					InheritVars: true,
					Children: []*ComponentNode{
						{
							Comp:        comp.CreateBox(),
							InheritVars: true,
							Ignore:      false,
							HasIf:       true,
							If:          true,
						},
					},
				},
			},
		},
	},
}

func isComponentNodeEqual(a *ComponentNode, b *ComponentNode) bool {
	if (a.Comp == nil && b.Comp != nil) || (a.Comp != nil && b.Comp == nil) {
		return false
	}

	if a.Comp != nil && b.Comp != nil && a.Comp.GetName() != b.Comp.GetName() {
		return false
	}

	if a.HasIf != b.HasIf || a.If != b.If {
		return false
	}

	if a.HasElseIf != b.HasElseIf || a.ElseIf != b.ElseIf {
		return false
	}

	if a.HasElse != b.HasElse || a.Else != b.Else {
		return false
	}

	if a.Ignore != b.Ignore {
		return false
	}

	if a.HasFor != b.HasFor {
		return false
	}

	if a.InheritVars != b.InheritVars {
		return false
	}

	if len(a.Vars) != len(b.Vars) {
		return false
	}

	for key, val := range a.Vars {
		bval, ok := b.Vars[key]
		if !ok || val != bval {
			return false
		}
	}

	if len(a.Children) != len(b.Children) {
		return false
	}

	for idx, child := range a.Children {
		if equal := isComponentNodeEqual(child, b.Children[idx]); !equal {
			return false
		}
	}

	return true
}

func TestGetVar(t *testing.T) {
	page := Page{
		def: testDef,
	}
	page.root = &ComponentNode{
		Comp:        comp.CreateTemplate(),
		Parent:      nil,
		InheritVars: true,
		Children: []*ComponentNode{
			{
				Comp:        comp.CreateFlex(),
				InheritVars: true,
				Vars: map[string]interface{}{
					"hello": "hello",
					"world": "world",
				},
				Children: []*ComponentNode{
					{
						Comp:        comp.CreateBox(),
						InheritVars: true,
					},
					{
						Comp:        comp.CreateBox(),
						InheritVars: false,
					},
				},
			},
		},
	}
	page.root.Children[0].Parent = page.root
	page.root.Children[0].Children[0].Parent = page.root.Children[0]
	page.root.Children[0].Children[1].Parent = page.root.Children[0]
	variable, err := page.getVarForNode(page.root.Children[0].Children[0], "Int8Var2")
	if err != nil {
		t.Fatal(err)
	}
	int8var2, ok := variable.(int8)
	if !ok {
		t.Fatalf("node:page.root.Children[0].Children[0] var:Int8Var2 expect:2 get: not int8")
	}
	if int8var2 != 2 {
		t.Fatalf("node:page.root.Children[0].Children[0] var:Int8Var2 expect:2 get: %d", int8var2)
	}

	variable, err = page.getVarForNode(page.root.Children[0].Children[0], "hello")
	if err != nil {
		t.Fatal(err)
	}
	hello, ok := variable.(string)
	if !ok {
		t.Fatalf("node:page.root.Children[0].Children[0] var: hello expect: \"hello\" get: not string")
	}
	if hello != "hello" {
		t.Fatalf("node:page.root.Children[0].Children[0] var: hello expect: \"hello\" get: %s", hello)
	}

	variable, err = page.getVarForNode(page.root.Children[0].Children[1], "Int32Var2")
	if err != nil {
		t.Fatal(err)
	}
	int32var2, ok := variable.(int32)
	if !ok {
		t.Fatalf("node:page.root.Children[0].Children[1] var:Int32Var2 expect: 2 get: not int32")
	}
	if int32var2 != 2 {
		t.Fatalf("node:page.root.Children[0].Children[1] var:Int32Var2 expect: 2 get: %d", int32var2)
	}

	_, err = page.getVarForNode(page.root.Children[0].Children[1], "world")
	if err == nil {
		t.Fatal("th 'world' variable shouldn't be available")
	}
}

func TestCreateNode(t *testing.T) {
	for _, testCase := range createNodeTestCases {
		t.Log("----------------------")
		t.Log(testCase.tpl)

		def := testDef
		def.Tpl = testCase.tpl

		page, err := NewPage(def)
		if err != nil && testCase.err == "" {
			t.Fatal(err)
		}

		if err != nil && testCase.err != "" {
			notExpectedErr := true
			if terr, ok := err.(*tperr.TypedError); ok && terr.Is(testCase.err) {
				notExpectedErr = false
			} else if derr, ok := err.(*ddl.DdlError); ok && derr.Is(testCase.err) {
				notExpectedErr = false
			}
			if notExpectedErr {
				t.Fatalf("the error occured during the test is not as expected.\n expect: %s\nactual: %s\n", testCase.err, err.Error())
			}
			continue
		}

		real := page.root
		expect := testCase.expect
		if !isComponentNodeEqual(real, expect) {
			t.Log(sprintComponentNode(real, 0))
			t.Log(sprintComponentNode(expect, 0))
			t.Fatalf("the generated component node is not as the expected\n")
		}
	}
}
