package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/id_generator/conf"
	proto "github.com/TensShinet/WeFile/service/id_generator/proto"
	"github.com/bwmarrin/snowflake"
)

type Service struct{}

func (s *Service) GenerateID(ctx context.Context, req *proto.IDReq, res *proto.IDResp) error {
	config := conf.GetConfig()
	node, err := snowflake.NewNode(int64(config.Service.NodeID))
	if err != nil {
		res.Id = -1
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
		return err
	}
	res.Id = node.Generate().Int64()
	return nil
}
