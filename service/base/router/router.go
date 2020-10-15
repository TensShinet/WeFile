package router

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/base/handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var logger = logging.GetLogger("base_service_router")

func Init() {
	config := conf.GetConfig()
	if config == nil {
		logger.Panic("config is nil")
		return
	}
	router := gin.Default()

	// 使用 session 中间件
	store, err := redis.NewStore(10, config.Redis.Network, config.Redis.Conn, config.Redis.Password, []byte(config.BaseAPI.SessionSecrete))
	if err != nil {
		logger.Panicf("init session store failed, for the reason:%v", err)
		return
	}
	store.Options(sessions.Options{
		MaxAge: config.BaseAPI.SessionMaxAge * 60,
	})

	router.Use(sessions.Sessions(config.BaseAPI.SessionName, store))

	v1 := router.Group("/api/v1")
	// user 相关
	v1.POST("/user/sign_in", handler.SignIn)
	v1.POST("/user/sign_up", handler.SignUp)

	// role 相关
	v1.GET("/role/all", handler.GetAllRoles)

	// 用户认证
	v1.Use(handler.Authorize())
	v1.GET("/file_list/:user_id", handler.GetUserFileList)
	v1.GET("/download_address/:user_id", handler.GetDownloadAddress)
	v1.GET("/upload_address/:user_id", handler.GetUploadAddress)

	if err := router.Run(config.BaseAPI.Address); err != nil {
		logger.Panicf("router run failed, for the reason:%v", err)
	}
}
