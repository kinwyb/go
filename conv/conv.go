package conv

import (
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// ToBool casts an interface to a bool type.
func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

// ToTime casts an interface to a time.Time type.
func ToTime(i interface{}) time.Time {
	v, _ := ToTimeE(i)
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i interface{}) float64 {
	v, _ := ToFloat64E(i)
	return v
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i interface{}) float32 {
	v, _ := ToFloat32E(i)
	return v
}

// ToInt64 casts an interface to an int64 type.
func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt32 casts an interface to an int32 type.
func ToInt32(i interface{}) int32 {
	v, _ := ToInt32E(i)
	return v
}

// ToInt16 casts an interface to an int16 type.
func ToInt16(i interface{}) int16 {
	v, _ := ToInt16E(i)
	return v
}

// ToInt8 casts an interface to an int8 type.
func ToInt8(i interface{}) int8 {
	v, _ := ToInt8E(i)
	return v
}

// ToInt casts an interface to an int type.
func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

// ToUint casts an interface to a uint type.
func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

// ToString casts an interface to a string type.
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

func ToBytes(i interface{}) []byte {
	v, _ := ToBytesE(i)
	return v
}

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapString(i interface{}) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i interface{}) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i interface{}) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapInt(i interface{}) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64(i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i interface{}) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i interface{}) []int {
	v, _ := ToIntSliceE(i)
	return v
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}

func checkByteArray(i interface{}, toNum bool) interface{} {
	if v, ok := indirect(i).([]byte); ok {
		i = string(v)
	}
	if toNum {
		if v, ok := i.(string); ok {
			if !strings.HasPrefix(v, "0x") &&
				!strings.HasPrefix(v, "0b") &&
				!strings.HasPrefix(v, "0o") {
				v = strings.TrimLeft(v, "0")
			}
			v = strings.TrimSpace(v)
			v1, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return v1
			}
		}
	}
	return i
}

// ToBool casts an interface to a bool type.
func ToBoolE(i interface{}) (bool, error) {
	i = checkByteArray(i, false)
	return cast.ToBoolE(i)
}

// ToTime casts an interface to a time.Time type.
func ToTimeE(i interface{}) (time.Time, error) {
	i = checkByteArray(i, false)
	return cast.ToTimeE(i)
}

