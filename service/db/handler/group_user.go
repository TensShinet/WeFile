package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
)

type GroupUserInfo struct {
	model.GroupUser
	Email string
	Name  string
}

func (s *Service) ListGroupUser(ctx context.Context, req *proto.UserIDGroupID, res *proto.ListGroupUserResp) error {
	var (
		users []*GroupUserInfo
		err   error
	)

	if err = db.Debug().Raw("SELECT users.email, users.name, users.id as user_id, t.join_at, t.group_id "+
		"FROM users INNER JOIN (SELECT * FROM group_users WHERE group_id = ?) AS t "+
		"ON users.id = t.user_id", req.GroupID).Scan(&users).Error; err != nil {
		return err
	}

	protoUsers := make([]*proto.GroupUserInfo, len(users))
	for i, v := range users {
		protoUsers[i] = &proto.GroupUserInfo{
			Email:   v.Email,
			Name:    v.Name,
			JoinAt:  v.JoinAt.Unix(),
			UserID:  v.UserID,
			GroupID: v.GroupID,
		}
	}
	res.Users = protoUsers
	return nil
}

func (s *Service) ListUserGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.ListUserGroupResp) error {
	var (
		err    error
		groups []*model.Group
	)

	if err = db.Raw("SELECT groups.* "+
		"FROM groups INNER JOIN (SELECT * FROM group_users WHERE user_id = ?) AS t "+
		"ON groups.id = t.group_id", req.UserID).Scan(&groups).Error; err != nil {
		return err
	}

	protoGroups := make([]*proto.Group, len(groups))
	for i, v := range groups {
		protoGroups[i] = &proto.Group{
			OwnerID:   v.OwnerID,
			Name:      v.Name,
			Password:  v.Password,
			CreatedAt: v.CreatedAt.Unix(),
			Status:    int32(v.Status),
			Id:        v.ID,
		}
	}
	res.Groups = protoGroups

	return nil
}

func (s *Service) CheckUserInGroup(ctx context.Context, req *proto.UserIDGroupID, res *proto.CheckUserInGroupResp) error {
	var (
		user = &GroupUserInfo{}
		err  error
	)

	if err = db.Debug().Raw("SELECT users.email, users.name, users.id, t.join_at, t.group_id "+
		"FROM users INNER JOIN (SELECT * FROM group_users WHERE group_id = ?) AS t "+
		"ON users.id = t.user_id", req.GroupID).Scan(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		}
		return err
	}

	res.GroupUserInfo = &proto.GroupUserInfo{

		Email:   user.Email,
		Name:    user.Name,
		JoinAt:  user.JoinAt.Unix(),
		UserID:  user.UserID,
		GroupID: user.GroupID,
	}
	return nil
}
