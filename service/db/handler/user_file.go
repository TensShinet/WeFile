package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
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

// 向 files 中插入一条数据
// 得到这条 file 数据
//
func (s *Service) InsertUserFile(ctx context.Context, req *proto.InsertUserFileMetaReq, res *proto.InsertUserFileMetaResp) (err error) {
	fileID, err := getID(ctx)
	defer func(err error) {
		if err != nil {
			res.Err = &proto.Error{
				Code:    -1,
				Message: err.Error(),
			}
			return
		}
	}(err)
	// TODO: 检查并发问题
	err = db.Transaction(func(tx *gorm.DB) error {
		// file 表
		var (
			err error
			id  int64
		)
		meta := req.FileMeta
		file := model.File{}
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where(model.File{
			Hash: meta.Hash,
		}).Attrs(model.File{
			ID:            fileID,
			Hash:          meta.Hash,
			HashAlgorithm: meta.HashAlgorithm,
			Size:          meta.Size,
			Location:      meta.Location,
			CreateAt:      time.Unix(meta.CreateAt, 0),
			UpdateAt:      time.Unix(meta.CreateAt, 0),
			Status:        int(meta.Status),
			Count:         0,
		}).FirstOrCreate(&file).Error; err != nil {
			// 回退
			return err
		}
		file.Count += 1

		if err := tx.Save(&file).Error; err != nil {
			return err
		}

		// user_files 表
		// 同一个目录不能重名
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

		if err == nil {
			res.FileMeta = req.UserFileMeta
			res.FileMeta.FileID = file.ID
		}
		return err
	})
	return err
}

// 删除 user_files 表中的一条记录
// 更新 files 表中的引用计数
func (s *Service) DeleteUserFile(ctx context.Context, req *proto.DeleteUserFileReq, res *proto.DeleteUserFileResp) (err error) {

	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		userFileMeta := &model.UserFile{}
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ? AND directory = ? AND file_name = ?", req.UserID, req.Directory, req.FileName).First(&userFileMeta).Error; err != nil {
			return err
		}

		if err = db.Delete(userFileMeta).Error; err != nil {
			return err
		}
		file := model.File{}
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userFileMeta.FileID).First(&file).Error; err != nil {
			return err
		}

		file.Count -= 1

		err = tx.Save(&file).Error
		// 删除 file 的信息
		if err == nil {
			res.FileMeta = &proto.UserFileMeta{
				FileID:       userFileMeta.FileID,
				FileName:     userFileMeta.FileName,
				IsDirectory:  userFileMeta.IsDirectory,
				UploadAt:     userFileMeta.UploadAt.Unix(),
				Directory:    userFileMeta.Directory,
				LastUpdateAt: userFileMeta.LastUpdateAt.Unix(),
				Status:       int32(userFileMeta.Status),
			}
		}

		return err

	})
	return
}

func (s *Service) QueryUserFile(ctx context.Context, req *proto.QueryUserFileReq, res *proto.QueryUserFileResp) (err error) {
	userFile := model.UserFile{}
	err = db.Where("user_id=? AND directory = ? AND file_name = ?", req.UserID, req.Directory, req.FileName).First(&userFile).Error
	if err != nil {
		res.Err = &proto.Error{Code: -1, Message: err.Error()}
		if err == gorm.ErrRecordNotFound {
			res.Err.Code = common.DBNotFoundCode
		}
		return
	}
	res.FileMeta = &proto.UserFileMeta{
		FileID:       userFile.FileID,
		FileName:     userFile.FileName,
		IsDirectory:  userFile.IsDirectory,
		UploadAt:     userFile.UploadAt.Unix(),
		Directory:    userFile.Directory,
		LastUpdateAt: userFile.LastUpdateAt.Unix(),
		Status:       int32(userFile.Status),
	}
	return
}
