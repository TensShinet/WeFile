package main

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/conf"
	"github.com/TensShinet/WeFile/service/db/conn"
	"github.com/TensShinet/WeFile/service/db/handler"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"testing"
)

func TestSession(t *testing.T) {
	conf.Init("db_conf.yml")
	config := conf.GetConfig()
	config.DB.MySQL.DSN = "tenshine:tenshine@tcp(127.0.0.1:3306)/wefile_test?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化
	conn.Init()
	model.Init()
	handler.Init()
	s := handler.Service{}
	// 插入
	res1 := &proto.InsertSessionResp{}
	if err := s.InsertSession(context.TODO(), &proto.InsertSessionReq{
		Session: &proto.Session{
			UserID:    123456,
			Token:     "token_123456",
			CreatedAt: 111,
			ExpireAt:  222,
			CSRFToken: "csrf_token_123456",
		},
	}, res1); err != nil {
		t.Fatalf("insert session failed, for the reason:%v", err)
	}
	// 查询已存在
	res2 := &proto.GetUserSessionResp{}
	if err := s.GetUserSession(context.TODO(), &proto.GetUserSessionReq{
		UserID: 123456,
	}, res2); err != nil {
		t.Fatalf("get user session failed, for the reason:%v", err)
	}
	// 查询不存在
	res3 := &proto.GetUserSessionResp{}
	if err := s.GetUserSession(context.TODO(), &proto.GetUserSessionReq{
		UserID: 456789,
	}, res3); err != nil {
		if res3.Err.Code == common.DBNotFoundCode {
			t.Log("DB Not Found")
		} else {
			t.Fatalf("get user session failed, for the reason:%v", err)
		}
	}
	// 删除 已存在
	res4 := &proto.DeleteUserSessionResp{}
	if err := s.DeleteUserSession(context.TODO(), &proto.DeleteUserSessionReq{UserID: 123456}, res4); err != nil {
		t.Fatal(err)
	}

	// 删除 不存在
	res5 := &proto.DeleteUserSessionResp{}
	if err := s.DeleteUserSession(context.TODO(), &proto.DeleteUserSessionReq{UserID: 456789}, res5); err != nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
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

	// 删除 User
	res2 := &proto.DeleteUserResp{User: &proto.User{}}
	if err := s.DeleteUser(context.TODO(), &proto.DeleteUserReq{Id: res1.Id}, res2); err != nil {
		t.Fatal(err)
	}
	u := res2.User
	if u.Id != res1.Id || u.RoleID != 1 || u.Name != "ts" || u.Password != "ts6666" || u.Email != "tanshunwork@gmail.com" ||
		u.SignUpAt != 1 || u.LasActiveAt != 2 || u.Profile != "ts666" || u.Status != 0 ||
		u.EmailValidated || u.PhoneValidated {
		t.Log(u.Id != res1.Id, u.RoleID != 1, u.Name != "ts", u.Password != "ts6666", u.Email != "tanshunwork@gmail.com",
			u.SignUpAt != 1, u.LasActiveAt != 2, u.Profile != "ts666", u.Status != 0, u.EmailValidated, u.PhoneValidated)
		t.Fatalf("delete failed")
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
		UserFileMeta: &proto.UserFileMeta{
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
