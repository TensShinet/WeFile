package main

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/auth/conf"
	"github.com/TensShinet/WeFile/service/auth/handler"
	"github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"time"
)

var logger = logging.GetLogger("auth_service")

func startRPCService() {
	config := conf.GetConfig()
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Service.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Name("go.micro.service.auth"),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*time.Duration(config.Service.RegisterTTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Service.RegisterInterval)),
	)
	// 服务初始化
	service.Init()

	// 注册服务
	if err := proto.RegisterServiceHandler(service.Server(), new(handler.Service)); err != nil {
		logger.Panicf("register failed, for the reason:%v", err.Error())
	}
	// 启动服务
	if err := service.Run(); err != nil {
		logger.Panicf("id generator start failed, for the reason:%v", err.Error())
	}
}

// jwt 相关
// jwt encode: 将一些关键信息 encode 之后返回
// jwt decode: 将一些关键信息 decode 之后返回，判断是否过期，判断签名是否一致
func main() {
	// 初始化配置
	conf.Init("auth_conf.yml")
	startRPCService()
}
