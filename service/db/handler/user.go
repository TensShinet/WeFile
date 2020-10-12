package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"time"
)

func (s *Service) InsertUser(ctx context.Context, req *proto.InsertUserReq, res *proto.InsertUserResp) (err error) {
	id, err := getID(ctx)
	if err != nil {
		logger.Error("generateIDService failed, for the reason:&v", err.Error())
		return
	}
	u := req.User
	err = db.Create(&model.User{
		ID:             id,
		RoleID:         u.RoleID,
		Name:           u.Name,
		Password:       u.Password,
		Email:          u.Email,
		Phone:          u.Phone,
		EmailValidated: u.EmailValidated,
		PhoneValidated: u.PhoneValidated,
		SignUpAt:       time.Unix(u.SignUpAt, 0),
		LastActiveAt:   time.Unix(u.LasActiveAt, 0),
		Profile:        u.Profile,
		Status:         int(u.Status),
	}).Error
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
	}
	return
}
func (s *Service) QueryUser(ctx context.Context, req *proto.QueryUserReq, res *proto.QueryUserResp) (err error) {
	u := model.User{}
	err = db.Where("user_id=? or email=?", req.Id, req.Email).First(&u).Error
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
		// 没找到就用 DBNotFoundCode
		if err == gorm.ErrRecordNotFound {
			res.Err.Code = common.DBNotFoundCode
		}
		return
	}
	res.User = &proto.User{
		Id:             u.ID,
		RoleID:         u.RoleID,
		Name:           u.Name,
		Password:       u.Password,
		Email:          u.Password,
		Phone:          u.Phone,
		EmailValidated: u.EmailValidated,
		PhoneValidated: u.PhoneValidated,
		SignUpAt:       u.SignUpAt.Unix(),
		LasActiveAt:    u.LastActiveAt.Unix(),
		Profile:        u.Profile,
		Status:         int32(u.Status),
	}
	return
}
