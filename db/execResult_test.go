package db

import (
	"testing"

	"github.com/kinwyb/go/err1"
)

func Test_ExecResultToBytes(t *testing.T) {
	r := &rusMsg{
		lastInsertId: 10,
		rowsAffected: 30,
		err:          err1.NewError(104, "错误内容"),
	}
	data := ExecResultToBytes(r)
	r2 := BytesToExecResult(data)
	t.Log(r2)
}
