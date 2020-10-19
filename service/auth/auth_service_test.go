package main

import (
	"context"
	"github.com/TensShinet/WeFile/service/auth/conf"
	"github.com/TensShinet/WeFile/service/auth/handler"
	"github.com/TensShinet/WeFile/service/auth/proto"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	conf.Init("auth_conf.yml")
	s := handler.Service{}
	res := &proto.EncodeResp{}
	if err := s.DownloadJWTEncode(context.TODO(), &proto.DownloadFileMeta{
		FileID:   123456,
		FileName: "ts666",
	}, res); err != nil {
		t.Fatal(err)
	}

	t.Log("download token ", res.Token)
	res1 := &proto.DownloadJWTDecodeResp{}
	if err := s.DownloadJWTDecode(context.TODO(), &proto.DecodeReq{Token: res.Token}, res1); err != nil {
		t.Fatal(err)
	}
	if res1.FileMeta.FileID != 123456 || res1.FileMeta.FileName != "ts666" {
		t.Fatal("Decode failed ", res1.FileMeta.FileID, res1.FileMeta.FileName)
	}
}

func TestUpload(t *testing.T) {

	s := handler.Service{}
	res := &proto.EncodeResp{}
	if err := s.UploadJWTEncode(context.TODO(), &proto.UploadFileMeta{
		UserID:    123456,
		Directory: "/dir1/dir2",
		FileName:  "ts666",
	}, res); err != nil {
		t.Fatal(err)
	}

	t.Log("upload token", res.Token)

	res1 := &proto.UploadJWTDecodeResp{}
	if err := s.UploadJWTDecode(context.TODO(), &proto.DecodeReq{Token: res.Token}, res1); err != nil {
		t.Fatal(err)
	}
	if res1.FileMeta.UserID != 123456 || res1.FileMeta.FileName != "ts666" || res1.FileMeta.Directory != "/dir1/dir2" {
		t.Fatal("Decode failed ", res1.FileMeta.UserID, res1.FileMeta.FileName)
	}
}

func TestDownloadAuth(t *testing.T) {
	// 设置超时时间 1s
	config := conf.GetConfig()
	config.JWT.ValidTime = 1
	s := handler.Service{}
	res := &proto.EncodeResp{}
	if err := s.DownloadJWTEncode(context.TODO(), &proto.DownloadFileMeta{
		FileID:   123456,
		FileName: "ts666",
	}, res); err != nil {
		t.Fatal(err)
	}

	t.Log("download token ", res.Token)
	time.Sleep(time.Second * 2)
	res1 := &proto.DownloadJWTDecodeResp{}
	if err := s.DownloadJWTDecode(context.TODO(), &proto.DecodeReq{Token: res.Token}, res1); err != nil {
		t.Logf("token is expired:%v", err.Error())
		return
	}
	if res1.FileMeta.FileID != 123456 || res1.FileMeta.FileName != "ts666" {
		t.Fatal("Decode failed ", res1.FileMeta.FileID, res1.FileMeta.FileName)
	}
}
