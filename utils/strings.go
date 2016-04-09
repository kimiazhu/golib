// Description: utils/strings.go
// Author: ZHU HAIHUA
// Since: 2016-04-08 19:45
package util

import (
	"encoding/json"
	"strconv"
	"fmt"
	. "reflect"
)

type Style int

const (
	StyleShort Style = iota
	StyleMedium
	StyleLong
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

// ReflectToString return the string format of the given argument,
// the default style is StyleMedium
//
// the long style may be a very long format like following:
//
//      Type{name=value}
//
// and the medium style would like:
//
//      {key=value}
//
// otherwise the short format will only print the value but no type
// and name information.
//
// since recursive call, this method would be pretty slow, so if you
// use it to print log, may be you need to check if the log level is
// enabled first
func ReflectToString(obj interface{}, args ...Style) string {
	style := StyleMedium
	if len(args) > 0 {
		style = args[0]
	}

	var result string
	switch style {
	case StyleShort:
		result = fmt.Sprintf("%v", obj)
	case StyleMedium, StyleLong:
		result = valueToString(ValueOf(obj), style)
	}
	return result
}

// valueToString recursively print all the value
func valueToString(val Value, style Style) string {
	if style == StyleShort {
		return "<not suitable for short style>"
	}
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
		if style == StyleLong {
			str = typ.String() + "("
		} else {
			str = "("
		}
		if v.IsNil() {
			str += "0"
		} else {
			str += "&" + valueToString(v.Elem(), style)
		}
		str += ")"
		return str
	case Array, Slice:
		v := val
		if style == StyleLong {
			str += typ.String()
		}
		str += "{"
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += ", "
			}
			str += valueToString(v.Index(i), style)
		}
		str += "}"
		return str
	case Map:
		t := typ
		if style == StyleLong{
			str = t.String()
		}
		str += "{"
		//str += "<can't iterate on maps>"
		keys := val.MapKeys()
		for i, _ := range keys {
			if i > 0 {
				str += ","
			}
			str += valueToString(keys[i], style)
			str += "="
			str += valueToString(val.MapIndex(keys[i]), style)
		}
		str += "}"
		return str
	case Chan:
		if style == StyleLong {
			str = typ.String()
		}
		return str
	case Struct:
		t := typ
		v := val
		if style == StyleLong {
			str += t.String()
		}
		str += "{"
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += ", "
			}
			str += val.Type().Field(i).Name
			str += "="
			str += valueToString(v.Field(i), style)
		}
		str += "}"
		return str
	case Interface:
		t := ""
		if style == StyleLong{
			t = typ.String()
		}
		return t + "(" + valueToString(val.Elem(), style) + ")"
	case Func:
		v := val
		t := ""
		if style == StyleLong{
			t = typ.String()
		}
		return t + "(" + strconv.FormatUint(uint64(v.Pointer()), 10) + ")"
	default:
		panic("valueToString: can't print type " + typ.String())
	}
}