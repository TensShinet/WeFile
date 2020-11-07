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

// swagger:response GetGroupListResponse
type GetGroupListResponse struct {
	// in: body
	Groups []*Group `json:"groups"`
}

// GetGroupList
//
// swagger:parameters GetGroupList
type GetGroupListParam struct {
	// 存储 session id
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
}

// swagger:route GET /user/group_list Group GetGroupList
//
// GetGroupList
//
// 得到用户的所有组
//     Responses:
//       200: GetGroupListResponse
//  	 401: UnauthorizedResponse
//       500: ErrorResponse
func GetGroupList(c *gin.Context) {
	var (
		userID int64
		err    error
		res    *db.ListUserGroupResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroupList failed, for the reason:%v", err)
		}
	}()

	userID = getUser(c).UserID
	logger.Info("GetGroupList userID:%v", userID)

	if res, err = dbService.ListUserGroup(c, &db.UserIDGroupID{UserID: userID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	groups := make([]*Group, len(res.Groups))
	for i, v := range res.Groups {
		groups[i] = &Group{
			GroupID:   v.Id,
			Name:      v.Name,
			OwnerID:   v.OwnerID,
			Status:    int(v.Status),
			CreatedAt: time.Unix(v.CreatedAt, 0),
		}
	}
	c.JSON(http.StatusOK, groups)

}

// swagger:parameters JoinGroup
type JoinGroupParam struct {
	// 存储 session id
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// in: body
	Body struct {
		// required: true
		GroupID int64 `json:"group_id"`
		// required: true
		Password string `json:"password"`
	}
}

// swagger:route POST /user/group_list Group JoinGroup
//
// JoinGroup
//
// 加入一个组
//     Responses:
//       200: GroupResponse
//  	 401: UnauthorizedResponse
//		 403: ForbiddenResponse
//		 404: NotFoundError
//		 409: ConflictError
//       500: ErrorResponse
func JoinGroup(c *gin.Context) {
	var (
		userID   int64
		groupID  int64
		err      error
		password string
		res      *db.GroupResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("JoinGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}

	userID = getUser(c).UserID
	password = utils.Digest256([]byte(c.Request.FormValue("password") + conf.GetConfig().BaseAPI.Salt))
	logger.Info("JoinGroup userID:%v groupID:%v password:%v", userID, groupID, password)

	if res, err = dbService.QueryGroup(c, &db.UserIDGroupID{GroupID: groupID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil {
		common.SetSimpleResponse(c, int(res.Err.Code), res.Err.Message)
		return
	}

	if password != res.Group.Password {
		common.SetSimpleResponse(c, http.StatusForbidden, "password or group_id error")
		return
	}

	if res, err = dbService.JoinGroup(c, &db.UserIDGroupID{
		UserID:  userID,
		GroupID: groupID,
	}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
	}

	if res.Err != nil {
		common.SetSimpleResponse(c, int(res.Err.Code), res.Err.Message)
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

// swagger:parameters LeaveGroup
type LeaveGroupParam struct {
	// 存储 session id
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// in: query
	// required: true
	GroupID int64 `json:"group_id"`
}

// swagger:route DELETE /user/group_list Group LeaveGroup
//
// QuitGroup
//
// 退出一个组
//     Responses:
//       200: GroupResponse
//  	 401: UnauthorizedResponse
//		 404: NotFoundError
//       500: ErrorResponse
func LeaveGroup(c *gin.Context) {
	var (
		userID  int64
		groupID int64
		err     error
		res     *db.GroupResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("QuitGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID
	logger.Info("QuitGroup userID:%v groupID:%v", userID, groupID)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res, err = dbService.LeaveGroup(c, &db.UserIDGroupID{
		UserID:  userID,
		GroupID: groupID,
	}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil {
		common.SetSimpleResponse(c, int(res.Err.Code), res.Err.Message)
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
