package handler

import (
	"encoding/gob"
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var (
	dbService   db.Service
	authService auth.Service
	logger      = logging.GetLogger("base_service_handler")
)

type UserSessionInfo struct {
	UserID    int64  `json:"user_id"`
	CSRFToken string `json:"csrf_token"`
}

const (
	defaultSessionKey = "user"
)

func Init() {
	config := conf.GetConfig()
	micReg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = config.Etcd.EndPoints
	})
	// 新建服务
	service := micro.NewService(
		micro.Registry(micReg),
	)

	dbService = db.NewService("go.micro.service.db", service.Client())
	authService = auth.NewService("go.micro.service.auth", service.Client())

	// 注册 user session struct
	gob.Register(UserSessionInfo{})
}
