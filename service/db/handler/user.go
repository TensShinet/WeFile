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
	logger.Infof("InsertUser id:%v", id)
	u := req.User
	modelUser := &model.User{}
	if err = db.Debug().Where("email = ?", u.Email).Take(modelUser).Error; err == gorm.ErrRecordNotFound {
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
			logger.Infof("InsertUser failed, for the reason:%v", err)
			res.Err = getProtoError(err, common.DBConflictCode)
			return nil
		} else {
			// 成功创建
			res.Id = id
		}
	} else if err != nil {
		return err
	} else {
		// 用户已经存在
		logger.Infof("InsertUser failed, for the reason:%v", common.ErrConflict)
		res.Err = getProtoError(common.ErrConflict, common.DBConflictCode)
		return nil
	}

	return nil
}

// 查询一个用户
//
// 检查这个用户存不存在 不存在 返回 DBNotFoundCode
// 存在就返回
func (s *Service) QueryUser(ctx context.Context, req *proto.QueryUserReq, res *proto.QueryUserResp) error {

	var (
		err error
	)

	u := model.User{}
	err = db.Where("id=? or email=?", req.Id, req.Email).First(&u).Error
	logger.Infof("QueryUser id=%v email=%v", req.Id, req.Email)
	if err != nil {
		logger.Errorf("QueryUser failed, for the reason:%v", err)
		if err == gorm.ErrRecordNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		} else {
			return err
		}
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
	return nil
}

// 删除一个用户
//
// 查询该用户存不存在 不存在返回 DBNotFoundCode
// 存在就上锁 删掉这个用户
// 事务中所有的 err 统一在外面处理
func (s *Service) DeleteUser(ctx context.Context, req *proto.DeleteUserReq, res *proto.DeleteUserResp) error {
	var (
		err      error
		affected int64
	)
	u := &model.User{}
	logger.Infof("DeleteUser id=%v", req.Id)
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.Id).First(u).Error; err != nil {
			return err
		}
		affected = tx.Delete(u).RowsAffected
		return nil
	})
	if err != nil {
		logger.Errorf("DeleteUser failed, for the reason:%v", err)
		return err
	}
	if affected == 0 {
		logger.Errorf("DeleteUser failed, for the reason:%v", gorm.ErrRecordNotFound)
		res.Err = getProtoError(gorm.ErrRecordNotFound, common.DBNotFoundCode)
		return nil
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
	return err
}
