package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "192.168.48.129:6379"
	redisPass = ""
)

//newRedisPool:创建redis连接池

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1.打开链接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			//2.访问验证：在这个地方其实可以不用，DO相当于使用一个命令，当不需要密码的时候不需要使用这个-
			if _, err = c.Do("", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
