package handler

import (
	"fmt"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// TODO: 生成 csrf token 更好的办法
func getCSRFToken() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(1 << 30))
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			id int64
			//CSRFToken string
			err error
		)
		defer func() {
			if err != nil {
				c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: err.Error()})
				c.Abort()
			}
		}()
		session := sessions.Default(c)
		u, ok := session.Get(defaultUserKey).(UserSessionInfo)
		if !ok {
			logger.Debugf("id:%v UserID:%v", id, u.UserID)
			err = fmt.Errorf("valid session")
			return
		}

		// 检查 csrf token
		if c.Request.Method == "DELETE" {
			if c.Query("csrf_token") != u.CSRFToken {
				err = fmt.Errorf("valid csrf_token")
				return
			}
		} else if c.Request.Method == "POST" {
			if c.Request.FormValue("csrf_token") != u.CSRFToken {
				err = fmt.Errorf("valid csrf_token")
				return
			}
		}

		c.Set(defaultUserKey, &User{
			UserID: u.UserID,
		})

		c.Next()
	}
}

func setSession(c *gin.Context, userID int64, key, csrfToken string) error {
	session := sessions.Default(c)
	session.Set(defaultUserKey, UserSessionInfo{
		UserID:    userID,
		CSRFToken: csrfToken,
	})
	if err := session.Save(); err != nil {
		return err
	}
	u, _ := session.Get(defaultUserKey).(UserSessionInfo)
	logger.Debugf("setSession sessionID: %v, UserID:%v, CSRFToken:%v", key, u.UserID, u.CSRFToken)
	return nil
}

func getUser(c *gin.Context) *User {
	val, _ := c.Get(defaultUserKey)
	user, _ := val.(*User)
	return user
}

func checkUserInGroup(c *gin.Context, userID, groupID int64) error {
	var (
		err error
		res *db.CheckUserInGroupResp
	)
	if res, err = dbService.CheckUserInGroup(c, &db.UserIDGroupID{UserID: userID, GroupID: groupID}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	if res.Err != nil && res.Err.Code == common.DBNotFoundCode {
		common.SetSimpleResponse(c, http.StatusForbidden, res.Err.Message)
		return err
	}

	return nil
}

func getGroupIDFromContext(c *gin.Context) (id int64, err error) {
	if c.Request.Method == "POST" {
		if id, err = utils.ParseInt64(c.Request.FormValue("group_id")); err != nil {
			common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
			return 0, err
		} else if id == 0 {
			common.SetSimpleResponse(c, http.StatusBadRequest, "invalid group id")
			return 0, fmt.Errorf("invalid group id")
		}
		return id, nil

	} else {
		if id, err = utils.ParseInt64(c.Query("group_id")); err != nil {
			common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
			return 0, err
		} else if id == 0 {
			common.SetSimpleResponse(c, http.StatusBadRequest, "invalid group id")
			return 0, fmt.Errorf("invalid group id")
		}
		return id, nil
	}
}

func getFileIDFromContext(c *gin.Context) (id int64, err error) {
	if c.Request.Method == "POST" {
		if id, err = utils.ParseInt64(c.Request.FormValue("file_id")); err != nil {
			common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
			return 0, err
		} else if id == 0 {
			common.SetSimpleResponse(c, http.StatusBadRequest, "invalid group id")
			return 0, fmt.Errorf("invalid file id")
		}
		return id, nil

	} else {
		if id, err = utils.ParseInt64(c.Query("file_id")); err != nil {
			common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
			return 0, err
		} else if id == 0 {
			common.SetSimpleResponse(c, http.StatusBadRequest, "invalid file id")
			return 0, fmt.Errorf("invalid file id")
		}
		return id, nil
	}
}
