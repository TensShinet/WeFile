package main

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/id_generator/handler"
	proto "github.com/TensShinet/WeFile/service/id_generator/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var logger = logging.GetLogger("id_generator")

func main() {
	config := conf.GetConfig("conf.yml")
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Name("go.micro.service.id_generator"),
		micro.Registry(micReg),
		micro.Version("latest"),
	)
	// 服务初始化
	service.Init()

	// 注册服务
	if err := proto.RegisterGenerateIDServiceHandler(service.Server(), new(handler.Service)); err != nil {
		logger.Panicf("register failed, for the reason:%v", err.Error())
	}
	// 启动服务
	if err := service.Run(); err != nil {
		logger.Panicf("id generator start failed, for the reason:%v", err.Error())
	}
}
