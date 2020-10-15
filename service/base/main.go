package main

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	baseConf "github.com/TensShinet/WeFile/service/base/conf"
	"github.com/TensShinet/WeFile/service/base/handler"
	"github.com/TensShinet/WeFile/service/base/router"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry/etcd"
	"time"
)

var logger = logging.GetLogger("base_service")

func registerService() {
	baseConf.Init()

	config := conf.GetConfig()

	// 使用etcd注册
	micReg := etcd.NewRegistry()

	// 新建服务
	service := micro.NewService(
		micro.Name("go.micro.service.base"),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*time.Duration(config.Service.RegisterTTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Service.RegisterInterval)),
	)
	// 服务初始化
	service.Init(
		micro.Action(func(c *cli.Context) error {
			handler.Init()
			router.Init()
			return nil
		}),
	)

	// 启动服务
	if err := service.Run(); err != nil {
		logger.Panicf("id generator start failed, for the reason:%v", err.Error())
	}
}

func main() {
	registerService()
}
