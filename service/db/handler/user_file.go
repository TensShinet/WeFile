package handler

import (
	"context"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
	"path/filepath"
	"time"
)

type UserFile struct {
	model.UserFile
	Size int64
}

// user_file 相关服务

// 列出 用户 所在目录的所有文件
func (s *Service) ListUserFile(ctx context.Context, req *proto.ListUserFileMetaReq, res *proto.ListUserFileMetaResp) (err error) {
	logger.Infof("ListUserFile directory:%v userID:%v", req.Directory, req.UserID)
	if req.Directory == "" {
		req.Directory = "/"
	}
	var files []*UserFile

	// slice 不报 ErrRecordNotFound
	if err = db.Raw("SELECT user_id, directory, file_name, file_id, is_directory, upload_at, last_update_at, user_files.status, files.size "+
		"FROM user_files LEFT JOIN files "+
		"ON user_files.file_id = files.id "+
		"WHERE user_id = ? and directory = ?", req.UserID, req.Directory).Scan(&files).Error; err != nil {
		logger.Errorf("ListUserFile failed, for the reason:%v", err)
		return err
	}

	res.UserFileMetaList = make([]*proto.ListFileMeta, len(files))
	for i, f := range files {
		res.UserFileMetaList[i] = &proto.ListFileMeta{
			FileID:       f.FileID,
			FileName:     f.FileName,
			IsDirectory:  f.IsDirectory,
			UploadAt:     f.UploadAt.Unix(),
			Directory:    f.Directory,
			LastUpdateAt: f.LastUpdateAt.Unix(),
			Status:       int32(f.Status),
			Size:         files[i].Size,
		}
		logger.Debugf("files[%v].Size:%v", i, files[i].Size)
	}

	return
}

// 向 user_files 表中插入一条数据
//
// 得到预插入的 id
// 如果不是目录 向 files 表中创建或获取一条数据
// 如果不是目录 更新 file 引用计数
// 向 user_files 表中插入一条数据 如果冲突 返回 ErrConflict 事务回退
func (s *Service) InsertUserFile(ctx context.Context, req *proto.InsertUserFileMetaReq, res *proto.InsertUserFileMetaResp) error {

	var (
		err                error
		fileID, userFileID int64
	)

	if fileID, err = getID(ctx); err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	}

	if userFileID, err = getID(ctx); err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	}
	logger.Infof("InsertUserFile userID:%v directory:%v filename:%v", req.UserID, req.UserFileMeta.Directory, req.UserFileMeta.FileName)
	err = db.Transaction(func(tx *gorm.DB) error {
		// file 表
		var (
			err error
		)
		fileMeta := req.FileMeta
		userFileMeta := req.UserFileMeta
		file := model.File{}

		// 检查父目录存不存在 根目录除外
		if req.UserFileMeta.Directory != "/" {
			parentHash := getUserFileHash(req.UserID, req.UserFileMeta.Directory, "")
			if err = tx.Where("hash = ?", parentHash).First(&model.UserFile{}).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return common.ErrParentDirNotFound
				}
				return err
			}
		}

		if !userFileMeta.IsDirectory {
			if err = tx.Set("gorm:query_option", "FOR UPDATE").Where(model.File{
				Hash: fileMeta.Hash,
			}).Attrs(model.File{
				ID:            fileID,
				Hash:          fileMeta.Hash,
				HashAlgorithm: fileMeta.HashAlgorithm,
				SamplingHash:  fileMeta.SamplingHash,
				Size:          fileMeta.Size,
				Location:      fileMeta.Location,
				CreateAt:      time.Unix(fileMeta.CreateAt, 0),
				UpdateAt:      time.Unix(fileMeta.CreateAt, 0),
				Status:        int(fileMeta.Status),
				Count:         0,
			}).FirstOrCreate(&file).Error; err != nil {
				// 回退
				return err
			}
			file.Count += 1
			// 更新引用计数
			if err := tx.Save(&file).Error; err != nil {
				return err
			}

		}

		// user_files 表
		// 同一个目录不能重名 如果有一样的就回退
		hash := getUserFileHash(req.UserID, req.UserFileMeta.Directory, req.UserFileMeta.FileName)
		err = tx.Where("hash = ?", hash).First(&model.UserFile{}).Error
		if err == gorm.ErrRecordNotFound {
			if err = tx.Create(&model.UserFile{
				ID:           userFileID,
				UserID:       req.UserID,
				FileID:       file.ID,
				FileName:     userFileMeta.FileName,
				IsDirectory:  userFileMeta.IsDirectory,
				Directory:    userFileMeta.Directory,
				UploadAt:     time.Unix(userFileMeta.UploadAt, 0),
				LastUpdateAt: time.Unix(userFileMeta.LastUpdateAt, 0),
				Hash:         hash, // 保证唯一性
				Status:       int(userFileMeta.Status),
			}).Error; err != nil {
				// 并发冲突
				return common.ErrConflict
			}
		} else if err != nil {
			return err
		} else {
			// 找到相同的文件
			return common.ErrConflict
		}
		// 成功插入
		res.FileMeta = req.UserFileMeta
		res.FileMeta.FileID = file.ID
		if !userFileMeta.IsDirectory {
			res.FileMeta.Size = fileMeta.Size
		}
		return nil
	})

	if err != nil {
		logger.Errorf("InsertUserFile failed, for the reason:%v", err)
		if err == common.ErrConflict {
			res.Err = getProtoError(err, common.DBConflictCode)
			return nil
		} else if err == common.ErrParentDirNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		}
		return err
	}
	return nil
}