// ToDuration casts an interface to a time.Duration type.
func ToDurationE(i interface{}) (time.Duration, error) {
	i = checkByteArray(i, true)
	return cast.ToDurationE(i)
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64E(i interface{}) (float64, error) {
	i = checkByteArray(i, true)
	return cast.ToFloat64E(i)
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32E(i interface{}) (float32, error) {
	i = checkByteArray(i, true)
	return cast.ToFloat32E(i)
}

// ToInt64 casts an interface to an int64 type.
func ToInt64E(i interface{}) (int64, error) {
	i = checkByteArray(i, true)
	return cast.ToInt64E(i)
}

// ToInt32 casts an interface to an int32 type.
func ToInt32E(i interface{}) (int32, error) {
	i = checkByteArray(i, true)
	return cast.ToInt32E(i)
}

// ToInt16 casts an interface to an int16 type.
func ToInt16E(i interface{}) (int16, error) {
	i = checkByteArray(i, true)
	return cast.ToInt16E(i)
}

// ToInt8 casts an interface to an int8 type.
func ToInt8E(i interface{}) (int8, error) {
	i = checkByteArray(i, true)
	return cast.ToInt8E(i)
}

// ToInt casts an interface to an int type.
func ToIntE(i interface{}) (int, error) {
	i = checkByteArray(i, true)
	return cast.ToIntE(i)
}

// ToUint casts an interface to a uint type.
func ToUintE(i interface{}) (uint, error) {
	i = checkByteArray(i, true)
	return cast.ToUintE(i)
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64E(i interface{}) (uint64, error) {
	i = checkByteArray(i, true)
	return cast.ToUint64E(i)
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32E(i interface{}) (uint32, error) {
	i = checkByteArray(i, true)
	return cast.ToUint32E(i)
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16E(i interface{}) (uint16, error) {
	i = checkByteArray(i, true)
	return cast.ToUint16E(i)
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8E(i interface{}) (uint8, error) {
	i = checkByteArray(i, true)
	return cast.ToUint8E(i)
}

// ToString casts an interface to a string type.
func ToStringE(i interface{}) (string, error) {
	return cast.ToStringE(i)
}

func ToBytesE(i interface{}) ([]byte, error) {
	v, err := cast.ToStringE(i)
	if err != nil {
		return nil, err
	}
	return []byte(v), err
}

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapStringE(i interface{}) (map[string]string, error) {
	return cast.ToStringMapStringE(i)
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSliceE(i interface{}) (map[string][]string, error) {
	return cast.ToStringMapStringSliceE(i)
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBoolE(i interface{}) (map[string]bool, error) {
	return cast.ToStringMapBoolE(i)
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapIntE(i interface{}) (map[string]int, error) {
	return cast.ToStringMapIntE(i)
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64E(i interface{}) (map[string]int64, error) {
	return cast.ToStringMapInt64E(i)
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMapE(i interface{}) (map[string]interface{}, error) {
	return cast.ToStringMapE(i)
}

// ToSlice casts an interface to a []interface{} type.
func ToSliceE(i interface{}) ([]interface{}, error) {
	ret, err := cast.ToSliceE(i)
	if err == nil {
		return ret, nil
	}
	switch v := i.(type) {
	case []string:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []int64:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []int:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []float64:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []float32:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []int32:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	case []bool:
		for _, v1 := range v {
			ret = append(ret, v1)
		}
		return ret, nil
	default:
		return ret, err
	}
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSliceE(i interface{}) ([]bool, error) {
	return cast.ToBoolSliceE(i)
}

// ToStringSlice casts an interface to a []string type.
func ToStringSliceE(i interface{}) ([]string, error) {
	return cast.ToStringSliceE(i)
}

// ToIntSlice casts an interface to a []int type.
func ToIntSliceE(i interface{}) ([]int, error) {
	return cast.ToIntSliceE(i)
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSliceE(i interface{}) ([]time.Duration, error) {
	return cast.ToDurationSliceE(i)
}

// ToBool casts an interface to a bool type.
func ToBoolDefault(i interface{}, defaultValue bool) bool {
	v, err := cast.ToBoolE(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDurationDefault(i interface{}, defaultValue time.Duration) time.Duration {
	v, err := ToDurationE(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64Default(i interface{}, defaultValue float64) float64 {
	v, err := ToFloat64E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32Default(i interface{}, defaultValue float32) float32 {
	v, err := ToFloat32E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt64 casts an interface to an int64 type.
func ToInt64Default(i interface{}, defaultValue int64) int64 {
	v, err := ToInt64E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt32 casts an interface to an int32 type.
func ToInt32Default(i interface{}, defaultValue int32) int32 {
	v, err := ToInt32E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt16 casts an interface to an int16 type.
func ToInt16Default(i interface{}, defaultValue int16) int16 {
	v, err := ToInt16E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt8 casts an interface to an int8 type.
func ToInt8Default(i interface{}, defaultValue int8) int8 {
	v, err := ToInt8E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToInt casts an interface to an int type.
func ToIntDefault(i interface{}, defaultValue int) int {
	v, err := ToIntE(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToUint casts an interface to a uint type.
func ToUintDefault(i interface{}, defaultValue uint) uint {
	v, err := ToUintE(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64Default(i interface{}, defaultValue uint64) uint64 {
	v, err := ToUint64E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32Default(i interface{}, defaultValue uint32) uint32 {
	v, err := ToUint32E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16Default(i interface{}, defaultValue uint16) uint16 {
	v, err := ToUint16E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8Default(i interface{}, defaultValue uint8) uint8 {
	v, err := ToUint8E(i)
	if err != nil {
		return defaultValue
	}
	return v
}

// ToString casts an interface to a string type.
func ToStringDefault(i interface{}, defaultValue string) string {
	v, err := ToStringE(i)
	if err != nil {
		return defaultValue
	}
	return v
}
