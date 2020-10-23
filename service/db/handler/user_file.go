package handler

import (
	"context"
	"fmt"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/TensShinet/WeFile/service/db/model"
	"github.com/TensShinet/WeFile/service/db/proto"
	"gorm.io/gorm"
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
			Size:         files[i].Size,
		}
		logger.Debugf("files[%v].Size:%v", i, files[i].Size)
	}

	return
}

var errConflict = fmt.Errorf("confilct error")

// 向 user_files 表中插入一条数据
//
// 得到预插入的 id
// 向 files 表中创建或获取一条数据
// 更新 file 引用计数
// 向 user_files 表中插入一条数据 如果冲突 返回 errConflict 事务回退
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
		meta := req.FileMeta
		file := model.File{}
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where(model.File{
			Hash: meta.Hash,
		}).Attrs(model.File{
			ID:            fileID,
			Hash:          meta.Hash,
			HashAlgorithm: meta.HashAlgorithm,
			SamplingHash:  meta.SamplingHash,
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
		// 更新引用计数
		if err := tx.Save(&file).Error; err != nil {
			return err
		}

		// user_files 表
		// 同一个目录不能重名
		userFile := req.UserFileMeta
		// 如果有一样的就回退
		hash := getUserFileHash(req.UserID, req.UserFileMeta.Directory, req.UserFileMeta.FileName)
		err = tx.Where("hash = ?", hash).First(&model.UserFile{}).Error
		if err == gorm.ErrRecordNotFound {
			if err = tx.Create(&model.UserFile{
				ID:           userFileID,
				UserID:       req.UserID,
				FileID:       file.ID,
				FileName:     userFile.FileName,
				IsDirectory:  userFile.IsDirectory,
				Directory:    userFile.Directory,
				UploadAt:     time.Unix(userFile.UploadAt, 0),
				LastUpdateAt: time.Unix(userFile.LastUpdateAt, 0),
				Hash:         hash, // 保证唯一性
				Status:       int(userFile.Status),
			}).Error; err != nil {
				// 并发冲突
				return errConflict
			}
		} else if err != nil {
			return err
		} else {
			// 找到相同的文件
			return errConflict
		}
		// 成功插入
		res.FileMeta = req.UserFileMeta
		res.FileMeta.FileID = file.ID
		return nil
	})

	if err != nil {
		logger.Errorf("InsertUserFile failed, for the reason:%v", err)
		if err == errConflict {
			res.Err = getProtoError(err, common.DBConflictCode)
			return nil
		}
		return err
	}
	return nil
}

// 删除 user_files 表中的一条记录
//
// 删除 user_files 中的一条记录如果没有 返回 not found 错误
// 更新 files 表中的引用计数
func (s *Service) DeleteUserFile(ctx context.Context, req *proto.DeleteUserFileReq, res *proto.DeleteUserFileResp) error {

	var (
		err          error
		affectedRows int64
	)

	userFileMeta := &model.UserFile{}
	logger.Infof("DeleteUserFile userID=%v directory=%v filename:%v", req.UserID, req.Directory, req.FileName)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)

		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ? AND directory = ? AND file_name = ?", req.UserID, req.Directory, req.FileName).First(&userFileMeta).Error; err != nil {
			return err
		}

		affectedRows = db.Delete(userFileMeta).RowsAffected

		file := model.File{}
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userFileMeta.FileID).First(&file).Error; err != nil {
			return err
		}
		// file 引用计数 - 1
		file.Count -= 1
		err = tx.Save(&file).Error
		// 成功删除
		return err

	})

	if err != nil {
		logger.Errorf("DeleteUserFile failed, for the reason:%v", err)
		return err
	}

	if affectedRows == 0 {
		res.Err = getProtoError(gorm.ErrRecordNotFound, common.DBNotFoundCode)
		return nil
	}

	// 删除成功
	res.FileMeta = &proto.UserFileMeta{
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
	err = db.Where("user_id=? AND directory = ? AND file_name = ?", req.UserID, req.Directory, req.FileName).First(&userFile).Error

	if err == gorm.ErrRecordNotFound {
		res.Err = getProtoError(err, common.DBNotFoundCode)
		return nil
	} else if err != nil {
		logger.Errorf("QueryUserFile failed, for the reason:%v", err)
		return err
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
	return nil
}
