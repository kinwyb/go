package conv

import (
	"testing"
)

func Test_num(t *testing.T) {
	v := ToInt64("000002220000")
	println(v)
	v = ToInt64("000002220000.655")
	println(v)
	v1 := ToFloat64("000002220000.000")
	println(v1)
	v2 := []string{"哈哈", "jhehe"}
	v3, e := ToSliceE(v2)
	if e != nil {
		t.Fatal(e)
	}
	println(v3)
}
