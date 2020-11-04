package main

import (
	"context"
	"github.com/TensShinet/WeFile/service/db/conf"
	"github.com/TensShinet/WeFile/service/db/conn"
	"github.com/TensShinet/WeFile/service/db/handler"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"testing"
	"time"
)

var (
	s       handler.Service
	user1ID int64
	user2ID int64
	user3ID int64

	group1ID int64
	group2ID int64
)

func TestInit(t *testing.T) {
	conf.Init("db_conf.yml")
	config := conf.GetConfig()
	config.DB.MySQL.DSN = "tenshine:tenshine@tcp(127.0.0.1:3306)/wefile_test?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化
	conn.Init()
	model.Init()
	handler.Init()

}

func TestCreateUser(t *testing.T) {
	// 插入 User
	res1 := &proto.InsertUserResp{}
	if err := s.InsertUser(context.TODO(), &proto.InsertUserReq{User: &proto.User{
		RoleID:         1,
		Name:           "ts1",
		Password:       "ts6666",
		Email:          "1@gmail.com",
		Phone:          "1888888",
		EmailValidated: false,
		PhoneValidated: false,
		SignUpAt:       1,
		LasActiveAt:    2,
		Profile:        "ts666",
		Status:         0,
	}}, res1); err != nil {
		t.Fatal(err)
	}
	t.Logf("user1 id: %v", res1.Id)
	user1ID = res1.Id

	if err := s.InsertUser(context.TODO(), &proto.InsertUserReq{User: &proto.User{
		RoleID:         1,
		Name:           "ts2",
		Password:       "ts6666",
		Email:          "2@gmail.com",
		Phone:          "1888888",
		EmailValidated: false,
		PhoneValidated: false,
		SignUpAt:       1,
		LasActiveAt:    2,
		Profile:        "ts666",
		Status:         0,
	}}, res1); err != nil {
		t.Fatal(err)
	}
	t.Logf("user2 id: %v", res1.Id)
	user2ID = res1.Id

	if err := s.InsertUser(context.TODO(), &proto.InsertUserReq{User: &proto.User{
		RoleID:         1,
		Name:           "ts3",
		Password:       "ts6666",
		Email:          "3@gmail.com",
		Phone:          "1888888",
		EmailValidated: false,
		PhoneValidated: false,
		SignUpAt:       1,
		LasActiveAt:    2,
		Profile:        "ts666",
		Status:         0,
	}}, res1); err != nil {
		t.Fatal(err)
	}
	t.Logf("user3 id: %v", res1.Id)
	user3ID = res1.Id
}

func TestCreateUserFile(t *testing.T) {
	// 插入 user_file
	res2 := &proto.InsertUserFileMetaResp{}
	if err := s.InsertUserFile(context.TODO(), &proto.InsertUserFileMetaReq{
		UserFileMeta: &proto.ListFileMeta{
			FileName:     "ts_file1",
			IsDirectory:  false,
			UploadAt:     1,
			Directory:    "/",
			LastUpdateAt: 2,
			Status:       0,
		},
		FileMeta: &proto.FileMeta{
			Hash:          "hash1",
			SamplingHash:  "samplingHash1",
			HashAlgorithm: "SHA256",
			Size:          128,
			Location:      "local",
			CreateAt:      1,
			Status:        2,
		},
		UserID: user1ID,
	}, res2); err != nil {
		t.Fatal(err)
	}

	// 插入 user_file
	if err := s.InsertUserFile(context.TODO(), &proto.InsertUserFileMetaReq{
		UserFileMeta: &proto.ListFileMeta{
			FileName:     "ts_file2",
			IsDirectory:  false,
			UploadAt:     1,
			Directory:    "/",
			LastUpdateAt: 2,
			Status:       0,
		},
		FileMeta: &proto.FileMeta{
			Hash:          "hash2",
			SamplingHash:  "samplingHash2",
			HashAlgorithm: "SHA256",
			Size:          256,
			Location:      "local",
			CreateAt:      1,
			Status:        2,
		},
		UserID: user1ID,
	}, res2); err != nil {
		t.Fatal(err)
	}
}

