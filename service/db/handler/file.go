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

func (s *Service) QueryFileMeta(ctx context.Context, req *proto.QueryFileMetaReq, res *proto.QueryFileMetaResp) error {

	var (
		err error
	)
	logger.Infof("QueryFileMeta id=%v hash=%v", req.Id, req.Hash)
	file := model.File{}
	err = db.Where("id = ? || hash = ?", req.Id, req.Hash).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Infof("QueryFileMeta id=%v hash=%v not found", req.Id, req.Hash)
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		} else {
			logger.Errorf("QueryFileMeta id=%v hash=%v failed, for the reason:%v", req.Id, req.Hash, err)
			return err
		}
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
	return nil
}
