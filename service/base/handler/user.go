// WeFile API
//
//     Schemes: http, https
//     BasePath: /api/v1
//     Version: 0.0.1
//
// swagger:meta
package handler

import (
	"github.com/gin-gonic/gin"
)

// swagger:parameters SignUp
type SignUpParam struct {
	// 用户名 长度小于 64
	//
	// Required: true
	// in: query
	Name string
	// 密码 长度小于 64
	//
	// Required: true
	// in: query
	Password string
	// 邮箱 需要检查格式是否正确
	//
	// Required: true
	// in: query
	Email string
	// 验证码 邮件拿验证码
	//
	// Required: true
	// in: query
	VerifyCode string
	// 个人介绍
	//
	// Required: false
	// in: query
	Profile string
}

// swagger:response User
type User struct {
	ID      string
	Name    string
	Email   string
	Profile string
}

// swagger:route POST /user/sign_up User SignUp
//
// SignUp
//
// 注册用户
//     Responses:
//       200: User
//		 409: ConflictError
//       500: ErrorResponse
func SignUp(c *gin.Context) {}

// swagger:route POST /user/sign_in User SignIn
//
// SignIn
//
// 用户登录
//     Responses:
//       200: User
//		 403: ForbiddenResponse
//       500: ErrorResponse
func SignIn(c *gin.Context) {}
