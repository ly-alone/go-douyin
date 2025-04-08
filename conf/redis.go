package conf

import (
	"github.com/gomodule/redigo/redis"
)

var (
	redisDb redis.Conn
)

func Redis() redis.Conn {
	return redisDb
}
func RedisDb() {
	redisDb, err := redis.Dial("tcp", ":63799")
	if err != nil {
		defer redisDb.Close()
		panic(err)
	}

}
