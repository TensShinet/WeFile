package handler

import (
	"github.com/TensShinet/WeFile/service/base/conf"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// swagger:parameters CreateGroup
type CreateGroupParam struct {
	// 存储 session id
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// in: body
	Body struct {
		// group 名字
		//
		// minimum length: 1
		// maximum length: 64
		// Required: true
		Name string `json:"name"`
		// group 密码 用于检验 能否加入
		//
		// minimum length: 8
		// maximum length: 64
		// Required: true
		Password string `json:"password"`
		// csrf_token
		// 登录的时候才会更新
		// Required: true
		CSRFToken string `json:"csrf_token"`
	}
}

// swagger:model Group
type Group struct {
	GroupID   int64     `json:"group_id"`
	Name      string    `json:"name"`
	OwnerID   int64     `json:"owner_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GroupResponse
//
// swagger:response GroupResponse
type GroupResponse struct {
	// in: body
	group *Group
}

// swagger:route POST /user/group Group CreateGroup
//
// CreateGroup
//
// 创建 group
//     Responses:
//		 200: GroupResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//       500: ErrorResponse
func CreateGroup(c *gin.Context) {

	user := getUser(c)
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	logger.Infof("CreateGroup userID:%v name:%v password:%v", user.UserID, name, password)

	if len(name) > 64 || len(name) < 1 || len(password) < 8 || len(password) > 64 {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "账号密码长度不符"})
		return
	}
	config := conf.GetConfig()
	password = utils.Digest256([]byte(password + config.BaseAPI.Salt))

	var (
		err error
		res *db.CreateGroupResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroup failed, for the reason:%v", err)
		}
	}()

	if res, err = dbService.CreateGroup(c, &db.CreateGroupReq{
		Group: &db.Group{
			OwnerID:   user.UserID,
			Name:      name,
			Password:  password,
			CreatedAt: time.Now().Unix(),
			Status:    0,
		},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "Server Error"})
		return
	}

	c.JSON(http.StatusOK, &Group{
		GroupID:   res.Group.Id,
		Name:      res.Group.Name,
		OwnerID:   res.Group.OwnerID,
		Status:    int(res.Group.Status),
		CreatedAt: time.Unix(res.Group.CreatedAt, 0),
	})
}

// swagger:parameters GetGroup GetGroupMemberList DeleteGroup
type GetGroupParam struct {
	// 存储 session id
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// group id
	//
	// Required: true
	GroupID int64 `json:"group_id"`
	// csrf_token
	//
	// 登录的时候才会更新
	//
	// GET 方法不需要，其他都需要
	CSRFToken string `json:"csrf_token"`
}

// swagger:route GET /user/group Group GetGroup
//
// GetGroup
//
// 得到创建信息
//     Responses:
//       200: GroupResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//		 403: ForbiddenResponse
//       500: ErrorResponse
func GetGroup(c *gin.Context) {
	var (
		groupID int64
		userID  int64
		err     error

		res1 *db.GroupResp
	)

	defer func() {
		if err != nil {
			logger.Errorf("GetGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID

	logger.Infof("CreateGroup userID:%v groupID:%v", userID, groupID)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res1, err = dbService.QueryGroup(c, &db.UserIDGroupID{GroupID: groupID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res1.Err != nil && res1.Err.Code == common.DBNotFoundCode {
		common.SetSimpleResponse(c, http.StatusNotFound, res1.Err.Message)
		return
	}

	c.JSON(http.StatusOK, &Group{
		GroupID:   res1.Group.Id,
		Name:      res1.Group.Name,
		OwnerID:   res1.Group.OwnerID,
		Status:    int(res1.Group.Status),
		CreatedAt: time.Unix(res1.Group.CreatedAt, 0),
	})

}

// swagger:model GroupUserInfo
type GroupUserInfo struct {
	UserID  int64     `json:"user_id"`
	GroupID int64     `json:"group_id"`
	JoinAt  time.Time `json:"join_at"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
}

// GroupMemberListResponse
//
// swagger:response GroupMemberListResponse
type GroupMemberListResponse struct {
	// in: body
	Users []*GroupUserInfo `json:"users"`
}

// swagger:route GET /user/group/member_list Group GetGroupMemberList
//
// GetGroupMemberList
//
// 得到组成员
//     Responses:
//       200: GroupMemberListResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//		 403: ForbiddenResponse
//       500: ErrorResponse
func GetGroupMemberList(c *gin.Context) {
	var (
		groupID int64
		userID  int64
		err     error
		res1    *db.ListGroupUserResp
	)

	defer func() {
		if err != nil {
			logger.Errorf("GetGroupMemberList failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID

	logger.Infof("GetGroupMemberList userID:%v groupID:%v", userID, groupID)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res1, err = dbService.ListGroupUser(c, &db.UserIDGroupID{GroupID: groupID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res1.Err != nil && res1.Err.Code == http.StatusNotFound {
		common.SetSimpleResponse(c, http.StatusNotFound, res1.Err.Message)
		return
	}
	users := make([]*GroupUserInfo, len(res1.Users))
	for i, u := range res1.Users {
		users[i] = &GroupUserInfo{
			UserID:  u.UserID,
			GroupID: u.GroupID,
			JoinAt:  time.Unix(u.JoinAt, 0),
			Email:   u.Email,
			Name:    u.Name,
		}
	}
	c.JSON(http.StatusOK, users)
}

// swagger:route Delete /user/group Group DeleteGroup
//
// DeleteGroup
//
// 删除一个组，只允许创建者删除
//     Responses:
//       200: GroupResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//       500: ErrorResponse
func DeleteGroup(c *gin.Context) {
	var (
		groupID int64
		userID  int64
		err     error
		res1    *db.GroupResp
	)

	defer func() {
		if err != nil {
			logger.Errorf("DeleteGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID

	logger.Infof("DeleteGroup userID:%v groupID:%v", userID, groupID)
	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res1, err = dbService.DeleteGroup(c, &db.UserIDGroupID{GroupID: groupID, UserID: userID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res1.Err != nil && res1.Err.Code == common.DBForbiddenCode {
		common.SetSimpleResponse(c, http.StatusForbidden, res1.Err.Message)
		return
	}

	c.JSON(http.StatusOK, &Group{
		GroupID:   groupID,
		Name:      res1.Group.Name,
		OwnerID:   res1.Group.OwnerID,
		Status:    int(res1.Group.Status),
		CreatedAt: time.Unix(res1.Group.CreatedAt, 0),
	})

}
