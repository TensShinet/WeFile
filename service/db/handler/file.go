package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
)

// file 相关服务
func (s *Service) InsertFileMeta(ctx context.Context, req *proto.InsertFileMetaReq, res *proto.InsertFileMetaResp) error {
	return nil
}

func (s *Service) QueryFileMeta(ctx context.Context, req *proto.QueryFileMetaReq, res *proto.QueryFileMetaResp) (err error) {
	file := model.File{}
	err = db.Where("hash = ?", req.Hash).First(&file).Error
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
		if err == gorm.ErrRecordNotFound {
			res.Err.Code = common.DBNotFoundCode
		}
		return
	}
	res.FileMeta = &proto.FileMeta{
		Id:            file.ID,
		Hash:          file.Hash,
		HashAlgorithm: file.HashAlgorithm,
		Size:          file.Size,
		Location:      file.Location,
		CreateAt:      file.CreateAt.Unix(),
		Status:        int32(file.Status),
	}
	return
}
