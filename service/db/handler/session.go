package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"time"
)

// TODO: redis 加速
// session 相关服务
func (s *Service) InsertSession(ctx context.Context, req *proto.InsertSessionReq, res *proto.InsertSessionResp) error {
	id, err := getID(ctx)
	if err != nil {
		logger.Error("generateIDService failed, for the reason:&v", err.Error())
		return err
	}
	session := req.Session
	err = db.Create(&model.Session{
		ID:        id,
		Token:     session.Token,
		UserID:    session.UserID,
		CreateAt:  time.Unix(session.CreatedAt, 0),
		ExpireAt:  time.Unix(session.ExpireAt, 0),
		CSRFToken: session.CSRFToken,
	}).Error
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
	}
	return err
}
func (s *Service) GetUserSession(ctx context.Context, req *proto.GetUserSessionReq, res *proto.GetUserSessionResp) (err error) {
	session := model.Session{}
	err = db.Where("user_id=?", req.UserID).First(&session).Error
	if err != nil {
		res.Err = &proto.Error{Code: -1, Message: err.Error()}
		if err == gorm.ErrRecordNotFound {
			res.Err.Code = common.DBNotFoundCode
		}
		return
	}
	res.Session = &proto.Session{
		UserID:    session.UserID,
		Token:     session.Token,
		CreatedAt: session.CreateAt.Unix(),
		ExpireAt:  session.ExpireAt.Unix(),
		CSRFToken: session.CSRFToken,
		SessionID: session.ID,
	}

	return
}

func (s *Service) DeleteUserSession(ctx context.Context, req *proto.DeleteUserSessionReq, res *proto.DeleteUserSessionResp) error {
	return nil
}