func TestUserFile(t *testing.T) {
	s := handler.Service{}

	// 插入 User
	res1 := &proto.InsertUserResp{}
	if err := s.InsertUser(context.TODO(), &proto.InsertUserReq{User: &proto.User{
		RoleID:         1,
		Name:           "ts",
		Password:       "ts6666",
		Email:          "tanshunwork@gmail.com",
		Phone:          "1888888",
		EmailValidated: false,
		PhoneValidated: false,
		SignUpAt:       1,
		LasActiveAt:    2,
		Profile:        "ts666",
		Status:         0,
	}}, res1); err != nil {
		t.Fatal(err)
	}
	t.Logf("user id: %v", res1.Id)

	// 插入 user_file
	res2 := &proto.InsertUserFileMetaResp{}
	if err := s.InsertUserFile(context.TODO(), &proto.InsertUserFileMetaReq{
		UserFileMeta: &proto.ListFileMeta{
			FileName:     "ts_file",
			IsDirectory:  false,
			UploadAt:     1,
			Directory:    "/",
			LastUpdateAt: 2,
			Status:       0,
		},
		FileMeta: &proto.FileMeta{
			Hash:          "hash1234",
			SamplingHash:  "samplingHash1234",
			HashAlgorithm: "SHA256",
			Size:          128,
			Location:      "local",
			CreateAt:      1,
			Status:        2,
		},
		UserID: res1.Id,
	}, res2); err != nil {
		t.Fatal(err)
	}
	m := res2.FileMeta
	t.Logf("file id: %v", m.FileID)
	if m.IsDirectory || m.FileName != "ts_file" || m.UploadAt != 1 || m.Directory != "/" ||
		m.LastUpdateAt != 2 || m.Status != 0 {
		t.Log(m.IsDirectory, m.FileName != "ts_file", m.UploadAt != 1, m.Directory != "/",
			m.LastUpdateAt != 2, m.Status)
		t.Fatal("insert user file failed")
	}

	// 删除 user file
	res3 := &proto.DeleteUserFileResp{}
	if err := s.DeleteUserFile(context.TODO(), &proto.DeleteUserFileReq{
		UserID:    res1.Id,
		Directory: "/",
		FileName:  "ts_file",
	}, res3); err != nil {
		t.Fatal(err)
	}
	m = res3.FileMeta
	if m.IsDirectory || m.FileName != "ts_file" || m.UploadAt != 1 || m.Directory != "/" ||
		m.LastUpdateAt != 2 || m.Status != 0 {
		t.Log(m.IsDirectory, m.FileName != "ts_file", m.UploadAt != 1, m.Directory != "/",
			m.LastUpdateAt != 2, m.Status)
		t.Fatal("insert user file failed")
	}

	// 删除 User
	res4 := &proto.DeleteUserResp{User: &proto.User{}}
	if err := s.DeleteUser(context.TODO(), &proto.DeleteUserReq{Id: res1.Id}, res4); err != nil {
		t.Fatal(err)
	}
	u := res4.User
	if u.Id != res1.Id || u.RoleID != 1 || u.Name != "ts" || u.Password != "ts6666" || u.Email != "tanshunwork@gmail.com" ||
		u.SignUpAt != 1 || u.LasActiveAt != 2 || u.Profile != "ts666" || u.Status != 0 ||
		u.EmailValidated || u.PhoneValidated {
		t.Log(u.Id != res1.Id, u.RoleID != 1, u.Name != "ts", u.Password != "ts6666", u.Email != "tanshunwork@gmail.com",
			u.SignUpAt != 1, u.LasActiveAt != 2, u.Profile != "ts666", u.Status != 0, u.EmailValidated, u.PhoneValidated)
		t.Fatalf("delete failed")
	}
}

