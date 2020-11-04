package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"time"
)

func (s *Service) CreateGroup(ctx context.Context, req *proto.CreateGroupReq, res *proto.CreateGroupResp) error {
	var (
		id  int64
		err error
	)
	if id, err = getID(ctx); err != nil {
		logger.Error("CreateGroup failed, for the reason:&v", err.Error())
		return err
	}
	g := req.Group
	logger.Infof("CreateGroup id:%v ownerID:%v name:%v", id, g.OwnerID, g.Name)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		// 创建 group
		if err = tx.Create(&model.Group{
			ID:        id,
			OwnerID:   g.OwnerID,
			Name:      g.Name,
			Password:  g.Password,
			CreatedAt: time.Unix(g.CreatedAt, 0),
			Status:    int(g.Status),
		}).Error; err != nil {
			logger.Errorf("CreateGroup failed, for the reason:%v", err)
			return err
		}

		// 向 groups user 里面插入一条数据
		if err = tx.Create(&model.GroupUser{
			UserID:  g.OwnerID,
			GroupID: id,
			JoinAt:  time.Now(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	res.Group = &proto.Group{
		OwnerID:   g.OwnerID,
		Name:      g.Name,
		Password:  g.Password,
		CreatedAt: g.CreatedAt,
		Status:    g.Status,
		Id:        id,
	}

	return err
}

func (s *Service) JoinGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.GroupResp) error {
	var (
		err   error
		group = model.Group{}
	)
	logger.Info("JoinGroup user_id:%v group_id:%v", req.UserID, req.GroupID)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err       error
			groupUser model.GroupUser
		)

		// 防止 group 被删掉 防止组 id 不存在
		if err = tx.Raw(" SELECT * FROM `groups` WHERE id = ? LIMIT 1 LOCK IN SHARE MODE;", req.GroupID).Scan(&group).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.ErrGroupNotExist
			} else {
				return err
			}
		}
		// 防止重复添加
		err = tx.Where("group_id = ? and user_id = ?", req.GroupID, req.UserID).Take(&groupUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// 用户已经加入了该组
		if err == nil {
			return common.ErrUserJoinedGroup
		}

		// 没有错误，加入该组
		// 向 groups user 里面插入一条数据
		if err = tx.Create(&model.GroupUser{
			UserID:  req.UserID,
			GroupID: req.GroupID,
			JoinAt:  time.Now(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Errorf("JoinGroup failed, for the reason:%v", err)
		if err == common.ErrUserJoinedGroup {
			res.Err = &proto.Error{
				Code:    common.DBConflictCode,
				Message: err.Error(),
			}
			return nil
		} else if err == common.ErrGroupNotExist {
			res.Err = &proto.Error{
				Code:    common.DBNotFoundCode,
				Message: err.Error(),
			}
			return nil
		}
		return err
	}

	res.Group = &proto.Group{
		OwnerID:   group.OwnerID,
		Name:      group.Name,
		Password:  group.Password,
		CreatedAt: group.CreatedAt.Unix(),
		Status:    int32(group.Status),
		Id:        req.GroupID,
	}

	return nil
}
func (s *Service) LeaveGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.GroupResp) error {
	var (
		err   error
		group = model.Group{}
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err       error
			groupUser model.GroupUser
		)

		// 防止 group 被删掉 防止组 id 不存在
		if err = tx.Raw(" SELECT * FROM `groups` WHERE id = ? LIMIT 1 LOCK IN SHARE MODE;", req.GroupID).Scan(&group).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.ErrGroupNotExist
			} else {
				return err
			}
		}
		// 防止重复添加
		if r := tx.Where("group_id = ? and user_id = ?", req.GroupID, req.UserID).Delete(&groupUser).RowsAffected; r == 0 {
			return common.ErrGroupNotExist
		}

		// 如果是 owner 离开 删除该组
		if group.OwnerID == req.UserID {
			// group 一定存在
			if err = tx.Delete(&group).Error; err != nil {
				return err
			}
			// 删除组内所有成员
			if err = tx.Where("group_id = ?", req.GroupID).Delete(&groupUser).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil
				}
			}

			// 删除组内所有文件
			if err := tx.Where("group_id = ?", req.GroupID).Delete(&GroupFile{}).Error; err != nil {
				// 组内可能没有文件
				if err == gorm.ErrRecordNotFound {
					return nil
				}
				return err
			}
			return err
		}

		return nil
	})

	if err != nil {
		if err == common.ErrGroupNotExist {
			res.Err = &proto.Error{
				Code:    common.DBNotFoundCode,
				Message: err.Error(),
			}
			return nil
		}
		return err
	}

	res.Group = &proto.Group{
		OwnerID:   group.OwnerID,
		Name:      group.Name,
		Password:  group.Password,
		CreatedAt: group.CreatedAt.Unix(),
		Status:    int32(group.Status),
		Id:        req.GroupID,
	}

	return nil
}
func (s *Service) DeleteGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.GroupResp) error {
	var (
		err   error
		group = model.Group{}
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			groupUser model.GroupUser
		)

		// 防止 group 被删掉 防止组 id 不存在
		if err = tx.Raw(" SELECT * FROM `groups` WHERE id = ? LIMIT 1 LOCK IN SHARE MODE;", req.GroupID).Scan(&group).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.ErrGroupNotExist
			} else {
				return err
			}
		}

		if group.OwnerID != req.UserID {
			return common.ErrNotOwner
		}

		tx.Delete(&group)

		// 删除组内所有成员
		if r := tx.Where("group_id = ?", req.GroupID).Delete(&groupUser).RowsAffected; r == 0 {
			return common.ErrGroupNotExist
		}

		// 删除组内所有文件
		if err = tx.Where("group_id = ?", req.GroupID).Delete(&GroupFile{}).Error; err != nil {
			// 组内可能没有文件
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		return nil
	})

	if err != nil {
		if err == common.ErrGroupNotExist {
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		} else if err == common.ErrNotOwner {
			res.Err = getProtoError(err, common.DBForbiddenCode)
			return nil
		}
		return err
	}

	res.Group = &proto.Group{
		OwnerID:   group.OwnerID,
		Name:      group.Name,
		Password:  group.Password,
		CreatedAt: group.CreatedAt.Unix(),
		Status:    int32(group.Status),
		Id:        req.GroupID,
	}

	return nil
}

func (s *Service) QueryGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.GroupResp) error {
	group := &model.Group{}
	if err := db.Where("id = ?", req.GroupID).Take(group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			res.Err = &proto.Error{
				Code:    common.DBNotFoundCode,
				Message: err.Error(),
			}
			return nil
		}
	}

	res.Group = &proto.Group{
		Id:        group.ID,
		OwnerID:   group.OwnerID,
		Name:      group.Name,
		Password:  group.Password,
		CreatedAt: group.CreatedAt.Unix(),
		Status:    int32(group.Status),
	}
	return nil
}
