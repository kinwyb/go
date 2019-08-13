package heldiamgo

import (
	"time"

	"github.com/gomodule/redigo/redis"

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
	rclient.Close()
	if err != nil && r.log != nil {
		r.log.Error("[SET]失败:%s", err.Error())
	}
	return err == nil
}

//SETEX 设置值和过期时间
func (r *RedisUtil) SETEX(key, value string, expireTime int64) bool {
	if r.debug {
		return true
	}
	rclient := r.pool.Get()
	_, err := rclient.Do("SETEX", r.prefix+key, expireTime, value)
	rclient.Close()
	if err != nil && r.log != nil {
		r.log.Error("[SETEX]失败:%s", err.Error())
	}
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
	rclient.Close()
	if err != nil {
		if r.log != nil {
			r.log.Error("[EXPIRE]失败:%s", err.Error())
		}
		return -1
	}
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
	rclient.Close()
	if err != nil {
		if r.log != nil {
			r.log.Error("[DEL]失败:%s", err.Error())
		}
		return -1
	}
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
	rclient.Close()
	if err != nil {
		if r.log != nil {
			r.log.Error("[KEYS]失败:%s", err.Error())
		}
		return nil, err
	}
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
	rclient.Close()
	if err != nil || ret == nil {
		if r.log != nil {
			r.log.Error("[GET]失败:%s", err.Error())
		}
		return ""
	}
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
	rclient.Close()
	if err != nil || ret == nil {
		if r.log != nil {
			r.log.Error("[GETSET]失败:%s", err.Error())
		}
		return ""
	}
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

// 获取集合所有值
func (r *RedisUtil) SMEMBERS(key string) ([]string, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("SMEMBERS", r.prefix+key)
	rclient.Close()
	if err != nil {
		return nil, err
	}
	return db.Strings(ret)
}

//原子操作
func (r *RedisUtil) INCRBY(key string, increment int64) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("INCRBY", r.prefix+key, increment)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 将元素添加到集合,返回0/1(失败/成功)
func (r *RedisUtil) SADD(key, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("SADD", r.prefix+key, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 将元素添加有序集合,返回0/1(失败/成功)
func (r *RedisUtil) ZADD(key string, score float64, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZADD", r.prefix+key, score, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 返回有序集合成员score值
func (r *RedisUtil) ZSCORE(key string, member string) (float64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZSCORE", r.prefix+key, member)
	rclient.Close()
	return db.Float64Default(ret), err
}

// 返回有序集合key中，指定区间的成员(按score从大到小排序)
func (r *RedisUtil) ZREVRANGE(key string, start, stop int64, withscores bool) ([]string, error) {
	args := []interface{}{
		r.prefix + key,
		start,
		stop,
	}
	if withscores {
		args = append(args, "WITHSCORES")
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZREVRANGE", args...)
	rclient.Close()
	if err != nil {
		return nil, err
	}
	return db.Strings(ret)
}

// 返回有序集合key中，指定区间的成员(按score从小到大排序)
func (r *RedisUtil) ZRANGE(key string, start, stop int64, withscores bool) ([]string, error) {
	args := []interface{}{
		r.prefix + key,
		start,
		stop,
	}
	if withscores {
		args = append(args, "WITHSCORES")
	}
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZRANGE", args...)
	rclient.Close()
	if err != nil {
		return nil, err
	}
	return db.Strings(ret)
}

// 返回有序集合key中,指定成员的排名(按score从小到大排序),第一名为0
func (r *RedisUtil) ZRANK(key string, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZRANK", r.prefix+key, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 返回有序集合key中,指定成员的排名(按score从大到小排序),第一名为0
func (r *RedisUtil) ZREVRANK(key string, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZREVRANK", r.prefix+key, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 移除有序集合key中指定成员的排名,返回移除数量
func (r *RedisUtil) ZREM(key string, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZREM", r.prefix+key, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 移除集合key中指定成员的排名,返回移除数量
func (r *RedisUtil) SREM(key string, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("SREM", r.prefix+key, member)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 返回集合key成员数量
func (r *RedisUtil) SCARD(key string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("SCARD", r.prefix+key)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 返回有序集合key成员数量
func (r *RedisUtil) ZCARD(key string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZCARD", r.prefix+key)
	rclient.Close()
	return db.Int64Default(ret), err
}

// 为有序集合key成员score增加值
func (r *RedisUtil) ZINCRBY(key string, increment int64, member string) (int64, error) {
	rclient := r.pool.Get()
	ret, err := rclient.Do("ZINCRBY", r.prefix+key, increment, member)
	rclient.Close()
	return db.Int64Default(ret), err
}
