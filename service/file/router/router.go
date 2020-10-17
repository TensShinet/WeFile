package router

import (
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/file/handler"
	"github.com/gin-gonic/gin"
)

var logger = logging.GetLogger("file_service_router")

func Init() {
	config := conf.GetConfig()
	if config == nil {
		logger.Panic("config is nil")
		return
	}
	router := gin.Default()
	// 允许跨域
	router.Use(handler.CorsMiddleware())

	v1 := router.Group("/api/v1")

	uploadAPI := v1.Group("/upload", handler.UploadAuthorizeMiddleware())
	uploadAPI.POST("", handler.Upload)

	downloadAPI := v1.Group("/download", handler.DownloadAuthorizeMiddleware())
	downloadAPI.GET("", handler.Download)
	if err := router.Run(config.FileAPI.Address); err != nil {
		logger.Panicf("router run failed, for the reason:%v", err)
	}
}
