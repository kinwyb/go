package tags

import (
	"bytes"
	"encoding/gob"
	"github.com/kinwyb/go/conv"
	"reflect"
	"sync"
)

//byteBuffer池
var byteBufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}
var dbTags map[string]map[string]string
var primaryTags map[string]string
var fields map[string][]string
var dbTagsLock *sync.RWMutex

func init() {
	dbTags = map[string]map[string]string{}
	primaryTags = map[string]string{}
	fields = map[string][]string{}
	dbTagsLock = &sync.RWMutex{}
}

//fieldObj 解析对象结构
//@param tp reflect.Type 对象类型
//@param vname string 对象名称
func fieldObj(tp reflect.Type, vname string) {
	x := tp.NumField()
	var field reflect.StructField
	var db, primary string
	tag := map[string]string{}
	fds := make([]string, x)
	for i := 0; i < x; i++ {
		field = tp.Field(i)
		fds[i] = field.Name
		db = field.Tag.Get("db")
		if db == "" {
			if field.Type.Kind() == reflect.Struct {
				//如果是结构体嵌套的,读取结构体当中的内容
				tp, _ := DbTag(field.Type)
				for k, v := range tp {
					tag[field.Name+":"+k] = v
				}
			}
			continue
		}
		tag[field.Name] = db
		if field.Tag.Get("primary") == "true" {
			primary = db
		}
	}
	dbTagsLock.Lock()
	dbTags[vname] = tag
	primaryTags[vname] = primary
	fields[vname] = fds
	dbTagsLock.Unlock()
}

//dbTag 解析对象数据库结构
//@param tp reflect.Type 对象类型
//@return map[string]string 对象数据库结构字段对照表
//@return string 对象数据库主键字段名
func DbTag(tp reflect.Type) (map[string]string, string) {
	dbTagsLock.RLock()
	if v, ok := dbTags[tp.String()]; ok {
		dbTagsLock.RUnlock()
		return v, primaryTags[tp.String()]
	}
	dbTagsLock.RUnlock()
	fieldObj(tp, tp.String())
	return dbTags[tp.String()], primaryTags[tp.String()]
}

//field 获取对象字段信息
//@param tp reflect.Type 对象类型
//@param vname string 对象名称
//@return []string 字段列表
func field(tp reflect.Type, vname string) []string {
	dbTagsLock.RLock()
	if v, ok := fields[vname]; ok {
		dbTagsLock.RUnlock()
		return v
	}
	dbTagsLock.RUnlock()
	fieldObj(tp, vname)
	return fields[vname]
}

//SetMapValue 设置对象值
func SetMapValue(obj interface{}, m map[string]interface{}) {
	var vp reflect.Value
	if v, ok := obj.(reflect.Value); ok {
		vp = v
	} else {
		vp = reflect.ValueOf(obj)
		if vp.CanInterface() {
			vp = vp.Elem()
		}
	}
	rtype := reflect.TypeOf(vp.Interface())
	tp := field(rtype, rtype.String())
	for _, v := range tp {
		d := vp.FieldByName(v)
		if d.Type().Kind() == reflect.Struct { //如果是结构体则再次解析里面内容
			SetMapValue(d, m)
			continue
		}
		switch d.Type().String() {
		case "*string":
			str, err := conv.ToStringE(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&str))
		case "*int":
			it, err := conv.ToIntE(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "*int64":
			it, err := conv.ToInt64E(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "*float64":
			it, err := conv.ToFloat64E(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "string":
			d.SetString(conv.ToString(m[v]))
		case "int", "int64":
			d.SetInt(conv.ToInt64(m[v]))
		case "float32", "float64":
			d.SetFloat(conv.ToFloat64(m[v]))
		}
	}
}

//Copy 转换对象值，对象字段名称相同并且类型相同时自动赋值相应字段
//@param dst interface{} 目标对象
//@param src interface{} 原始对象
func Copy(dst, src interface{}) {
	var buf = byteBufferPool.Get().(*bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		buf.Truncate(0)
		byteBufferPool.Put(buf)
		return
	}
	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
	buf.Truncate(0)
	byteBufferPool.Put(buf)
}

//getPtrInterface 获取指针的原始对象
func GetPtrInterface(v reflect.Value) interface{} {
	for {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil
			}
			v = v.Elem()
		} else {
			break
		}
	}
	return v.Interface()
}

//判断值是否为空
func IsEmpty(v reflect.Value) bool {
	return GetPtrInterface(v) == reflect.Zero(v.Type()).Interface()
}
