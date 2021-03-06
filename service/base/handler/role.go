package handler

import (
	"github.com/TensShinet/WeFile/service/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// swagger:model Role
type Role struct {
	// 角色 id
	RoleID int64 `json:"role_id"`
	// 角色名
	Name string `json:"name"`
}

// Role List
//
// swagger:response RoleList
type RoleList struct {
	// in: body
	Roles []*Role `json:"body"`
}

// swagger:route GET /role/all Role GetAllRoles
//
// GetAllRoles
//
// 获得所有角色
//     Responses:
//       200: RoleList
//       500: ErrorResponse
func GetAllRoles(c *gin.Context) {
	c.JSON(http.StatusOK, []Role{
		{
			RoleID: common.GeneralUserRoleID,
			Name:   common.GeneralUserRoleName,
		},
		{
			RoleID: common.AdminRoleID,
			Name:   common.AdminRoleName,
		},
	})
}
