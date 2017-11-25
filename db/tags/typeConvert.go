package tags

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kinwyb/go/err1"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("nil returned")

var errNegativeInt = errors.New("unexpected value for Uint64")

//interface指针转换成interface
func convertInterfacePoint(it interface{}) interface{} {
	switch it.(type) {
	case *interface{}:
		return *it.(*interface{})
	default:
		return it
	}
}

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    ExecResult
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface{}) (int, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case int:
		return reply, nil
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case string:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return int(n), err
	case nil:
		return 0, ErrNil
	case err1.Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int, got type %T", reply)
}

// IntDefault is a helper that converts interface to an integer. If err is not
// equal to nil, then return default value
func IntDefault(reply interface{}, def ...int) int {
	result, err := Int(reply)
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
func Int64(reply interface{}) (int64, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case int:
		return int64(reply), nil
	case int64:
		return reply, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case string:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case err1.Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int64, got type %T", reply)
}

// Int64Default is a helper that converts 64 bit integer. If err is
// not equal to nil, then return default value
func Int64Default(reply interface{}, def ...int64) int64 {
	result, err := Int64(reply)
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
func Uint64(reply interface{}) (uint64, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case int64:
		if reply < 0 {
			return 0, errNegativeInt
		}
		return uint64(reply), nil
	case []byte:
		n, err := strconv.ParseUint(string(reply), 10, 64)
		return n, err
	case string:
		n, err := strconv.ParseUint(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case err1.Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Uint64, got type %T", reply)
}

// Uint64Default is a helper that converts to 64 bit integer. If err is
// not equal to nil, then return default value
func Uint64Default(reply interface{}, def ...uint64) uint64 {
	result, err := Uint64(reply)
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
func Float64(reply interface{}) (float64, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []byte:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case string:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case float64:
		return float64(reply), nil
	case nil:
		return 0, ErrNil
	case err1.Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Float64, got type %T", reply)
}

// Float64Default is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then return default value
func Float64Default(reply interface{}, def ...float64) float64 {
	result, err := Float64(reply)
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
func String(reply interface{}) (string, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case int64:
		return strconv.FormatInt(reply, 10), nil
	case int32:
		return strconv.FormatInt(int64(reply), 10), nil
	case nil:
		return "", ErrNil
	case err1.Error:
		return "", reply
	}
	return "", fmt.Errorf("redigo: unexpected type for String, got type %T", reply)
}

// StringDefault is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns default value
func StringDefault(reply interface{}, def ...string) string {
	result, err := String(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
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
func Bytes(reply interface{}) ([]byte, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []byte:
		return reply, nil
	case string:
		return []byte(reply), nil
	case nil:
		return nil, ErrNil
	case err1.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Bytes, got type %T", reply)
}

// BytesDefault is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns default value
func BytesDefault(reply interface{}, def ...[]byte) []byte {
	result, err := Bytes(reply)
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
func Bool(reply interface{}) (bool, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case int64:
		return reply != 0, nil
	case []byte:
		return strconv.ParseBool(string(reply))
	case nil:
		return false, ErrNil
	case err1.Error:
		return false, reply
	}
	return false, fmt.Errorf("redigo: unexpected type for Bool, got type %T", reply)
}

// BoolDefault is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns default value
func BoolDefault(reply interface{}, def ...bool) bool {
	result, err := Bool(reply)
	if err != nil && def != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

// MultiBulk is a helper that converts an array command reply to a []interface{}.
//
// Deprecated: Use Values instead.
func MultiBulk(reply interface{}) ([]interface{}, error) { return Values(reply) }

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      ExecResult
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface{}) ([]interface{}, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []interface{}:
		return reply, nil
	case nil:
		return nil, ErrNil
	case err1.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Values, got type %T", reply)
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface{}) ([]string, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []interface{}:
		result := make([]string, len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			p, ok := reply[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("redigo: unexpected element type for Strings, got type %T", reply[i])
			}
			result[i] = string(p)
		}
		return result, nil
	case nil:
		return nil, ErrNil
	case err1.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Strings, got type %T", reply)
}

// ByteSlices is a helper that converts an array command reply to a [][]byte.
// If err is not equal to nil, then ByteSlices returns nil, err. Nil array
// items are stay nil. ByteSlices returns an error if an array item is not a
// bulk string or nil.
func ByteSlices(reply interface{}) ([][]byte, error) {
	switch reply := convertInterfacePoint(reply).(type) {
	case []interface{}:
		result := make([][]byte, len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			p, ok := reply[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("redigo: unexpected element type for ByteSlices, got type %T", reply[i])
			}
			result[i] = p
		}
		return result, nil
	case nil:
		return nil, ErrNil
	}
	return nil, fmt.Errorf("redigo: unexpected type for ByteSlices, got type %T", reply)
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(result interface{}) (map[string]string, error) {
	values, err := Values(result)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		m[string(key)] = string(value)
	}
	return m, nil
}

// IntMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]int. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func IntMap(result interface{}) (map[string]int, error) {
	values, err := Values(result)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: IntMap expects even number of values result")
	}
	m := make(map[string]int, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		value, err := Int(values[i+1])
		if err != nil {
			return nil, err
		}
		m[string(key)] = value
	}
	return m, nil
}

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func Int64Map(result interface{}) (map[string]int64, error) {
	values, err := Values(result)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: Int64Map expects even number of values result")
	}
	m := make(map[string]int64, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		value, err := Int64(values[i+1])
		if err != nil {
			return nil, err
		}
		m[string(key)] = value
	}
	return m, nil
}
