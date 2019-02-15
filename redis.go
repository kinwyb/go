package heldiamgo

import (
	"time"

	"github.com/gomodule/redigo/redis"

	"encoding/json"

	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/logs"
)

//RedisUtil redis操作工具
type RedisUtil struct {
	pool   *redis.Pool //redis连接池
	prefix string      //前缀
	debug  bool
	log    logs.Logger
}

//初始化redis,host=地址，password密码,db数据名称，maxidle最大有效连接数,active 最大可用连接数
func InitializeRedis(host, password, db string, maxidle, maxactive int) *RedisUtil {
	if maxidle < 10 {
		maxidle = 10
	}
	if maxactive < 10 {
		maxactive = 10
	}
	return &RedisUtil{
		pool: &redis.Pool{ // 建立连接池
			MaxIdle:     maxidle,
			MaxActive:   maxactive,
			IdleTimeout: 180 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host)
				if err != nil {
					return nil, err
				}
				if password != "" {
					//认证
					c.Do("AUTH", password)
				}
				// 选择db
				c.Do("SELECT", db)
				return c, nil
			},
		},
		debug: false,
	}
}

//设置是否调试，调试状态不进入真正缓存
func (r *RedisUtil) Debug(b bool) {
	r.debug = b
}

//设置缓存前缀
func (r *RedisUtil) Prefix(prefix string) {
	r.prefix = prefix
}

//设置日志
func (r *RedisUtil) SetLogger(lg logs.Logger) {
	r.log = lg
}

//SET 设置值
func (r *RedisUtil) SET(key, value string) bool {
	if r.debug {
		return true
	}
	rclient := r.pool.Get()
	_, err := rclient.Do("SET", r.prefix+key, value)
	if err != nil && r.log != nil {
		r.log.Error("[SET]失败:%s", err.Error())
	}
	rclient.Close()
	return err == nil
}

//SETEX 设置值和过期时间
func (r *RedisUtil) SETEX(key, value string, expireTime int64) bool {
	if r.debug {
		return true
	}
	rclient := r.pool.Get()
	_, err := rclient.Do("SETEX", r.prefix+key, expireTime, value)
	if err != nil && r.log != nil {
		r.log.Error("[SETEX]失败:%s", err.Error())
	}
	rclient.Close()
	return err == nil
}

