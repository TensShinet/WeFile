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

type GroupFile struct {
	model.GroupFile
	Size int64
}

func (s *Service) ListGroupFile(ctx context.Context, req *proto.ListGroupFileReq, res *proto.ListGroupFileResp) error {
	var (
		err error
	)

	logger.Infof("ListGroupFile directory:%v groupID:%v", req.Directory, req.GroupID)
	if req.Directory == "" {
		req.Directory = "/"
	}
	var files []*GroupFile

	// slice 不报 ErrRecordNotFound
	if err = db.Raw("SELECT group_id, directory, file_name, file_id, is_directory, upload_at, last_update_at, group_files.status, files.size "+
		"FROM group_files LEFT JOIN files "+
		"ON group_files.file_id = files.id "+
		"WHERE group_id = ? and directory = ?", req.GroupID, req.Directory).Scan(&files).Error; err != nil {
		logger.Errorf("ListGroupFile failed, for the reason:%v", err)
		return err
	}

	res.GroupFileMetaList = make([]*proto.ListFileMeta, len(files))
	for i, f := range files {
		res.GroupFileMetaList[i] = &proto.ListFileMeta{
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

	return nil
}

func (s *Service) InsertGroupFile(ctx context.Context, req *proto.InsertGroupFileReq, res *proto.InsertGroupFileResp) error {
	var (
		err                 error
		fileID, groupFileID int64
	)

	if fileID, err = getID(ctx); err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	}

	if groupFileID, err = getID(ctx); err != nil {
		res.Err = getProtoError(err, common.DBServiceError)
		return err
	}
	logger.Infof("InsertGroupFile groupID:%v directory:%v filename:%v", req.GroupID, req.GroupFileMeta.Directory, req.GroupFileMeta.FileName)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err   error
			group model.Group
		)
		fileMeta := req.FileMeta
		groupFileMeta := req.GroupFileMeta
		file := model.File{}

		// 检查 group 是不是存在
		// 防止 group 被删掉 防止组 id 不存在
		if err = tx.Raw(" SELECT * FROM `groups` WHERE id = ? LIMIT 1 LOCK IN SHARE MODE;", req.GroupID).Scan(&group).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.ErrGroupNotExist
			} else {
				return err
			}
		}

		// 检查父目录存不存在 根目录除外
		if groupFileMeta.Directory != "/" {
			parentHash := getGroupFileHash(req.GroupID, groupFileMeta.Directory, "")
			if err = tx.Debug().Where("hash = ?", parentHash).First(&model.GroupFile{}).Error; err != nil {
				return common.ErrParentDirNotFound
			}
		}

		if !groupFileMeta.IsDirectory {
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

		// group_files 表
		// 同一个目录不能重名 如果有一样的就回退
		hash := getGroupFileHash(req.GroupID, groupFileMeta.Directory, groupFileMeta.FileName)
		err = tx.Where("hash = ?", hash).First(&model.GroupFile{}).Error
		if err == gorm.ErrRecordNotFound {
			if err = tx.Create(&model.GroupFile{
				ID:           groupFileID,
				GroupID:      req.GroupID,
				FileID:       file.ID,
				FileName:     groupFileMeta.FileName,
				IsDirectory:  groupFileMeta.IsDirectory,
				Directory:    groupFileMeta.Directory,
				UploadAt:     time.Unix(groupFileMeta.UploadAt, 0),
				LastUpdateAt: time.Unix(groupFileMeta.LastUpdateAt, 0),
				Hash:         hash, // 保证唯一性
				Status:       int(groupFileMeta.Status),
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
		res.FileMeta = groupFileMeta
		res.FileMeta.FileID = file.ID
		if !groupFileMeta.IsDirectory {
			res.FileMeta.Size = fileMeta.Size
		}
		return nil
	})

	if err != nil {
		logger.Errorf("InsertGroupFile failed, for the reason:%v", err)
		if err == common.ErrParentDirNotFound {
			res.Err = getProtoError(err, common.DBNotFoundCode)
			return nil
		} else if err == common.ErrConflict {
			res.Err = getProtoError(err, common.DBConflictCode)
			return nil
		}
		return err
	}
	return nil
}

func (s *Service) DeleteGroupFile(ctx context.Context, req *proto.DeleteGroupFileReq, res *proto.DeleteGroupFileResp) error {
	var (
		err error
	)

	groupFileMeta := &model.GroupFile{}
	logger.Infof("DeleteGroupFile group_id=%v directory=%v filename=%v", req.GroupID, req.Directory, req.FileName)
	err = db.Transaction(func(tx *gorm.DB) error {
		var (
			err   error
			group model.Group
		)

		// 检查 group 是不是存在
		// 防止 group 被删掉 防止组 id 不存在
		if err = tx.Raw(" SELECT * FROM `groups` WHERE id = ? LIMIT 1 LOCK IN SHARE MODE;", req.GroupID).Scan(&group).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.ErrGroupNotExist
			} else {
				return err
			}
		}

		hash := getGroupFileHash(req.GroupID, req.Directory, req.FileName)
		if err = tx.Where("hash = ?", hash).First(groupFileMeta).Error; err != nil {
			return err
		}
		if groupFileMeta.IsDirectory {
			if err = tx.Debug().Where("group_id = ? and directory LIKE ?", groupFileMeta.GroupID, filepath.Join(groupFileMeta.Directory, groupFileMeta.FileName)+"%").Delete(&GroupFile{}).Error; err != nil {
				return err
			}
		}
		// 删除自己 触发删除钩子
		if err = tx.Delete(&groupFileMeta).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("DeleteGroupFile failed, for the reason:%v", err)
		return err
	}

	if err == gorm.ErrRecordNotFound {
		res.Err = getProtoError(err, common.DBNotFoundCode)
		return nil
	}

	// 删除成功
	res.FileMeta = &proto.ListFileMeta{
		FileID:       groupFileMeta.FileID,
		FileName:     groupFileMeta.FileName,
		IsDirectory:  groupFileMeta.IsDirectory,
		UploadAt:     groupFileMeta.UploadAt.Unix(),
		Directory:    groupFileMeta.Directory,
		LastUpdateAt: groupFileMeta.LastUpdateAt.Unix(),
		Status:       int32(groupFileMeta.Status),
	}
	return nil
}

func (s *Service) QueryGroupFile(ctx context.Context, req *proto.QueryGroupFileReq, res *proto.QueryGroupFileResp) error {

	var (
		err error
	)
	groupFile := model.GroupFile{}
	logger.Infof("QueryGroupFile groupID=%v directory=%v filename:%v", req.GroupID, req.Directory, req.FileName)
	err = db.Where("hash = ?", getGroupFileHash(req.GroupID, req.Directory, req.FileName)).First(&groupFile).Error

	if err == gorm.ErrRecordNotFound {
		res.Err = getProtoError(err, common.DBNotFoundCode)
		return nil
	} else if err != nil {
		logger.Errorf("QueryGroupFile failed, for the reason:%v", err)
		return err
	}

	res.FileMeta = &proto.ListFileMeta{
		FileID:       groupFile.FileID,
		FileName:     groupFile.FileName,
		IsDirectory:  groupFile.IsDirectory,
		UploadAt:     groupFile.UploadAt.Unix(),
		Directory:    groupFile.Directory,
		LastUpdateAt: groupFile.LastUpdateAt.Unix(),
		Status:       int32(groupFile.Status),
	}
	return nil
}
