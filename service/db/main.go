package main

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/db/conf"
	"github.com/TensShinet/WeFile/service/db/conn"
	"github.com/TensShinet/WeFile/service/db/handler"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"time"
)

var logger = logging.GetLogger("db_service")

// db service 只针对各个模块对数据库或者内存数据库的操作
// 基本没有逻辑代码
func main() {
	conf.Init("db_conf.yml")

	config := conf.GetConfig()
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Service.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Name("go.micro.service.db"),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*time.Duration(config.Service.RegisterTTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Service.RegisterInterval)),
	)
	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// 初始化 db 连接
			conn.Init()
			model.Init()
			handler.Init()
			return nil
		}),
	)

	// 注册服务
	if err := proto.RegisterServiceHandler(service.Server(), new(handler.Service)); err != nil {
		logger.Panicf("register failed, for the reason:%v", err.Error())
	}
	// 启动服务
	if err := service.Run(); err != nil {
		logger.Panicf("start failed, for the reason:%v", err.Error())
	}

	logger.Info("running db service")
}
