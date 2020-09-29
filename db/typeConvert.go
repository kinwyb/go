package db

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/kinwyb/go/conv"
	"github.com/micro/protobuf/ptypes"
	"reflect"
	"strings"
	"time"
)

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    ExecResult
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
// Deprecated: use conv package function
func Int(reply interface{}) (int, error) {
	return conv.ToIntE(reply)
}

// IntDefault is a helper that converts interface to an integer. If err is not
// equal to nil, then return default value
// Deprecated: use conv package function
func IntDefault(reply interface{}, def ...int) int {
	result, err := conv.ToIntE(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    ExecResult
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
// Deprecated: use conv package function
func Int64(reply interface{}) (int64, error) {
	return conv.ToInt64E(reply)
}

// Int64Default is a helper that converts 64 bit integer. If err is
// not equal to nil, then return default value
// Deprecated: use conv package function
func Int64Default(reply interface{}, def ...int64) int64 {
	result, err := conv.ToInt64E(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// Uint64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    ExecResult
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
// Deprecated: use conv package function
func Uint64(reply interface{}) (uint64, error) {
	return conv.ToUint64E(reply)
}

// Uint64Default is a helper that converts to 64 bit integer. If err is
// not equal to nil, then return default value
// Deprecated: use conv package function
func Uint64Default(reply interface{}, def ...uint64) uint64 {
	result, err := conv.ToUint64E(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to an int as follows:
//
//  Reply type    ExecResult
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
// Deprecated: use conv package function
func Float64(reply interface{}) (float64, error) {
	return conv.ToFloat64E(reply)
}

// Float64Default is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then return default value
// Deprecated: use conv package function
func Float64Default(reply interface{}, def ...float64) float64 {
	result, err := conv.ToFloat64E(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      ExecResult
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
// Deprecated: use conv package function
func String(reply interface{}) (string, error) {
	if v, ok := reply.(time.Time); ok {
		return v.Format("2006-01-02 15:04:05"), nil
	}
	return conv.ToStringE(reply)
}

// StringDefault is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns default value
// Deprecated: use conv package function
func StringDefault(reply interface{}, def ...string) string {
	result, err := String(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return strings.TrimSpace(result)
}

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      ExecResult
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
// Deprecated: use conv package function
func Bytes(reply interface{}) ([]byte, error) {
	return conv.ToBytesE(reply)
}

// BytesDefault is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns default value
// Deprecated: use conv package function
func BytesDefault(reply interface{}, def ...[]byte) []byte {
	result, err := conv.ToBytesE(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      ExecResult
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
// Deprecated: use conv package function
func Bool(reply interface{}) (bool, error) {
	return conv.ToBoolE(reply)
}

// BoolDefault is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns default value
// Deprecated: use conv package function
func BoolDefault(reply interface{}, def ...bool) bool {
	result, err := conv.ToBoolE(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
// Deprecated: use conv package function
func Strings(reply interface{}) ([]string, error) {
	return conv.ToStringSliceE(reply)
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
// Deprecated: use conv package function
func StringMap(result interface{}) (map[string]string, error) {
	return conv.ToStringMapStringE(result)
}

// IntMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]int. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
// Deprecated: use conv package function
func IntMap(result interface{}) (map[string]int, error) {
	return conv.ToStringMapIntE(result)
}

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
// Deprecated: use conv package function
func Int64Map(result interface{}) (map[string]int64, error) {
	return conv.ToStringMapInt64E(result)
}

// ToValue converts an interface{} to a ptypes.Value
func InterfaceToPytesStruceValue(v interface{}) *st.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	default:
		// Fallback to reflection for other types
		return toValue(reflect.ValueOf(v))
	}
}

func toValue(v reflect.Value) *st.Value {
	switch v.Kind() {
	case reflect.Bool:
		return &st.Value{
			Kind: &st.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return toValue(reflect.Indirect(v))
	case reflect.Array, reflect.Slice:
		size := v.Len()
		if size == 0 {
			return nil
		}
		values := make([]*st.Value, size)
		for i := 0; i < size; i++ {
			values[i] = toValue(v.Index(i))
		}
		return &st.Value{
			Kind: &st.Value_ListValue{
				ListValue: &st.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = toValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*st.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = toValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &st.Value{
			Kind: &st.Value_StructValue{
				StructValue: &st.Struct{
					Fields: fields,
				},
			},
		}
	default:
		// Last resort
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}

// DecodeToMap converts a pb.Struct to a map from strings to Go types.
// DecodeToMap panics if s is invalid.
func DecodePytesStruceValueToMap(s *st.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	m := map[string]interface{}{}
	for k, v := range s.Fields {
		m[k] = decodeValue(v)
	}
	return m
}

func decodeValue(v *st.Value) interface{} {
	switch k := v.Kind.(type) {
	case *st.Value_NullValue:
		return nil
	case *st.Value_NumberValue:
		return k.NumberValue
	case *st.Value_StringValue:
		return k.StringValue
	case *st.Value_BoolValue:
		return k.BoolValue
	case *st.Value_StructValue:
		return DecodePytesStruceValueToMap(k.StructValue)
	case *st.Value_ListValue:
		s := make([]interface{}, len(k.ListValue.Values))
		for i, e := range k.ListValue.Values {
			s[i] = decodeValue(e)
		}
		return s
	default:
		panic("protostruct: unknown kind")
	}
}

func InterfaceToProtoAny(v interface{}) (*any.Any, error) {
	return ptypes.MarshalAny(InterfaceToPytesStruceValue(v))
}

func InterfaceToProtoAnyDefault(v interface{}) *any.Any {
	ret, _ := ptypes.MarshalAny(InterfaceToPytesStruceValue(v))
	return ret
}

func ProtoAnyToInterface(v *any.Any) interface{} {
	s := &st.Value{}
	ptypes.UnmarshalAny(v, s)
	return decodeValue(s)
}