func TestCreateGroup(t *testing.T) {
	res2 := &proto.CreateGroupResp{}

	// 创建组
	if err := s.CreateGroup(context.TODO(), &proto.CreateGroupReq{
		Group: &proto.Group{
			OwnerID:   user1ID,
			Name:      "group1",
			Password:  "123",
			CreatedAt: time.Now().Unix(),
			Status:    0,
		},
	}, res2); err != nil {
		t.Fatal(err)
	}

	t.Logf("group1 id: %v", res2.Group.Id)
	group1ID = res2.Group.Id

	// 创建组
	if err := s.CreateGroup(context.TODO(), &proto.CreateGroupReq{
		Group: &proto.Group{
			OwnerID:   user1ID,
			Name:      "group2",
			Password:  "123",
			CreatedAt: time.Now().Unix(),
			Status:    0,
		},
	}, res2); err != nil {
		t.Fatal(err)
	}

	t.Logf("group2 id: %v", res2.Group.Id)
	group2ID = res2.Group.Id

}

func TestJoinGroup(t *testing.T) {
	// 加入组
	res5 := &proto.GroupResp{}
	if err := s.JoinGroup(context.TODO(), &proto.UserIDGroupID{
		UserID:  user2ID,
		GroupID: group1ID,
	}, res5); err != nil {
		t.Fatal(err)
	}
	t.Logf("join group id: %v", res5.Group.Id)

}

func TestListGroupUser(t *testing.T) {

	// 列出组内所有用户
	res7 := &proto.ListGroupUserResp{}
	if err := s.ListGroupUser(context.TODO(), &proto.UserIDGroupID{
		GroupID: group1ID,
	}, res7); err != nil {
		t.Fatal(err)
	}

	t.Logf("user list length:%v", len(res7.Users))
	for _, v := range res7.Users {
		t.Logf("groupID:%v name:%v, email:%v userID:%v, joinAt:%v", v.GroupID, v.Name, v.Email, v.UserID, v.JoinAt)
	}

}

func TestListUserGroup(t *testing.T) {
	// 列出一个用户的所有组
	res := &proto.ListUserGroupResp{}
	if err := s.ListUserGroup(context.TODO(), &proto.UserIDGroupID{
		UserID: user1ID,
	}, res); err != nil {
		t.Fatal(err)
	}

	t.Logf("group list length:%v", len(res.Groups))
	for _, v := range res.Groups {
		t.Logf("groupID:%v name:%v, ownerID:%v, password:%v", v.Id, v.Name, v.OwnerID, v.Password)
	}
}

func TestCreateGroupFile(t *testing.T) {
	// 插入 group_file1
	res2 := &proto.InsertGroupFileResp{}
	if err := s.InsertGroupFile(context.TODO(), &proto.InsertGroupFileReq{
		GroupFileMeta: &proto.ListFileMeta{
			FileName:     "group_file1",
			IsDirectory:  false,
			UploadAt:     1,
			Directory:    "/",
			LastUpdateAt: 2,
			Status:       0,
		},
		FileMeta: &proto.FileMeta{
			Hash: "hash1",
		},
		GroupID: group1ID,
	}, res2); err != nil {
		t.Fatal(err)
	}

	// 插入 group_file2
	if err := s.InsertGroupFile(context.TODO(), &proto.InsertGroupFileReq{
		GroupFileMeta: &proto.ListFileMeta{
			FileName:     "group_file2",
			IsDirectory:  false,
			UploadAt:     1,
			Directory:    "/",
			LastUpdateAt: 2,
			Status:       0,
		},
		FileMeta: &proto.FileMeta{
			Hash: "hash2",
		},
		GroupID: group1ID,
	}, res2); err != nil {
		t.Fatal(err)
	}
}

