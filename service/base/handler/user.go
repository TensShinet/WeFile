// WeFile API
//
//     Schemes: http, https
//     BasePath: /api/v1
//     Version: 0.0.1
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package handler

import (
	"github.com/gin-gonic/gin"
)

// swagger:model NewUser
type NewUser struct {
	// 用户名
	//
	// minimum length: 1
	// maximum length: 64
	// Required: true
	Name string `json:"name"`
	// 密码
	//
	// minimum length: 8
	// maximum length: 64
	// Required: true
	Password string `json:"password"`
	// 邮箱
	//
	// 邮箱唯一
	//
	// Required: true
	Email string `json:"email"`
	// 验证码
	//
	// 邮件拿验证码
	//
	// Required: true
	VerifyCode string `json:"verify_code"`
}

// swagger:model User
type User struct {
	// 用户 id
	UserID int64 `json:"user_id"`
	// 角色 id
	RoleID int `json:"role_id"`
	// 角色名
	RoleName string `json:"role_name"`
	// 用户名
	Name string `json:"name"`
	// 邮箱
	Email string `json:"email"`
	// 个人简介
	Profile string `json:"profile"`
}

// swagger:parameters SignUp
type SignUpParam struct {
	// 新用户
	// in: body
	// required: true
	User *NewUser `json:"user"`
}

// swagger:parameters SignIn
type SignInParam struct {
	// in: body
	Body struct {
		// 邮箱
		//
		// Required: true
		Email string `json:"email"`
		// 密码
		//
		// Required: true
		Password string `json:"password"`
	} `json:"body"`
}

// User
// swagger:response UserResponse
type UserResponse struct {
	Cookie string `json:"cookie"`
	// in: body
	User *User `json:"user"`
}

// swagger:route POST /user/sign_up User SignUp
//
// SignUp
//
// 用户注册
//     Responses:
//       200: UserResponse
//		 409: ConflictError
//       500: ErrorResponse
func SignUp(c *gin.Context) {}

// swagger:route POST /user/sign_in User SignIn
//
// SignIn
//
// 用户登录
//     Responses:
//       200: UserResponse
//		 403: ForbiddenResponse
//       500: ErrorResponse
func SignIn(c *gin.Context) {}
