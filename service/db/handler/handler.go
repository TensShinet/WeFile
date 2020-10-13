package handler

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/db/conn"
	idg "github.com/TensShinet/WeFile/service/id_generator/proto"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"gorm.io/gorm"
)

type Service struct{}

var (
	db                *gorm.DB
	redisCli          *redis.Client
	generateIDService idg.GenerateIDService
	logger            = logging.GetLogger("db_service_handler")
)

// 调用 handler Init 之前一定要初始化 conn
func Init() {
	redisCli = conn.GetRedisCli()
	db = conn.GetDB()
	if redisCli == nil || db == nil {
		logger.Panic("db or redis cli is not initialized")
	}
	config := conf.GetConfig()
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Registry(micReg),
	)
	generateIDService = idg.NewGenerateIDService("go.micro.service.id_generator", service.Client())
}