func TestListGroupFile(t *testing.T) {
	res := &proto.ListGroupFileResp{}
	if err := s.ListGroupFile(context.TODO(), &proto.ListGroupFileReq{
		GroupID:   group1ID,
		Directory: "/",
	}, res); err != nil {
		t.Fatal(err)
	}

	t.Logf("group file len:%v", len(res.GroupFileMetaList))
	if len(res.GroupFileMetaList) != 2 {
		t.Fatal("group file len error")
	}

	for _, v := range res.GroupFileMetaList {
		t.Logf("fileName:%v fileID:%v", v.FileName, v.FileID)
	}

}

func TestDeleteGroupFile(t *testing.T) {
	res := &proto.DeleteGroupFileResp{}
	if err := s.DeleteGroupFile(context.TODO(), &proto.DeleteGroupFileReq{
		GroupID:   group1ID,
		Directory: "/",
		FileName:  "group_file1",
	}, res); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteGroupFile(context.TODO(), &proto.DeleteGroupFileReq{
		GroupID:   group1ID,
		Directory: "/",
		FileName:  "group_file2",
	}, res); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteGroup(t *testing.T) {
	// 尝试删除不是自己的组
	res6 := &proto.GroupResp{}
	if err := s.DeleteGroup(context.TODO(), &proto.UserIDGroupID{
		UserID:  user3ID,
		GroupID: group1ID,
	}, res6); err != nil {
		t.Fatal(err)
	}

	t.Log("删除不是自己的组: ", res6.Err.Message)
	if err := s.DeleteGroup(context.TODO(), &proto.UserIDGroupID{
		UserID:  user1ID,
		GroupID: group1ID,
	}, res6); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteGroup(context.TODO(), &proto.UserIDGroupID{
		UserID:  user1ID,
		GroupID: group2ID,
	}, res6); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteUser(t *testing.T) {
	res2 := &proto.DeleteUserResp{User: &proto.User{}}
	if err := s.DeleteUser(context.TODO(), &proto.DeleteUserReq{Id: user1ID}, res2); err != nil {
		t.Fatal(err)
	}
	u := res2.User
	if u.Id != user1ID || u.RoleID != 1 || u.Name != "ts1" || u.Password != "ts6666" || u.Email != "1@gmail.com" ||
		u.SignUpAt != 1 || u.LasActiveAt != 2 || u.Profile != "ts666" || u.Status != 0 ||
		u.EmailValidated || u.PhoneValidated {
		t.Log(u.Id != user1ID, u.RoleID != 1, u.Name != "ts1", u.Password != "ts6666", u.Email != "1@gmail.com",
			u.SignUpAt != 1, u.LasActiveAt != 2, u.Profile != "ts666", u.Status != 0, u.EmailValidated, u.PhoneValidated)
		t.Fatalf("delete failed")
	}

	if err := s.DeleteUser(context.TODO(), &proto.DeleteUserReq{Id: user2ID}, res2); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteUser(context.TODO(), &proto.DeleteUserReq{Id: user3ID}, res2); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteUserFile(t *testing.T) {
	res3 := &proto.DeleteUserFileResp{}
	if err := s.DeleteUserFile(context.TODO(), &proto.DeleteUserFileReq{
		UserID:    user1ID,
		Directory: "/",
		FileName:  "ts_file1",
	}, res3); err != nil {
		t.Fatal(err)
	}
	m := res3.FileMeta
	if m.IsDirectory || m.FileName != "ts_file1" || m.UploadAt != 1 || m.Directory != "/" ||
		m.LastUpdateAt != 2 || m.Status != 0 {
		t.Log(m.IsDirectory, m.FileName != "ts_file1", m.UploadAt != 1, m.Directory != "/",
			m.LastUpdateAt != 2, m.Status)
		t.Fatal("insert user file failed")
	}

	if err := s.DeleteUserFile(context.TODO(), &proto.DeleteUserFileReq{
		UserID:    user2ID,
		Directory: "/",
		FileName:  "ts_file2",
	}, res3); err != nil {
		t.Fatal(err)
	}

}
