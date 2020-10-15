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
	"encoding/base64"
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	// csrf Token
	CSRFToken string `json:"csrf_token"`
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
// 		 400: BadRequestResponse
//		 409: ConflictError
//       500: ErrorResponse
func SignUp(c *gin.Context) {
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")
	logger.Debugf("SignUp %v %v %v", name, password, email)
	// TODO: 复杂检查
	if len(name) > 64 || len(name) < 1 || len(password) < 8 || len(password) > 64 {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "账号密码长度不符"})
		return
	}
	config := conf.GetConfig()
	password = base64.StdEncoding.EncodeToString(digest256(password + config.BaseAPI.Salt))
	var (
		err  error
		res1 *db.InsertUserResp
	)

	if res1, err = dbService.InsertUser(c, &db.InsertUserReq{
		User: &db.User{
			RoleID:         common.GeneralUserRoleID,
			Name:           name,
			Password:       password,
			Email:          email,
			EmailValidated: false,
			PhoneValidated: false,
			SignUpAt:       time.Now().Unix(),
			LasActiveAt:    time.Now().Unix(),
		},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "Server Error"})
		return
	}
	if res1.Err != nil && res1.Err.Code == common.DBConflictCode {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "邮箱已被注册"})
		return
	}

	// 注册成功
	csrfToken := getCSRFToken()
	if err := setSession(c, res1.Id, defaultSessionKey, csrfToken); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, User{
		UserID:    res1.Id,
		RoleID:    common.GeneralUserRoleID,
		RoleName:  "普通用户",
		Name:      name,
		Email:     email,
		CSRFToken: csrfToken,
	})
}

// swagger:route POST /user/sign_in User SignIn
//
// SignIn
//
// 用户登录
//     Responses:
//       200: UserResponse
//		 403: ForbiddenResponse
//       500: ErrorResponse
func SignIn(c *gin.Context) {
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")
	logger.Debug("SignIn ", password, email)
	var (
		res1 *db.QueryUserResp
		err  error
	)
	// 查询用户
	if res1, err = dbService.QueryUser(c, &db.QueryUserReq{
		Email: email,
	}); err != nil {
		// 登录失败
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}
	config := conf.GetConfig()

	if res1.Err != nil && res1.Err.Code == common.DBNotFoundCode {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse{
			Message: "该邮箱未注册",
		})
		return
	}

	if res1.User.Password != base64.StdEncoding.EncodeToString(digest256(password+config.BaseAPI.Salt)) {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse{
			Message: "账户密码错误",
		})
		return
	}

	// 登录成功
	csrfToken := getCSRFToken()
	if err := setSession(c, res1.User.Id, defaultSessionKey, csrfToken); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}

	u := res1.User

	c.JSON(http.StatusOK, User{
		UserID:    u.Id,
		RoleID:    int(u.RoleID),
		RoleName:  common.GeneralUserRoleName,
		Name:      u.Name,
		Email:     u.Email,
		Profile:   u.Profile,
		CSRFToken: csrfToken,
	})

}
