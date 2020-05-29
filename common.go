package heldiamgo

import (
	"encoding/json"
	"regexp"
	"strings"

	"math/rand"

	"time"

	"fmt"

	"unicode"

	"strconv"

	"github.com/coreos/etcd/pkg/idutil"
	"github.com/pborman/uuid"
)

//时间格式
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
)

var idGen = idutil.NewGenerator(uint16(rand.Uint32()>>16), time.Now())

func UUID() string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1)
}

//生成数字ID带日期
func IDGen() string {
	next := idGen.Next()
	return fmt.Sprintf("%s%d", time.Now().Format("060102"), next)
}

//生成唯一数字ID
func ID() uint64 {
	return idGen.Next()
}

//生成随机字符串,len字符串长度,onlynumber 是否只包含数字
func RandCode(len int, onlynumber bool) string {
	if len < 1 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	ret := make([]byte, len)
	if onlynumber {
		chars := []byte("0123456789")
		for i := 0; i < len; i++ {
			ret[i] = chars[rand.Int31n(9)]
		}
	} else {
		chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
		for i := 0; i < len; i++ {
			ret[i] = chars[rand.Int31n(35)]
		}
	}
	return string(ret[:])
}

//验证是否是有效的邮箱
func IsEmail(emailaddress string) bool {
	regexpstring := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}`
	ret, _ := regexp.MatchString(regexpstring, emailaddress)
	return ret
}

//验证是否是有效的手机号
func IsPhone(phonenumber string) bool {
	//regexpstring := `^1[3|4|5|7|8|9][0-9]\d{8}$`
	regexpstring := `^1\d{10}$`
	ret, _ := regexp.MatchString(regexpstring, phonenumber)
	return ret
}

//验证是否有效的身份证号码
func IsIdCard(Identnumber string) bool {
	regexpstring := `^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$|^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}([0-9]|X)$`
	ret, _ := regexp.MatchString(regexpstring, Identnumber)
	return ret
}

// 验证是否为家庭电话
func IsHomePhone(phone string) bool {
	regexpstring := `^(\d{3}-)(\d{8})$|^(\d{4}-)(\d{7})$|^(\d{4}-)(\d{8})$`
	ret, _ := regexp.MatchString(regexpstring, phone)
	return ret
}

//字符串数组去重
func RemoveStringArrayDuplicate(list []string) []string {
	var x []string
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//判断是否有中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) ||
			(regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

// 当前日期时间
func TimeNow() string {
	return time.Now().Format(DateTimeFormat)
}

//json数据
func JsonString(obj interface{}) string {
	ret, _ := json.Marshal(obj)
	return string(ret)
}

//验证是否是时间字符串
func IsTimeString(timeString string, timeFormat string) bool {
	_, err := time.ParseInLocation(timeFormat, timeString, time.Local)
	return err == nil
}

// 判断字符串类型的切片元素是否重复
func StringSliceIsRepeat(slc []string) bool {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	if len(slc) == len(result) {
		return false
	}
	return true
}

// 浮点数保留指定的位数
// @value  float64 要转换的浮点数值
// @preciseDigits int 要保留的位数
func Decimal(value float64, preciseDigits int) float64 {
	precise := "%." + strconv.Itoa(preciseDigits) + "f"
	value, _ = strconv.ParseFloat(fmt.Sprintf(precise, value), 64)
	return value
}

//日期起始时间
func DateStartTimeString(t time.Time) string {
	return t.Format(DateFormat) + " 00:00:00"
}

//日期结束时间
func DateEndTimeString(t time.Time) string {
	return t.Format(DateFormat) + " 23:59:59"
}

// int64数组去除重复
func RemoveInt64ArrayDuplicate(arr []int64) (newArr []int64) {
	newArr = make([]int64, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// int32数组去除重复
func RemoveInt32ArrayDuplicate(arr []int32) (newArr []int32) {
	newArr = make([]int32, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// int数组去除重复
func RemoveIntArrayDuplicate(arr []int) (newArr []int) {
	newArr = make([]int, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// float64数组去除重复
func RemoveFloat64ArrayDuplicate(arr []float64) (newArr []float64) {
	newArr = make([]float64, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// float32数组去除重复
func RemoveFloat32ArrayDuplicate(arr []float32) (newArr []float32) {
	newArr = make([]float32, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
