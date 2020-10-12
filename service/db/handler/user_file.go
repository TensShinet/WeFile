package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"time"
)

// user_file 相关服务

func (s *Service) ListUserFile(ctx context.Context, req *proto.ListUserFileMetaReq, res *proto.ListUserFileMetaResp) (err error) {
	if req.Directory == "" {
		req.Directory = "/"
	}
	var files []*model.UserFile
	err = db.Where("user_id = ? and directory = ?", req.UserID, req.Directory).Find(&files).Error
	res.UserFileMetaList = make([]*proto.UserFileMeta, len(files))
	for i, f := range files {
		res.UserFileMetaList[i] = &proto.UserFileMeta{
			FileID:       f.FileID,
			FileName:     f.FileName,
			IsDirectory:  f.IsDirectory,
			UploadAt:     f.UploadAt.Unix(),
			Directory:    f.Directory,
			LastUpdateAt: f.LastUpdateAt.Unix(),
			Status:       int32(f.Status),
		}
	}
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
	}
	return
}

// 同时向 user_file 表和 file 表内插入一条数据
func (s *Service) InsertUserFile(ctx context.Context, req *proto.InsertUserFileMetaReq, res *proto.InsertUserFileMetaResp) (err error) {
	fileID, err := getID(ctx)
	defer func(err error) {
		if err != nil {
			res.Err = &proto.Error{
				Code:    -1,
				Message: err.Error(),
			}
		}
	}(err)

	err = db.Transaction(func(tx *gorm.DB) error {
		// file 表
		var (
			err error
			id  int64
		)
		meta := req.FileMeta
		if err = tx.Where(model.File{
			Hash: req.FileMeta.Hash,
		}).FirstOrCreate(&model.File{
			ID:            fileID,
			Hash:          meta.Hash,
			HashAlgorithm: meta.HashAlgorithm,
			Size:          meta.Size,
			Location:      meta.Location,
			CreateAt:      time.Unix(meta.CreateAt, 0),
			UpdateAt:      time.Unix(meta.CreateAt, 0),
			Status:        int(meta.Status),
		}).Error; err != nil {
			// 回退
			return err
		}

		// 拿到 fileID
		file := model.File{}
		tx.Where("hash = ?", meta.Hash).First(&file)

		// user_files 表
		// TODO:上传同一个文件文件名不能一样
		userFile := req.UserFileMeta
		if id, err = getID(ctx); err != nil {
			res.Err = &proto.Error{
				Code:    -1,
				Message: err.Error(),
			}
			return err
		}
		// 如果有一样的就回退
		err = tx.Create(&model.UserFile{
			ID:           id,
			UserID:       req.UserID,
			FileID:       file.ID,
			FileName:     userFile.FileName,
			IsDirectory:  userFile.IsDirectory,
			Directory:    userFile.Directory,
			UploadAt:     time.Unix(userFile.UploadAt, 0),
			LastUpdateAt: time.Unix(userFile.LastUpdateAt, 0),
			Status:       int(userFile.Status),
		}).Error
		return err
	})
	return err
}