// 删除 user_files 表中的记录
//
// 如果是目录除了删除自己，还删除以他为前缀
func (s *Service) DeleteUserFile(ctx context.Context, req *proto.DeleteUserFileReq, res *proto.DeleteUserFileResp) error {

	var (
		err error
	)

	userFileMeta := &model.UserFile{}
	logger.Infof("DeleteUserFile user_id=%v directory=%v filename:%v", req.UserID, req.Directory, req.FileName)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		hash := getUserFileHash(req.UserID, req.Directory, req.FileName)
		if err = tx.Where("hash = ?", hash).First(userFileMeta).Error; err != nil {
			return err
		}
		if userFileMeta.IsDirectory {
			if err = tx.Debug().Where("user_id = ? and directory LIKE ?", userFileMeta.UserID, filepath.Join(userFileMeta.Directory, userFileMeta.FileName)+"%").Delete(&UserFile{}).Error; err != nil {
				return err
			}
		}
		// 删除自己 触发删除钩子
		if err = tx.Delete(&userFileMeta).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("DeleteUserFile failed, for the reason:%v", err)
		return err
	}

	if err == gorm.ErrRecordNotFound {
		res.Err = getProtoError(err, common.DBNotFoundCode)
		return nil
	}

	// 删除成功
	res.FileMeta = &proto.ListFileMeta{
		FileID:       userFileMeta.FileID,
		FileName:     userFileMeta.FileName,
		IsDirectory:  userFileMeta.IsDirectory,
		UploadAt:     userFileMeta.UploadAt.Unix(),
		Directory:    userFileMeta.Directory,
		LastUpdateAt: userFileMeta.LastUpdateAt.Unix(),
		Status:       int32(userFileMeta.Status),
	}
	return nil
}

func (s *Service) QueryUserFile(ctx context.Context, req *proto.QueryUserFileReq, res *proto.QueryUserFileResp) error {

	var (
		err error
	)
	userFile := model.UserFile{}
	logger.Infof("QueryUserFile userID=%v directory=%v filename:%v", req.UserID, req.Directory, req.FileName)
	err = db.Where("hash = ?", getUserFileHash(req.UserID, req.Directory, req.FileName)).First(&userFile).Error

	if err == gorm.ErrRecordNotFound {
		res.Err = getProtoError(err, common.DBNotFoundCode)
		return nil
	} else if err != nil {
		logger.Errorf("QueryUserFile failed, for the reason:%v", err)
		return err
	}

	res.FileMeta = &proto.ListFileMeta{
		FileID:       userFile.FileID,
		FileName:     userFile.FileName,
		IsDirectory:  userFile.IsDirectory,
		UploadAt:     userFile.UploadAt.Unix(),
		Directory:    userFile.Directory,
		LastUpdateAt: userFile.LastUpdateAt.Unix(),
		Status:       int32(userFile.Status),
	}
	return nil
}
