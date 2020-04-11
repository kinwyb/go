package db

import (
	"errors"
	"testing"
)

func Test_ExecResultToBytes(t *testing.T) {
	r := &rusMsg{
		lastInsertId: 10,
		rowsAffected: 30,
		err:          errors.New("错误内容"),
	}
	data := ExecResultToBytes(r)
	r2 := BytesToExecResult(data)
	t.Log(r2)
}
