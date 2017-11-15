package db

import (
	"database/sql"

	perr "github.com/kinwyb/go/err"
)

//查询结果返回接口
type QueryResult interface {
	//逐条获取结果
	//如果参数func返回true，并且还有下一条结果则再次调用func返回下一条
	ForEach(func(map[string]interface{}) bool) QueryResult
	//出错时回调参数方法
	Error(func(perr.Error)) QueryResult
	//是否出错
	HasError() perr.Error
	//是否为空
	IsEmpty() bool
	//结果空是回调参数方法
	Empty(func()) QueryResult
	//关闭查询结果.
	//如果读取了结果内容查询会自动关闭,只有不需要获取查询结果的时候才需要手动调用关闭查询结果
	Close()
	//获取字段列表
	Columns() []string
	//获取所有数据
	Rows() [][]interface{}
	//读取某行的指定字段值.
	//columnName表示字段名称，index表示第几行默认第一行，如果结果不存在返回nil
	Get(columnName string, index ...int) interface{}
	//读取某行的所有数据.
	//index代表第几行默认第一行，返回的map中key是数据字段名称，value是值
	GetMap(index ...int) map[string]interface{}
	//获取结果长度
	Length() int
}

//读取查询结果
func NewQueryResult(rows *sql.Rows, fmterr FormatError) QueryResult {
	ret := &res{}
	if rows == nil {
		ret.columns = []string{}
		ret.data = nil
	} else {
		var err error
		ret.columns, err = rows.Columns()
		if err != nil {
			if fmterr != nil {
				ret.err = fmterr.FormatError(err)
			} else {
				ret.err = perr.NewError(1, "查询字段读取错误", err)
			}
		}else{
			ret.rows = rows
		}
	}
	return ret
}

//返回一个查询错误
func ErrQueryResult(err perr.Error) QueryResult {
	return &res{
		err: err,
	}
}

type res struct {
	columns    []string        //查询字段内容
	data       [][]interface{} //查询结果内容
	datalength int             //结果长度
	rows       *sql.Rows       //查询结果对象
	err        perr.Error      //查询错误
	errFmt     FormatError     //错误格式化
}

func (r *res) Error(f func(perr.Error)) QueryResult {
	if r.err != nil && f != nil {
		f(r.err)
	}
	return r
}

func (r *res) HasError() perr.Error {
	return r.err
}

func (r *res) IsEmpty() bool {
	return r.datalength < 1
}

func (r *res) Empty(f func()) QueryResult {
	if r.datalength < 1 && f != nil {
		f()
	}
	return r
}

//解析查询结果
func (r *res) passRows() {
	if r.rows != nil {
		r.data = make([][]interface{}, 0)
		for r.rows.Next() {
			row := make([]interface{}, len(r.columns))
			for i := range row {
				var ref interface{}
				row[i] = &ref
			}
			err := r.rows.Scan(row...)
			if err != nil {
				r.datalength = 0
				r.rows.Close()
				r.rows = nil
				if r.errFmt != nil {
					r.err = r.errFmt.FormatError(err)
				} else {
					r.err = perr.NewError(1, "数据读取错误", err)
				}
				return
			}
			for k, v := range row {
				row[k] = *v.(*interface{})
			}
			r.data = append(r.data, row)
		}
		r.datalength = len(r.data)
		r.rows.Close()
		r.rows = nil
	}
}

//读取某行的指定字段值.
//columnName表示字段名称，index表示第几行默认第一行，如果结果不存在返回nil
func (r *res) Get(columnName string, index ...int) interface{} {
	if len(index) < 1 {
		index = []int{0}
	}
	r.passRows()
	if index[0] >= r.datalength { //超出数据返回nil
		return nil
	}
	for i, v := range r.columns {
		if v == columnName {
			return r.data[index[0]][i]
		}
	}
	return nil
}

//读取某行的所有数据.
//index代表第几行默认第一行，返回的map中key是数据字段名称，value是值
func (r *res) GetMap(index ...int) (ret map[string]interface{}) {
	if len(index) < 1 {
		index = []int{0}
	}
	r.passRows()
	if index[0] >= r.datalength {
		return
	}
	ret = make(map[string]interface{})
	for i, v := range r.columns {
		ret[v] = r.data[index[0]][i]
	}
	return
}

//获取字段列表
func (r *res) Columns() []string {
	return r.columns
}

//获取所有数据
func (r *res) Rows() [][]interface{} {
	r.passRows()
	return r.data
}

//获取结果长度
func (r *res) Length() int {
	r.passRows()
	return r.datalength
}

//循环读取所有数据
//返回的map中key是数据字段名称，value是值,回调函数中如果返回false则停止循环后续数据
func (r *res) ForEach(f func(map[string]interface{}) bool) QueryResult {
	if f == nil {
		return r
	}
	r.passRows()
	if r.datalength < 1 { //没有数据结果直接返回
		return r
	}
	ret := map[string]interface{}{}
	for j, v := range r.data {
		if j >= r.datalength {
			break
		}
		for i, vv := range r.columns {
			ret[vv] = v[i]
		}
		if !f(ret) {
			break
		}
	}
	return r
}

//关闭查询结果
func (r *res) Close() {
	if r.rows != nil {
		r.rows.Close()
	}
}