//EXPIRE 设置过期时间
//@return -1=错误，0=key不存在， 1=成功
func (r *RedisUtil) EXPIRE(key string, expireTime int64) int {
	if r.debug {
		return -1
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("EXPIRE", r.prefix+key, expireTime)
	if err != nil {
		if r.log != nil {
			r.log.Error("[EXPIRE]失败:%s", err.Error())
		}
		rclient.Close()
		return -1
	}
	rclient.Close()
	return db.IntDefault(ret, -1)
}

//DEL 删除
//@return 返回删除的个数,-1代表错误
func (r *RedisUtil) DEL(key ...interface{}) int {
	if r.debug {
		return -1
	}
	keys := make([]interface{}, len(key))
	for i, v := range key {
		switch v.(type) {
		case string:
			keys[i] = r.prefix + v.(string)
		default:
			keys[i] = v
		}
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("DEL", keys...)
	if err != nil {
		if r.log != nil {
			r.log.Error("[DEL]失败:%s", err.Error())
		}
		rclient.Close()
		return -1
	}
	rclient.Close()
	return db.IntDefault(ret, -1)
}

//KEYS 获取key集合
//@param {string} pattern 匹配模式
// Supported glob-style patterns:
// h?llo matches hello, hallo and hxllo
// h*llo matches hllo and heeeello
// h[ae]llo matches hello and hallo, but not hillo
// h[^e]llo matches hallo, hbllo, ... but not hello
// h[a-b]llo matches hallo and hbllo
// Use \ to escape special characters if you want to match them verbatim.
//@return 返回符合patten的key
func (r *RedisUtil) KEYS(pattern string) ([]string, error) {
	if r.debug {
		return []string{}, nil
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("KEYS", r.prefix+pattern)
	if err != nil {
		if r.log != nil {
			r.log.Error("[KEYS]失败:%s", err.Error())
		}
		rclient.Close()
		return nil, err
	}
	rclient.Close()
	return db.Strings(ret)
}

//GET 获取
//@return 如果key不存在返回空
func (r *RedisUtil) GET(key string) string {
	if r.debug {
		return ""
	}
	rclient := r.pool.Get()
	key = r.prefix + key
	ret, err := rclient.Do("GET", key)
	if err != nil || ret == nil {
		if r.log != nil {
			r.log.Error("[GET]失败:%s", err.Error())
		}
		rclient.Close()
		return ""
	}
	rclient.Close()
	return db.StringDefault(ret, "")
}

//GETSET 获取旧值并设置新值
//@return 如果key不存在返回空
func (r *RedisUtil) GETSET(key, value string) string {
	if r.debug {
		return ""
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("GETSET", r.prefix+key, value)
	if err != nil || ret == nil {
		if r.log != nil {
			r.log.Error("[GETSET]失败:%s", err.Error())
		}
		rclient.Close()
		return ""
	}
	rclient.Close()
	return db.StringDefault(ret, "")
}

//GETEXP 获取值并更新过期时间
//@return 如果key不存在返回空
func (r *RedisUtil) GETEXP(key string, expireTime int64) string {
	if r.debug {
		return ""
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("GET", r.prefix+key)
	if err != nil || ret == nil {
		if r.log != nil {
			r.log.Error("[GETEXP]:%s", err.Error())
		}
		rclient.Close()
		return ""
	}
	rclient.Do("EXPIRE", r.prefix+key, expireTime)
	rclient.Close()
	return db.StringDefault(ret, "")
}

//GetConn 获取redis链接
func (r *RedisUtil) GetConn(fun func(redis.Conn)) {
	rclient := r.pool.Get()
	fun(rclient)
	rclient.Close()
}

//GETOBJ 获取缓存
func (r *RedisUtil) GETOBJ(key string, obj interface{}) bool {
	if obj == nil || r.debug {
		return false
	}
	str := r.GET(r.prefix + key)
	if str != "" {
		if v, ok := obj.(json.Unmarshaler); ok {
			return v.UnmarshalJSON([]byte(str)) == nil
		} else {
			return json.Unmarshal([]byte(str), obj) == nil
		}
	}
	return false
}

//GETOBJEXP 获取缓存并更新缓存时间
func (r *RedisUtil) GETOBJEXP(key string, obj interface{}, expireTime int64) bool {
	if obj == nil || r.debug {
		return false
	}
	str := r.GETEXP(r.prefix+key, expireTime)
	if str != "" {
		if v, ok := obj.(json.Unmarshaler); ok {
			return v.UnmarshalJSON([]byte(str)) == nil
		} else {
			return json.Unmarshal([]byte(str), obj) == nil
		}
	}
	return false
}

//PUT 设置缓存
func (r *RedisUtil) PUT(key string, obj interface{}) {
	if obj == nil || r.debug {
		return
	}
	var b []byte
	var err error
	if v, ok := obj.(json.Marshaler); ok {
		b, err = v.MarshalJSON()
	} else {
		b, err = json.Marshal(obj)
	}
	if err != nil {
		if r.log != nil {
			r.log.Error("[PUT]:%s", err.Error())
		}
		return
	}
	r.SET(r.prefix+key, string(b))
}

//PUTEXP 设置缓存
func (r *RedisUtil) PUTEXP(key string, obj interface{}, time int64) {
	if obj == nil || r.debug {
		return
	}
	var b []byte
	var err error
	if v, ok := obj.(json.Marshaler); ok {
		b, err = v.MarshalJSON()
	} else {
		b, err = json.Marshal(obj)
	}
	if err != nil {
		if r.log != nil {
			r.log.Error("[PUTEXP]:%s", err.Error())
		}
		return
	}
	r.SETEX(r.prefix+key, string(b), time)
}
