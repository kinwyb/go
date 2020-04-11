package mysql

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/kinwyb/go/db/tags"
)

//SetSQL 转换成插入语句
func SetSQL(obj interface{}) (string, []interface{}) {
	vp := reflect.ValueOf(obj)
	if vp.Kind() == reflect.Interface || vp.Kind() == reflect.Ptr {
		vp = vp.Elem()
	}
	retinterface := make([]interface{}, 0)
	buf := &bytes.Buffer{}
	rtype := reflect.TypeOf(vp.Interface())
	tp, _ := tags.DbTag(rtype)
	var value, srcValue reflect.Value
	if tp != nil {
		for k, v := range tp {
			srcValue = vp
			for {
				ks := strings.SplitN(k, ":", 2)
				if len(ks) > 1 {
					srcValue = srcValue.FieldByName(ks[0])
					k = ks[1]
				} else {
					value = srcValue.FieldByName(k)
					break
				}
			}
			if (value.Kind() == reflect.Ptr && value.IsNil()) || tags.IsEmpty(value) {
				continue
			}
			buf.WriteString("`")
			buf.WriteString(v)
			buf.WriteString("` = ?,")
			retinterface = append(retinterface, tags.GetPtrInterface(value))
		}
	}
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String(), retinterface
}

//Update 更新数据
func Update(table string, obj interface{}) (string, []interface{}, error) {
	vp := reflect.ValueOf(obj)
	if v, ok := obj.(reflect.Value); ok {
		vp = v
	} else {
		vp = reflect.ValueOf(obj)
		if vp.CanInterface() {
			vp = vp.Elem()
		}
	}
	retinterface := make([]interface{}, 0)
	rtype := reflect.TypeOf(vp.Interface())
	tp, primary := tags.DbTag(rtype)
	if primary == "" {
		return "", nil, errors.New("更新数据无主键")
	}
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(table)
	buf.WriteString(" SET ")
	var paramkey interface{}
	var value, srcValue reflect.Value
	if tp != nil {
		for k, v := range tp {
			srcValue = vp
			for {
				ks := strings.SplitN(k, ":", 2)
				if len(ks) > 1 {
					srcValue = srcValue.FieldByName(ks[0])
					k = ks[1]
				} else {
					value = srcValue.FieldByName(k)
					break
				}
			}
			if value.Kind() == reflect.Ptr && value.IsNil() {
				continue
			} else if primary == v {
				paramkey = tags.GetPtrInterface(value)
				continue
			}
			vl := tags.GetPtrInterface(value)
			if _, ok := vl.(string); ok && vl.(string) == "" {
				buf.WriteString(v + " = null,")
			} else {
				buf.WriteString(v)
				buf.WriteString(" = ?,")
				retinterface = append(retinterface, vl)
			}
		}

	}
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteString(" WHERE ")
	buf.WriteString(primary)
	buf.WriteString(" = ?")
	retinterface = append(retinterface, paramkey)
	return buf.String(), retinterface, nil
}

//SelectSQL 获取对象
func SelectSQL(obj interface{}, tablename ...string) *bytes.Buffer {
	vp := reflect.ValueOf(obj)
	if vp.CanInterface() {
		vp = vp.Elem()
	}
	buf := &bytes.Buffer{}
	if tablename != nil && len(tablename) > 0 {
		buf.WriteString("SELECT ")
	}
	rtype := reflect.TypeOf(vp.Interface())
	tp, _ := tags.DbTag(rtype)
	if tp != nil {
		for k, v := range tp {
			buf.WriteString(v)
			buf.WriteString(" `")
			if strings.Contains(k, ":") {
				ks := strings.Split(k, ":")
				l := len(ks) - 1
				if len(ks) > 1 && ks[1] != "" {
					k = ks[l]
				}
			}
			buf.WriteString(k)
			buf.WriteString("`,")
		}
	}
	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}
	if tablename != nil && len(tablename) > 0 {
		buf.WriteString(" FROM ")
		buf.WriteString(tablename[0])
	}
	return buf
}

//拼接 where in 条件,field要查询的字段名称以及拼接条件，value表示in中的字段值多个值按char[默认逗号]分割,args表示参数集合
//返回拼接结果和参数集合，eg: whereIN("AND material_code", "xxx,xxxxx", args, bs)
func WhereIN(field string, values string, args []interface{}, bs *bytes.Buffer, char ...string) (*bytes.Buffer, []interface{}) {
	if bs == nil {
		bs = bytes.NewBufferString("")
	}
	if args == nil {
		args = make([]interface{}, 0)
	}
	if values == "" {
		return bs, args
	} else if len(char) < 1 {
		char = []string{","}
	}
	bs.WriteString(field)
	bs.WriteString(" IN (")
	value := strings.Split(values, char[0])
	for _, v := range value {
		bs.WriteString("?,")
		args = append(args, v)
	}
	bs.Truncate(bs.Len() - 1)
	bs.WriteString(")")
	return bs, args
}
