package tags

import (
	"reflect"
	"sync"
)

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
	vp := reflect.ValueOf(obj)
	if vp.CanInterface() {
		vp = vp.Elem()
	}
	rtype := reflect.TypeOf(vp.Interface())
	tp := field(rtype, rtype.String())
	for _, v := range tp {
		d := vp.FieldByName(v)
		switch d.Type().String() {
		case "*string":
			str, err := String(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&str))
		case "*int":
			it, err := Int(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "*int64":
			it, err := Int64(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "*float64":
			it, err := Float64(m[v])
			if err != nil {
				break
			}
			d.Set(reflect.ValueOf(&it))
		case "string":
			d.SetString(StringDefault(m[v], ""))
		case "int", "int64":
			d.SetInt(Int64Default(m[v], 0))
		case "float32", "float64":
			d.SetFloat(Float64Default(m[v], 0))
		}
	}
}

//Copy 转换对象值，对象字段名称相同并且类型相同时自动赋值相应字段
//@param obj interface{} 原数据对象
//@param obj2 interface{} 待赋值对象
func Copy(obj interface{}, obj2 interface{}) {
	vp1 := reflect.ValueOf(obj)
	if vp1.CanInterface() {
		vp1 = vp1.Elem()
	}
	vp2 := reflect.ValueOf(obj2)
	if vp2.CanInterface() {
		vp2 = vp2.Elem()
	}
	type1 := reflect.TypeOf(vp1.Interface())
	tp1 := field(type1, type1.String())
	if tp1 != nil {
		for _, k := range tp1 {
			v1 := GetPtrInterface(vp1.FieldByName(k))
			if v1 == nil {
				continue
			}
			v2 := vp2.FieldByName(k)
			if v2.IsValid() {
				switch v2.Type().String() {
				case "*string":
					str, err := String(v1)
					if str == "" || err != nil {
						break
					}
					v2.Set(reflect.ValueOf(&str))
				case "*int":
					switch v1.(type) {
					case int:
						it := v1.(int)
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case int32:
						it := int(v1.(int32))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case int64:
						it := int(v1.(int64))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					}
				case "*int64":
					switch v1.(type) {
					case int:
						it := int64(v1.(int))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case int32:
						it := int64(v1.(int32))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case int64:
						it := v1.(int64)
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					}

				case "*float32":
					switch v1.(type) {
					case float32:
						it := v1.(float32)
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case float64:
						it := float32(v1.(float64))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					}
				case "*float64":
					switch v1.(type) {
					case float32:
						it := float64(v1.(float32))
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					case float64:
						it := v1.(float64)
						if it != 0 {
							v2.Set(reflect.ValueOf(&it))
						}
					}
				case "string":
					v2.SetString(StringDefault(v1, ""))
				case "int", "int32", "int64":
					switch v1.(type) {
					case int:
						v2.SetInt(int64(v1.(int)))
					case int32:
						v2.SetInt(int64(v1.(int32)))
					case int64:
						v2.SetInt(v1.(int64))
					}
				case "float32", "float64":
					switch v1.(type) {
					case float32:
						v2.SetFloat(float64(v1.(float32)))
					case float64:
						v2.SetFloat(v1.(float64))
					}
				}
			}
		}
	}
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
