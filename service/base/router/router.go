package router

import (
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/base/conf"
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
		MaxAge: config.BaseAPI.SessionMaxAge * 60 * 60,
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
	v1.GET("/user/file_list", handler.GetUserFileList)
	v1.DELETE("/user/file_list", handler.DeleteUserFile)
	v1.POST("/user/file_list", handler.CreateDirectory)
	v1.GET("/user/download_address", handler.GetDownloadAddress)
	v1.GET("/user/upload_address", handler.GetUploadAddress)

	// user group 相关
	// 创建 group
	v1.POST("/user/group", handler.CreateGroup)
	v1.DELETE("/user/group", handler.DeleteGroup)
	v1.GET("/user/group", handler.GetGroup)
	v1.GET("/user/group/member_list", handler.GetGroupMemberList)

	// 获取所有 group
	v1.GET("/user/group_list", handler.GetGroupList)
	v1.POST("/user/group_list", handler.JoinGroup)
	v1.DELETE("/user/group_list", handler.LeaveGroup)

	// group file 相关
	v1.DELETE("/user/group/file_list", handler.DeleteGroupFile)
	v1.POST("/user/group/file_list", handler.CreateGroupDirectory)
	v1.GET("/user/group/file_list", handler.GetGroupFileList)
	v1.GET("/user/group/download_address", handler.GetGroupDownloadAddress)
	v1.GET("/user/group/upload_address", handler.GetGroupUploadAddress)

	if err := router.Run(config.BaseAPI.Address); err != nil {
		logger.Panicf("router run failed, for the reason:%v", err)
	}
}
