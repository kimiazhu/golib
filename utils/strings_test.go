// Description: util
// Author: ZHU HAIHUA
// Since: 2016-04-08 22:33
package util

import (
	"testing"
)

type TestType int

type TestStruct struct {
	V  string
	V3 map[string]TestStructA
	v2 []TestType
}

type TestStructA struct {
	v []int
}

// since the map is unsorted, the last one may be different from expected
func TestReflectToStringLong(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected string
	}{
		{true, `true`},
		{32.68, `32.68`},
		{0xFA, `250`},
		{struct{ v string }{v: "ss"}, `struct { v string }{v=ss}`},
		{struct{ v int }{v: -1}, `struct { v int }{v=-1}`},
		{struct {
			v1 float64
			V2 bool
		}{v1: 1.01, V2: true}, `struct { v1 float64; V2 bool }{v1=1.01, V2=true}`},
		{struct {
			v1 float64
			V2 []string
		}{v1: 1.01, V2: []string{"a", "v"}}, `struct { v1 float64; V2 []string }{v1=1.01, V2=[]string[a,v]}`},
		{struct {
			v1 TestType
			V2 []interface{}
		}{v1: 2, V2: []interface{}{"a", 1, true, struct{ v TestType }{-1}}},
			`struct { v1 util.TestType; V2 []interface {} }{v1=2, V2=[]interface {}[interface {}(a),interface {}(1),interface {}(true),interface {}(struct { v util.TestType }{v=-1})]}`},
		{TestStruct{V: "valV", v2: []TestType{9, 8, 7}, V3: map[string]TestStructA{"k2": TestStructA{[]int{0, 1, 2}}, "k1": TestStructA{[]int{3, 4, 5}}}},
			`util.TestStruct{V=valV, V3=map[string]util.TestStructA{k2=util.TestStructA{v=[]int[0,1,2]},k1=util.TestStructA{v=[]int[3,4,5]}}, v2=[]util.TestType[9,8,7]}`},
	}

	for _, test := range tests {
		if got := ReflectToString(test.input, StyleLong); got != test.expected {
			t.Errorf("ReflectToString(%v), expect: %v, but got: %v", test.input, test.expected, got)
		}
	}
}

// since the map is unsorted, the last one may be different from expected
func TestReflectToStringMedium(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected string
	}{
		{true, `true`},
		{32.68, `32.68`},
		{0xFA, `250`},
		{struct{ v string }{v: "ss"}, `{v=ss}`},
		{struct{ v int }{v: -1}, `{v=-1}`},
		{struct {
			v1 float64
			V2 bool
		}{v1: 1.01, V2: true}, `{v1=1.01, V2=true}`},
		{struct {
			v1 float64
			V2 []string
		}{v1: 1.01, V2: []string{"a", "v"}}, `{v1=1.01, V2=[a,v]}`},
		{struct {
			v1 TestType
			V2 []interface{}
		}{v1: 2, V2: []interface{}{"a", 1, true, struct{ v TestType }{-1}}},
			`{v1=2, V2=[a,1,true,{v=-1}]}`},
		{TestStruct{V: "valV", v2: []TestType{9, 8, 7}, V3: map[string]TestStructA{"k2": TestStructA{[]int{0, 1, 2}}, "k1": TestStructA{[]int{3, 4, 5}}}},
			`{V=valV, V3={k2={v=[0,1,2]},k1={v=[3,4,5]}}, v2=[9,8,7]}`},
	}

	for _, test := range tests {
		if got := ReflectToString(test.input); got != test.expected {
			t.Errorf("ReflectToString(%v), expect: %v, but got: %v", test.input, test.expected, got)
		}
	}
}

// since the map is unsorted, the last one may be different from expected
func TestReflectToStringShort(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected string
	}{
		{true, `true`},
		{32.68, `32.68`},
		{0xFA, `250`},
		{[]int{1,2,3,4,5}, `[1;2;3;4;5]`},
		{struct{ v string }{v: "ss"}, `{ss}`},
		{struct{ v int }{v: -1}, `{-1}`},
		{struct {
			v1 float64
			V2 bool
		}{v1: 1.01, V2: true}, `{1.01, true}`},
		{struct {
			v1 float64
			V2 []string
		}{v1: 1.01, V2: []string{"a", "v"}}, `{1.01, [a;v]}`},
		{struct {
			v1 TestType
			V2 []interface{}
		}{v1: 2, V2: []interface{}{"a", 1, true, struct{ v TestType }{-1}}},
			`{2, [a;1;true;{-1}]}`},
		{TestStruct{V: "valV", v2: []TestType{9, 8, 7}, V3: map[string]TestStructA{"k2": TestStructA{[]int{0, 1, 2}}, "k1": TestStructA{[]int{3, 4, 5}}}},
			`{valV, {k2={[0;1;2]};k1={[3;4;5]}}, [9;8;7]}`},
	}

	for _, test := range tests {
		if got := ReflectToString(test.input, StyleShort, &Conf{SepElem:";"}); got != test.expected {
			t.Errorf("ReflectToString(%v), expect: %v, but got: %v", test.input, test.expected, got)
		}
	}
}

func TestConfig(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected string
	}{
		{true, `true`},
		{[]int{1,2,3,4,5}, `1-2-3-4-5`},
		{0xFA, `250`},
	}

	c := &Conf{SepElem:"-", BoundaryArrayAndSliceStart:NONE, BoundaryArrayAndSliceEnd:NONE}
	for _, test := range tests {
		if got := ReflectToString(test.input, StyleShort, c); got != test.expected {
			t.Errorf("ReflectToString(%v), expect: %v, but got: %v", test.input, test.expected, got)
		}
	}
}