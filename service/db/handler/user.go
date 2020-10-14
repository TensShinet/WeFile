package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"time"
)

// 插入一个用户
//
// 检查是不是存在这个用户
// 1. 插入之前就检查到
// 2. 插入之后才检查到
// 插入一个用户
func (s *Service) InsertUser(ctx context.Context, req *proto.InsertUserReq, res *proto.InsertUserResp) error {
	var (
		id  int64
		err error
	)
	if id, err = getID(ctx); err != nil {
		logger.Error("generateIDService failed, for the reason:&v", err.Error())
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	}
	u := req.User
	modelUser := &model.User{}
	if err = db.Where("email = ?", u.Email).Take(modelUser).Error; err == gorm.ErrRecordNotFound {
		if err = db.Create(&model.User{
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
		}).Error; err != nil {
			// 并发创建，冲突
			res.Err = getProtoError(err, common.DBConflictCode)
			return err
		} else {
			// 成功创建
			res.Id = id
		}
	} else if err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	} else {
		// 用户已经存在
		res.Err = getProtoError(err, common.DBConflictCode)
	}

	return nil
}

// 查询一个用户
//
// 检查这个用户存不存在 不存在 返回 DBNotFoundCode
// 存在就返回
func (s *Service) QueryUser(ctx context.Context, req *proto.QueryUserReq, res *proto.QueryUserResp) (err error) {
	u := model.User{}
	err = db.Where("id=? or email=?", req.Id, req.Email).First(&u).Error
	if err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		// 没找到就用 DBNotFoundCode
		if err == gorm.ErrRecordNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
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

// 删除一个用户
//
// 查询该用户存不存在 不存在返回 DBNotFoundCode
// 存在就上锁 删掉这个用户
// 事务中所有的 err 统一在外面处理
func (s *Service) DeleteUser(ctx context.Context, req *proto.DeleteUserReq, res *proto.DeleteUserResp) (err error) {
	u := &model.User{}
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.Id).First(u).Error; err != nil {
			return err
		}
		err = tx.Delete(u).Error
		return err
	})
	if err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		if err == gorm.ErrRecordNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
		}
		return
	}
	res.User = &proto.User{
		Id:             u.ID,
		RoleID:         u.RoleID,
		Name:           u.Name,
		Password:       u.Password,
		Email:          u.Email,
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
