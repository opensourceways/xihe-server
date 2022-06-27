package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	//----------------------------------
	RedisNotFoundItem = "redisNotFoundItem"
)

var (
	redisURL            string
	redisMaxIdle        int = 3   //最大空闲连接数
	redisIdleTimeoutSec int = 240 //最大空闲连接时间
	redisDB             int = 0
	// RedisPrefix         string = "yc"
	redisPool *redis.Pool
)

type DBPojo interface {
	TableName() string
}

//InitRedis 初始化redis连接池
func InitRedis() error {

	redisHost := GetConfig().RedisConfig.Addr
	redisPort := GetConfig().RedisConfig.Port
	redisMaxIdle = 100
	redisDB = GetConfig().RedisConfig.Db
	redisIdleTimeoutSec = 100
	redisPassword := GetConfig().RedisConfig.Password
	redisURL = fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
	return NewRedisPool(redisURL, redisPassword, redisDB)
}

// NewRedisPool return redis pool
func NewRedisPool(redisURL, pswd string, db int) (redisErr error) {
	redisPool = &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: time.Duration(redisIdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				err = fmt.Errorf("redis connection error: %s", err)
				Log.Error(err)
				return nil, err
			}
			//check password
			if _, err = c.Do("AUTH", pswd); err != nil {
				err = fmt.Errorf("redis auth password error: %s", err)
				Log.Error(err)
				return nil, err
			}
			if _, err = c.Do("SELECT", db); err != nil {
				c.Close()
				err = fmt.Errorf("redis SELECT   error: %s", err)
				Log.Error(err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, redisErr = c.Do("PING")
			if redisErr != nil {
				return fmt.Errorf("ping redis error: %s", redisErr)
			}
			return nil
		},
	}
	return
}

//Set key:value
func Set(k string, data interface{}) error {
	c := redisPool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	_, err := c.Do("SET", k, value)
	if err != nil {
		return err
	}
	return nil
}

//BatchSet
func BatchSet(k string, data interface{}, c redis.Conn) error {

	value, _ := json.Marshal(data)
	_, err := c.Do("SET", k, value)
	if err != nil {
		return err
	}
	return nil
}

//HMSet 统一redis开关
func HMSet(k string, filed interface{}, value interface{}, c redis.Conn) error {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	value, _ = json.Marshal(value)
	_, err := c.Do("HMSET", k, filed, value)
	if err != nil {
		return err
	}
	return nil
}

//HMHINCRBY 哈希表的 某个field 自增
func HMHINCRBY(k string, filed interface{}, incre int) error {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("HINCRBY", k, filed, incre)
	if err != nil {
		return err
	}
	return nil
}

//HMGet 哈希表get
func HMGet(k string, filed interface{}, c redis.Conn) (interface{}, error) {

	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	result, err := (c.Do("HMGet", k, filed))
	if err != nil {
		return nil, err
	}
	return result, nil
}

//HVALS 哈希表get
func HVALS(k string, c redis.Conn) ([]string, error) {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	result, err := redis.Strings(c.Do("HVALS", k))
	if err != nil {
		Log.Error("HVALS:", k, err)
		return nil, err
	}
	return result, nil
}

//HGETALL 哈希表get
func HGETALL(k string) ([]interface{}, error) {

	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Values(c.Do("HGETALL", k))
	if err != nil {
		Log.Error("HGETALL:", k, err)
		return nil, err
	}
	return result, nil
}

//HDEL 哈希表get
func HDEL(k string, filed interface{}, c redis.Conn) error {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	_, err := c.Do("HDEL", k, filed)
	if err != nil {
		Log.Error("HDEL:", k, err)
		return err
	}
	return nil
}

//RedisOpen 打开redis连接
func RedisOpen() redis.Conn {
	return redisPool.Get()
}

