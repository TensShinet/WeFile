package main

import (
	"context"
	"github.com/TensShinet/WeFile/service/id_generator/handler"
	proto "github.com/TensShinet/WeFile/service/id_generator/proto"
	"testing"
	"time"
)

// 测试并发 id 生成
func TestGenerate(t *testing.T) {
	service := handler.Service{}
	var last int64 = 0
	for i := 0; i < 100000; i++ {
		temp := proto.IDResp{}
		err := service.GenerateID(context.TODO(), nil, &temp)
		if err != nil {
			t.Error("test failed, for the reason:" + err.Error())
			return
		}
		if temp.Id < last {
			t.Error("test failed, id doesn't not increase automatically")
			break
		}
		last = temp.Id
		time.Sleep(time.Microsecond)
	}
}
