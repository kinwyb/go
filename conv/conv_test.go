package conv

import "testing"

func Test_num(t *testing.T) {
	v := ToInt64("000002220000")
	println(v)
	v = ToInt64("000002220000.655")
	println(v)
	v1 := ToFloat64("000002220000.000")
	println(v1)
}
