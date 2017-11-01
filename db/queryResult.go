package db

import (
	"database/sql"

	perr "github.com/kinwyb/go/err"
)

//查询结果数据存放对象
//查询结果对象
type QueryResult struct {
	columns    []string
	data       [][]interface{}
	datalength int
	rows       *sql.Rows
	res        *res
}

//读取查询结果
func (q *QueryResult) setResult(rows *sql.Rows) (err error) {
	if rows == nil {
		q.columns = []string{}
		q.data = nil
		return
	}
	q.columns, err = rows.Columns()
	q.rows = rows
	return
}

//解析查询结果
func (q *QueryResult) passRows() {
	if q.rows != nil {
		q.data = make([][]interface{}, 0)
		for q.rows.Next() {
			row := make([]interface{}, len(q.columns))
			for i := range row {
				var ref interface{}
				row[i] = &ref
			}
			err := q.rows.Scan(row...)
			if err != nil {
				q.datalength = 0
				q.rows.Close()
				q.rows = nil
				if q.res != nil {
					if q.res.errFmt != nil {
						q.res.err = q.res.errFmt.FormatError(err)
					} else {
						q.res.err = perr.NewError(1, "数据读取错误", err)
					}
				}
				return
			}
			for k, v := range row {
				row[k] = *v.(*interface{})
			}
			q.data = append(q.data, row)
		}
		q.datalength = len(q.data)
		q.rows.Close()
		q.rows = nil
	}
}

//读取某行的指定字段值.
//columnName表示字段名称，index表示第几行默认第一行，如果结果不存在返回nil
func (q *QueryResult) Get(columnName string, index ...int) interface{} {
	if len(index) < 1 {
		index = []int{0}
	}
	q.passRows()
	if index[0] >= q.datalength { //超出数据返回nil
		return nil
	}
	for i, v := range q.columns {
		if v == columnName {
			return q.data[index[0]][i]
		}
	}
	return nil
}

//读取某行的所有数据.
//index代表第几行默认第一行，返回的map中key是数据字段名称，value是值
func (q *QueryResult) GetMap(index ...int) (ret map[string]interface{}) {
	if len(index) < 1 {
		index = []int{0}
	}
	q.passRows()
	if index[0] >= q.datalength {
		return
	}
	ret = make(map[string]interface{})
	for i, v := range q.columns {
		ret[v] = q.data[index[0]][i]
	}
	return
}

//获取字段列表
func (q *QueryResult) Columns() []string {
	return q.columns
}

//获取所有数据
func (q *QueryResult) Rows() [][]interface{} {
	q.passRows()
	return q.data
}

//获取结果长度
func (q *QueryResult) Length() int {
	q.passRows()
	return q.datalength
}

//循环读取所有数据
//返回的map中key是数据字段名称，value是值,回调函数中如果返回false则停止循环后续数据
func (q *QueryResult) ForEach(f func(map[string]interface{}) bool) {
	if f == nil {
		return
	}
	q.passRows()
	ret := map[string]interface{}{}
	for j, v := range q.data {
		if j >= q.datalength {
			break
		}
		for i, vv := range q.columns {
			ret[vv] = v[i]
		}
		if !f(ret) {
			break
		}
	}
}

//关闭查询结果
func (q *QueryResult) Close() {
	if q.rows != nil {
		q.rows.Close()
	}
}
