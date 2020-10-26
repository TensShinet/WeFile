package handler

import (
	"github.com/TensShinet/WeFile/logging"
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/service/file/conf"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var (
	dbService   db.Service
	authService auth.Service
	fileStore   store.Store
	logger      = logging.GetLogger("file_service_handler")
)

func Init() {
	config := conf.GetConfig()
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Service.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Registry(micReg),
	)

	dbService = db.NewService("go.micro.service.db", service.Client())
	authService = auth.NewService("go.micro.service.auth", service.Client())

	var (
		err error
	)
	if fileStore, err = store.NewLocalStore(config.FileAPI.LocalTempStore, config.SamplingChunkSize); err != nil {
		logger.Panicf("init failed, for the reason:%v", err)
	}
}
