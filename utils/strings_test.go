// Description: util
// Author: ZHU HAIHUA
// Since: 2016-04-08 22:33
package util

import (
	"testing"
)

type TestType int

type TestStruct struct {
	V string
	V3 map[string]TestStructA
	v2 []TestType
}

type TestStructA struct {
	v []int
}

func TestReflectToString(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected string
	}{
		{true, `true`},
		{32.68, `32.68`},
		{0xFA, `250`},
		{struct{ v string }{v: "ss"}, `struct { v string }{v=ss}`},
		{struct{ v int }{v: -1}, `struct { v int }{v=-1}`},
		{struct{ v1 float64; V2 bool; }{v1: 1.01, V2: true}, `struct { v1 float64; V2 bool }{v1=1.01, V2=true}`},
		{struct{ v1 float64; V2 []string; }{v1: 1.01, V2: []string{"a", "v"}}, `struct { v1 float64; V2 []string }{v1=1.01, V2=[]string{a, v}}`},
		{struct{ v1 TestType; V2 []interface{}; }{v1: 2, V2: []interface{}{"a", 1, true, struct {v TestType}{-1}}},
			`struct { v1 util.TestType; V2 []interface {} }{v1=2, V2=[]interface {}{interface {}(a), interface {}(1), interface {}(true), interface {}(struct { v util.TestType }{v=-1})}}`},
		{TestStruct{V: "valV", v2: []TestType{9,8,7}, V3: map[string]TestStructA{"k2": TestStructA{[]int{0,1,2}}, "k1": TestStructA{[]int{3,4,5}}}},
			`util.TestStruct{V=valV, V3=map[string]util.TestStructA{k2=util.TestStructA{v=[]int{0, 1, 2}},k1=util.TestStructA{v=[]int{3, 4, 5}}}, v2=[]util.TestType{9, 8, 7}}`},
	}

	for _, test := range tests {
		if got := ReflectToString(test.input, StyleLong); got != test.expected {
			t.Errorf("ReflectToString(%v), expect: %v, but got: %v", test.input, test.expected, got)
		}
	}
}