//HMGetPojo 在不关闭redis的情况下 读取数据
func HMGetPojo(k string, filed interface{}, pojo DBPojo, c redis.Conn) error {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	dbpojoJSON, err := redis.Strings(c.Do("HMGet", k, filed))
	if err != nil {
		return fmt.Errorf("HMGetPojo  k:%v, filed:%v. err:%v", k, filed, err)
	}

	if len(dbpojoJSON) == 0 {
		return fmt.Errorf("HMGetPojo  k:%v, filed:%v", k, filed)
	}
	if len(dbpojoJSON[0]) == 0 {
		Log.Info(fmt.Errorf("%s,   key:%s, field:%v ", RedisNotFoundItem, k, filed))
		return nil
	}
	err = json.Unmarshal([]byte(dbpojoJSON[0]), pojo)
	if err != nil {
		return fmt.Errorf("HMGetPojo Unmarshal dbpojoJSON:%v   --- error:%v", dbpojoJSON[0], err)
	}
	return nil
}

//RedisClose 关闭redis连接
func RedisClose(c redis.Conn) {
	c.Close()
}

//Lpush 插入多个
func Lpush(k string, data interface{}) error {
	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("LPUSH", k, data)
	if err != nil {
		return err
	}
	return nil
}

//LBatchpush 插入多个
func LBatchpush(k string, data interface{}, c redis.Conn) error {
	_, err := c.Do("LPUSH", k, data)
	if err != nil {
		return err
	}
	return nil
}

//PlusOne 值累加1
func PlusBatchOne(k string, c redis.Conn) int {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	result, err := redis.Int(c.Do("INCR", k))
	if err != nil {
		return -1
	}
	return result
}

//GetInt 获取数字
func GetInt(k string) (int, error) {

	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Int(c.Do("GET", k))
	if err != nil {
		return 0, err
	}
	return result, nil
}

//GetFloat 获取数字
func GetFloat(k string) (float64, error) {

	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Float64(c.Do("GET", k))
	if err != nil {
		return 0, err
	}
	return result, nil
}

//Rpush 插入多个
func Rpush(k string, data int) error {

	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("RPUSH", k, data)
	if err != nil {
		return err
	}
	return nil
}

//RpushString 插入多个Object
func RpushString(k string, data string) (count int, err error) {

	c := redisPool.Get()
	defer c.Close()
	count, err = redis.Int(c.Do("RPUSH", k, data))
	return
}

//LpushString 插入多个Object
func LpushString(k string, data string) (count int, err error) {

	c := redisPool.Get()
	defer c.Close()
	count, err = redis.Int(c.Do("LPUSH", k, data))
	return
}

//LREM 移除
func LREM(k string, count int, value interface{}) (err error) {

	c := redisPool.Get()
	defer c.Close()
	_, err = c.Do("LREM", k, count, value)
	return
}

//LTRIM 只保留左边的多少个
func LTRIM(k string, start, stop int) (err error) {

	c := redisPool.Get()
	defer c.Close()
	_, err = c.Do("LTRIM", k, start, stop)
	return
}

//Lrange 获取
func Lrange(k string, start, end int) ([]int, error) {
	c := redisPool.Get()
	defer c.Close()

	idlist, err := redis.Ints(c.Do("LRANGE", k, start, end))
	if err != nil {
		Log.Error("Get Error: ", err.Error(), k)
		return nil, err
	}
	return idlist, nil
}

//LrangeString 从左边获取多个
func LrangeString(k string, start, end int) ([]string, error) {
	c := redisPool.Get()
	defer c.Close()
	idlist, err := redis.Strings(c.Do("LRANGE", k, start, end))
	if err != nil {
		Log.Error("Get Error: ", err.Error(), k)
		return nil, err
	}
	return idlist, nil
}

//LSET 插入多个
func LSET(k string, index int, value string) error {

	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("LSET", k, index, value)
	if err != nil {
		Log.Error("Get Error: ", err.Error(), k)
		return err
	}
	return nil
}

//Rrange 从右边获取最新的数据
func Rrange(k string, start, end int) ([]int, error) {

	c := redisPool.Get()
	defer c.Close()

	idlist, err := redis.Ints(c.Do("RRANGE", k, start, end))
	if err != nil {
		Log.Error("Get Error: ", err.Error(), k)
		return nil, err
	}
	return idlist, nil
}

