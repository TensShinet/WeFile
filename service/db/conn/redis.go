package conn

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/go-redis/redis"
	"sync"
)

var (
	redisOnce sync.Once
	redisCli  *redis.Client
)

func redisInit() {
	redisOnce.Do(func() {
		config := conf.GetConfig()
		redisInfo := config.Redis
		// TODO: 支持哨兵模式
		redisCli = redis.NewClient(&redis.Options{
			Addr:     redisInfo.Conn,
			Password: redisInfo.Password,
			DB:       redisInfo.DBNum,
		})
		// 检测连接
		if _, err := redisCli.Ping().Result(); err != nil {
			logger.Panicf("connecting redis failed, for the reason:%v", err.Error())
		}

	})
}

func GetRedisCli() *redis.Client {
	return redisCli
}
