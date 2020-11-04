package handler

import (
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"net/http"
)

const defaultAuthKey = "auth"

// 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Max-Age", "10000")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept")
		c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
		// 允许跨域 OPTIONS 直接返回
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusNoContent, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// upload jwt 认证中间件
func UploadAuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		a, ok := c.Request.Header["Authorization"]
		if !ok || len(a) < 1 || len(a[0]) < len("Bearer ") || a[0][:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: "authorization failed"})
			c.Abort()
			return
		}
		jwtToken := a[0][7:]
		res, err := authService.UploadJWTDecode(c, &auth.DecodeReq{Token: jwtToken})
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: "authorization failed"})
			c.Abort()
			return
		}

		spew.Dump("res.FileMeta ", res.FileMeta)

		c.Set(defaultAuthKey, res.FileMeta)
		c.Next()
	}
}

// 更新 Authorize 信息
func UpdateAuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get(defaultAuthKey)
		fileMeta, _ := val.(*auth.UploadFileMeta)

		res, _ := authService.UploadJWTEncode(c, fileMeta)
		logger.Debugf("UpdateAuthorizeMiddleware token:%v", res.Token)
		c.Header("Authorization", "Bearer "+res.Token)
		c.Next()
	}
}

// download jwt 认证中间件
func DownloadAuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		a, ok := c.Request.Header["Authorization"]
		if !ok || len(a) < 1 || len(a[0]) < len("Bearer ") || a[0][:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: "authorization failed"})
			c.Abort()
			return
		}
		jwtToken := a[0][7:]
		res, err := authService.DownloadJWTDecode(c, &auth.DecodeReq{Token: jwtToken})
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: "authorization failed"})
			c.Abort()
			return
		}
		c.Set(defaultAuthKey, res.FileMeta)
		c.Next()
	}
}
