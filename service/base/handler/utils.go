package handler

import (
	"fmt"
	"github.com/TensShinet/WeFile/service/common"
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
			id  int64
			err error
		)
		defer func() {
			if err != nil {
				c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse{Message: err.Error()})
				c.Abort()
			}
		}()
		if id, err = strconv.ParseInt(c.Param("user_id"), 10, 64); err != nil {
			logger.Debug("Authorize ", c.Param("user_id"))
			return
		}
		session := sessions.Default(c)
		u, ok := session.Get(defaultSessionKey).(UserSessionInfo)
		if !ok || u.UserID != id {
			logger.Debugf("id:(%v) UserID:(%v)", id, u.UserID)
			err = fmt.Errorf("valid session")
			return
		}

		c.Next()
	}
}

func setSession(c *gin.Context, userID int64, key, csrfToken string) error {
	session := sessions.Default(c)
	session.Set(defaultSessionKey, UserSessionInfo{
		UserID:    userID,
		CSRFToken: csrfToken,
	})
	if err := session.Save(); err != nil {
		return err
	}
	u, _ := session.Get(defaultSessionKey).(UserSessionInfo)
	logger.Debugf("setSession sessionID: %v, UserID:%v, CSRFToken:%v", key, u.UserID, u.CSRFToken)
	return nil
}