//GetKeyTTL 获取这个key的剩余过期时间，秒
func GetKeyTTL(k string) int {

	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Int(c.Do("TTL", k))
	if err != nil {
		Log.Error("GetKeyTTL:", err)
		return 0
	}
	return result
}

func SetKeyExpire(k string, ex int, c redis.Conn) {

	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		Log.Error("set error", err.Error())
	}
}
func SetKeyWithExpire(k string, data interface{}, ex int, c redis.Conn) error {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", k, value, "ex", ex)

	return err

}

func CheckKey(k string, c redis.Conn) bool {

	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		Log.Error(err)
		return false
	}
	return exist

}

func DelKey(k string, c redis.Conn) error {

	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	_, err := c.Do("DEL", k)
	if err != nil {
		Log.Error(err)
		return err
	}
	return nil
}

//GetJsonByte GetJsonByte
func GetJsonByte(k string) ([]byte, error) {
	c := redisPool.Get()
	defer c.Close()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return jsonGet, nil
}

//RedisGet GetJsonByte
func RedisGet(k string, c redis.Conn) ([]byte, error) {
	if c == nil {
		c = redisPool.Get()
		defer c.Close()
	}
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return jsonGet, nil
}

//SADD 用来向set里面增加
func SADD(k string, fileID int) error {

	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("SADD", k, fileID)
	if err != nil {
		return err
	}
	return nil
}

//SADDString 用来向set里面增加string
func SADDString(k string, value string) error {

	c := redisPool.Get()
	defer c.Close()
	_, err := c.Do("SADD", k, value)
	if err != nil {
		return err
	}
	return nil
}

//SMEMBERS 获取k的所有memebers
func SMEMBERS(k string) (interface{}, error) {

	c := redisPool.Get()
	defer c.Close()
	result, err := c.Do("SMEMBERS", k)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//ZADD 用来向Zset里面增加
func ZADD(key string, score int, value interface{}, c redis.Conn) error {

	value, _ = json.Marshal(value)
	_, err := c.Do("ZADD", key, score, value)
	if err != nil {
		return err
	}
	return nil
}

//CacheSetDataPlus  缓存一些数据，每隔一段时间再更新数据库，避免db较大压力
func CacheSetDataPlus(dataname string, dataField string) int {
	k := fmt.Sprintf("cache:%s", dataname)
	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Int(c.Do("HGET", k, dataField))
	if err != nil {
		result = 0
	}
	_, err = c.Do("HSET", k, dataField, result+1)
	if err != nil {
		Log.Error("CacheSetDataPlus HSET " + dataname + " field:" + dataField)
		return 0
	}
	return result + 1
}

//DianZhan 点赞，点过的就取消点赞
func DianZhan(tablename, channelStr, docid string, userid int, c redis.Conn) int {

	key := fmt.Sprintf("%s:%s:%s", tablename, channelStr, docid)
	result, err := redis.Int(c.Do("GETBIT", key, userid))
	if err != nil {
		Log.Error("DianZhan： ", err)
		return -1
	}
	if result == 0 {
		result = 1
	} else {
		result = 0
	}
	_, err = c.Do("SETBIT", key, userid, result)
	Log.Error(key, userid, result, err)

	return result
}

//GetDianZhan 查看某个作品是否已经点赞
func GetDianZhan(tablename, channelStr, docid string, userid int, c redis.Conn) int {

	key := fmt.Sprintf("%s:%s:%s", tablename, channelStr, docid)
	result, err := redis.Int(c.Do("GETBIT", key, userid))
	if err != nil {
		Log.Error("DianZhan： ", err)
		return -1
	}
	return result
}

//CacheGetAllKey 获取一个cache中的所有field
func CacheGetAllKey(dataname string) []string {

	k := fmt.Sprintf("cache:%s", dataname)
	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Strings(c.Do("HKEYS", k))
	if err != nil {
		// Log.Error("CacheGetData HGET " + dataname + " field:" + dataField)
		return nil
	}
	return result
}

//GetKeyFromRedisByLabel 根据名称获取它在redis里面的Key标签
func GetKeyFromRedisByLabel(tabelname string, label string) string {
	labelKey := fmt.Sprintf("%s_label:%s", tabelname, label)
	return labelKey
}
