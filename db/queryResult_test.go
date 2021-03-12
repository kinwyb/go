package db

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestQueryResultToByte(t *testing.T) {
	r := &res{
		columns: []string{"string",
			"[]byte", "bool", "float64", "float32",
			"int64", "int32"},
		datalength: 2,
		err:        errors.New("错误内容"),
		data: [][]interface{}{
			{
				"字符串1",
				[]byte("bytes"),
				false,
				10.23,
				1.2,
				64,
				1,
			},
			{
				"字符串2",
				[]byte("bytes2"),
				true,
				5.1,
				8.3,
				6,
				19,
			},
		},
	}
	b := QueryResultToBytes(r)
	fmt.Println(hex.EncodeToString(b))
	r2 := BytesToQueryResult(b)
	fmt.Println(r2.Columns())
	fmt.Println(r2.Length())
	it := r.Iterator()
	for it.HasNext() {
		data, err := it.Next()
		if err == io.EOF {
			break
		}
		fmt.Printf("%v\n", data)
	}
}
