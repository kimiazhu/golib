// Description: utils/strings
// Author: ZHU HAIHUA
// Since: 2016-04-08 19:45
package strings

import (
	"encoding/json"
	"strconv"
	"fmt"
	. "reflect"
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

func ReflectToString(obj interface{}) string {
	fmt.Printf("========%v\n", ValueOf(obj).Type().Field(2).Name)
	return valueToString(ValueOf(obj))
}

func valueToString(val Value) string {

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
		str = typ.String() + "("
		if v.IsNil() {
			str += "0"
		} else {
			str += "&" + valueToString(v.Elem())
		}
		str += ")"
		return str
	case Array, Slice:
		v := val
		str += typ.String()
		str += "{"
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += ", "
			}
			str += valueToString(v.Index(i))
		}
		str += "}"
		return str
	case Map:
		t := typ
		str = t.String()
		str += "{"
		str += "<can't iterate on maps>"
		str += "}"
		return str
	case Chan:
		str = typ.String()
		return str
	case Struct:
		t := typ
		v := val
		str += t.String()
		str += "{"
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += ", "
			}
			str += valueToString(v.Field(i))
		}
		str += "}"
		return str
	case Interface:
		return typ.String() + "(" + valueToString(val.Elem()) + ")"
	case Func:
		v := val
		return typ.String() + "(" + strconv.FormatUint(uint64(v.Pointer()), 10) + ")"
	default:
		panic("valueToString: can't print type " + typ.String())
	}
}