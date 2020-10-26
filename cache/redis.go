package cache

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

// TODO: 详细配置
type RedisConfig struct {
	Network  string
	Address  string
	Password string
}

var (
	once   sync.Once
	pool   *redis.Pool
	logger = logging.GetLogger("cache_redis")
)

func InitRedisPool(config RedisConfig) {
	once.Do(func() {
		pool = &redis.Pool{
			MaxIdle:     50,
			MaxActive:   30,
			IdleTimeout: 300 * time.Second,
			Dial: func() (redis.Conn, error) {
				// 1. 打开连接
				c, err := redis.Dial(config.Network, config.Address)
				if err != nil {
					logger.Panic(err)
					return nil, err
				}

				// 2. 访问认证
				if config.Password != "" {
					if _, err = c.Do("AUTH", config.Password); err != nil {
						logger.Panic(err)
						_ = c.Close()
						return nil, err
					}
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
	})
}

func GetRedisPool() *redis.Pool {
	if pool == nil {
		logger.Panic("GetRedisPool before InitRedisPool")
		return nil
	}
	return pool
}
