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
	commaAndSpace = ", "
	comma         = ","
	equals        = "="

	// the NONE used to set an empty boundary or separator.
	// e.g: you want to set the output of slice with no boundary,
	// you NEED to set Conf as:
	//
	//      &Conf {
	//          BoundaryArrayAndSliceStart: NONE, // NOT ""
	//          BoundaryArrayAndSliceEnd: NONE, // NOT ""
	//      }
	NONE = "<none>"
)

type Conf struct {
	SepElem     string
	SepField    string
	SepKeyValue string

	BoundaryStructStart        string
	BoundaryStructEnd          string
	BoundaryMapStart           string
	BoundaryMapEnd             string
	BoundaryArrayAndSliceStart string
	BoundaryArrayAndSliceEnd   string
	BoundaryPointerFuncStart   string
	BoundaryPointerFuncEnd     string
	BoundaryInterfaceStart     string
	BoundaryInterfaceEnd       string
}

var global *Conf = &Conf{
	SepElem:     comma,
	SepField:    commaAndSpace,
	SepKeyValue: equals,

	BoundaryStructStart:        "{",
	BoundaryStructEnd:          "}",
	BoundaryMapStart:           "{",
	BoundaryMapEnd:             "}",
	BoundaryArrayAndSliceStart: "[",
	BoundaryArrayAndSliceEnd:   "]",
	BoundaryPointerFuncStart:   "(",
	BoundaryPointerFuncEnd:     ")",
	BoundaryInterfaceStart:     "(",
	BoundaryInterfaceEnd:       ")",
}

// updateConfig will update the global configuration using the given conf.
// any empty config field will be ignored, so if you really want to set a
// config value to be an empty string, you need to set it to be const NONE, l
// ike this:
//
//      &Conf {
//          BoundaryArrayAndSliceStart: NONE, // DO NOT SET TO: ""
//          BoundaryArrayAndSliceEnd: NONE, // DO NOT SET TO: ""
//      }
func updateConfig(conf *Conf) {
	arg := ValueOf(conf).Elem()
	gcnf := ValueOf(global).Elem()
	for i := 0; i < arg.NumField(); i++ {
		av := arg.Field(i).String()
		if av == NONE {
			gcnf.Field(i).SetString("")
		} else if av != "" {
			gcnf.Field(i).SetString(av)
		}
	}
}

// ReflectToString return the string formatted by the given argument,
// the args may be one or two
//
// the first argument is the print style, and it's default value is
// StyleMedium. the second argument is the style configuration pointer.
//
// The long style may be a very long format like following:
//
//      Type{name=value}
//
// it's some different from fmt.Printf("%#v\n", value),
// it's separated by comma and equal
//
// Then, the medium style would like:
//
//      {key=value}
//
// it's some different from fmt.Printf("%+v\n", value),
// it's separated by comma and equal
//
// Otherwise the short format will only print the value but no type
// and name information.
//
// since recursive call, this method would be pretty slow, so if you
// use it to print log, may be you need to check if the log level is
// enabled first
//
// examples:
//
//   - ReflectToString(input)
//   - ReflectToString(input, StyleLong)
//   - ReflectToString(input, StyleMedium, &Conf{SepElem:";", SepField:",", SepKeyValue:":"})
//   - ReflectToString(input, StyleLong, &Conf{SepField:","})
func ReflectToString(obj interface{}, args ...interface{}) string {
	style := StyleMedium
	switch len(args) {
	case 1:
		style = args[0].(Style)
	case 2:
		style = args[0].(Style)
		cnf := args[1].(*Conf)
		updateConfig(cnf)
	}

	return valueToString(ValueOf(obj), style)
}

// valueToString recursively print all the value
func valueToString(val Value, style Style) string {
	//if style == StyleShort {
	//	return "<not suitable for short style>"
	//}
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
			str += typ.String() + global.BoundaryPointerFuncStart
		}
		if v.IsNil() {
			str += "0"
		} else {
			str += "&" + valueToString(v.Elem(), style)
		}
		if style == StyleLong {
			str += global.BoundaryPointerFuncEnd
		}
		return str
	case Array, Slice:
		v := val
		if style == StyleLong {
			str += typ.String()
		}
		str += global.BoundaryArrayAndSliceStart
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str += global.SepElem
			}
			str += valueToString(v.Index(i), style)
		}
		str += global.BoundaryArrayAndSliceEnd
		return str
	case Map:
		t := typ
		if style == StyleLong {
			str += t.String()
		}
		str += global.BoundaryMapStart
		//str += "<can't iterate on maps>"
		keys := val.MapKeys()
		for i, _ := range keys {
			if i > 0 {
				str += global.SepElem
			}
			str += valueToString(keys[i], style)
			str += global.SepKeyValue
			str += valueToString(val.MapIndex(keys[i]), style)
		}
		str += global.BoundaryMapEnd
		return str
	case Chan:
		if style == StyleLong {
			str += typ.String()
		}
		return str
	case Struct:
		t := typ
		v := val
		if style == StyleLong {
			str += t.String()
		}
		str += global.BoundaryStructStart
		for i, n := 0, v.NumField(); i < n; i++ {
			if i > 0 {
				str += global.SepField
			}
			if style == StyleLong || style == StyleMedium {
				str += val.Type().Field(i).Name
				str += global.SepKeyValue
			}
			str += valueToString(v.Field(i), style)
		}
		str += global.BoundaryStructEnd
		return str
	case Interface:
		//t := ""
		if style == StyleLong {
			str += typ.String() + global.BoundaryInterfaceStart
		}
		str += valueToString(val.Elem(), style)
		if style == StyleLong {
			str += global.BoundaryInterfaceEnd
		}
		return str
	case Func:
		v := val
		if style == StyleLong {
			str += typ.String() + global.BoundaryPointerFuncStart
		}
		str += strconv.FormatUint(uint64(v.Pointer()), 10)
		if style == StyleLong {
			str += global.BoundaryPointerFuncEnd
		}
		return str
	default:
		panic("valueToString: can't print type " + typ.String())
	}
}
