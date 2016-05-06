// Description: utils/strings.go
// Author: ZHU HAIHUA
// Since: 2016-04-08 19:45
package util

import (
	"encoding/json"
	"fmt"
	. "reflect"
	"strconv"
)

type StringStyle int

const (
	StringStyleShort StringStyle = iota
	StringStyleMedium
	StringStyleLong
)

// ToJson return the json format of the obj
// when error occur it will return empty.
// Notice: unexported field will not be marshaled
func ToJson(obj interface{}) string {
	result, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("<no value with error: %v>", err)
	}
	return string(result)
}

// ToString return the common string format of the obj according
// to the given arguments
//
// by default obj.ToString() will be called if this method exists.
// otherwise we will call ReflectToString() to get it's string
// representation
//
// the args please refer to the ReflectToString() function.
func ToString(obj interface{}, args ...interface{}) string {
	if v, ok := obj.(fmt.Stringer); ok {
		return v.String()
	}
	return ReflectToString(obj, args)
}

const (
	// the NONE used to set an empty boundary or separator.
	// e.g: you want to set the output of slice with no boundary,
	// you NEED to set StringConf as:
	//
	//      &StringConf {
	//          BoundaryArrayAndSliceStart: NONE, // NOT ""
	//          BoundaryArrayAndSliceEnd: NONE, // NOT ""
	//      }
	NONE string = "<none>"
)

type StringConf struct {
	SepElem     string `defaults:","`
	SepField    string `defaults:", "`
	SepKeyValue string `defaults:"="`

	BoundaryStructStart      string `defaults:"{"`
	BoundaryStructEnd        string `defaults:"}"`
	BoundaryMapStart         string `defaults:"{"`
	BoundaryMapEnd           string `defaults:"}"`
	BoundaryArraySliceStart  string `defaults:"["`
	BoundaryArraySliceEnd    string `defaults:"]"`
	BoundaryPointerFuncStart string `defaults:"("`
	BoundaryPointerFuncEnd   string `defaults:")"`
	BoundaryInterfaceStart   string `defaults:"("`
	BoundaryInterfaceEnd     string `defaults:")"`
}

// v return the value of StringConf by given key
func (conf *StringConf) v(name string) string {
	cnf := ValueOf(conf).Elem()
	if field, ok := cnf.Type().FieldByName(name); ok {
		value := cnf.FieldByName(name).String()
		if value == NONE {
			return ""
		} else if value == "" {
			return field.Tag.Get("defaults")
		} else {
			return cnf.FieldByName(name).String()
		}
	} else {
		panic(fmt.Errorf("<no such key: %s>", name))
	}
}

// ReflectToString return the string formatted by the given arguments,
// the number of optional arguments can be one or two
//
// the first argument is the print style, and it's default value is
// StyleMedium. the second argument is the style configuration pointer.
//
// The long style may be a very long format like following:
//
//      Type{name=value}
//
// it's some different from fmt.Printf("%#v\n", value),
// and separated by comma and equal by default.
//
// Then, the medium style would like:
//
//      {key=value}
//
// it's some different from fmt.Printf("%+v\n", value),
// and separated by comma and equal by default.
//
// Otherwise the short format will only print the value but no type
// and name information.
//
// since recursive calling, this method would be pretty slow, so if you
// use it to print log, may be you need to check if the log level is
// enabled firstly
//
// examples:
//
//   - ReflectToString(input)
//   - ReflectToString(input, StringStyleLong)
//   - ReflectToString(input, StringStyleMedium, &StringConf{SepElem:";", SepField:",", SepKeyValue:":"})
//   - ReflectToString(input, StringStyleLong, &StringConf{SepField:","})
func ReflectToString(obj interface{}, args ...interface{}) string {
	style := StringStyleMedium
	cnf := &StringConf{}
	switch len(args) {
	case 1:
		style = args[0].(StringStyle)
	case 2:
		style = args[0].(StringStyle)
		cnf = args[1].(*StringConf)
	}
	return valueToString(ValueOf(obj), style, cnf)
}

// valueToString recursively print all the value
func valueToString(val Value, style StringStyle, cnf *StringConf) string {
	var str string
	if !val.IsValid() {
		return "<zero Value>"
	}
	typ := val.Type()
	switch val.Kind() {
	case Int, Int8, Int16, Int32, Int64:
		return strconv.FormatInt(val.Int(), 10)
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return strconv.FormatUint(val.Uint(), 10)
	case Float32, Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)
	case Complex64, Complex128:
		c := val.Complex()
		return strconv.FormatFloat(real(c), 'g', -1, 64) + "+" + strconv.FormatFloat(imag(c), 'g', -1, 64) + "i"
	case String:
		return val.String()
	case Bool:
		if val.Bool() {
			return "true"
		} else {
			return "false"
		}
	case Ptr:
		v := val
		if style == StringStyleLong {
			str += typ.String() + cnf.v("BoundaryPointerFuncStart")
		}
		if v.IsNil() {
			str += "0"
		} else {
			str += "&" + valueToString(v.Elem(), style, cnf)
		}
		if style == StringStyleLong {
			str += cnf.v("BoundaryPointerFuncEnd")
		}
		return str
	case Array, Slice:
		v := val
		if style == StringStyleLong {
			str += typ.String()
		}
		str += cnf.v("BoundaryArraySliceStart")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += cnf.v("SepElem")
			}
			str += valueToString(v.Index(i), style, cnf)
		}
		str += cnf.v("BoundaryArraySliceEnd")
		return str
	case Map:
		t := typ
		if style == StringStyleLong {
			str += t.String()
		}
		str += cnf.v("BoundaryMapStart")
		//str += "<can't iterate on maps>"
		keys := val.MapKeys()
		for i, _ := range keys {
			if i > 0 {
				str += cnf.v("SepElem")
			}
			str += valueToString(keys[i], style, cnf)
			str += cnf.v("SepKeyValue")
			str += valueToString(val.MapIndex(keys[i]), style, cnf)
		}
		str += cnf.v("BoundaryMapEnd")
		return str
	case Chan:
		if style == StringStyleLong {
			str += typ.String()
		}
		return str
	case Struct:
		t := typ
		v := val
		if style == StringStyleLong {
			str += t.String()
		}
		str += cnf.v("BoundaryStructStart")
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += cnf.v("SepField")
			}
			if style == StringStyleLong || style == StringStyleMedium {
				str += val.Type().Field(i).Name
				str += cnf.v("SepKeyValue")
			}
			str += valueToString(v.Field(i), style, cnf)
		}
		str += cnf.v("BoundaryStructEnd")
		return str
	case Interface:
		//t := ""
		if style == StringStyleLong {
			str += typ.String() + cnf.v("BoundaryInterfaceStart")
		}
		str += valueToString(val.Elem(), style, cnf)
		if style == StringStyleLong {
			str += cnf.v("BoundaryInterfaceEnd")
		}
		return str
	case Func:
		v := val
		if style == StringStyleLong {
			str += typ.String() + cnf.v("BoundaryPointerFuncStart")
		}
		str += strconv.FormatUint(uint64(v.Pointer()), 10)
		if style == StringStyleLong {
			str += cnf.v("BoundaryPointerFuncEnd")
		}
		return str
	default:
		panic("valueToString: can't print type " + typ.String())
	}
}
