package handler

import (
	"fmt"
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	errChunkNumMore = fmt.Errorf("the amount of uploaded chunks is more than the amount of needed chunks")
	errChunkNumLess = fmt.Errorf("the amount of uploaded chunks is less than the amount of needed chunks")
)

const (
	maxSmallFileSize = 100 * 1024 * 1024               // 100 M
	maxFileSize      = 100 * 1024 * 1024 * 1024 * 1024 // 100 TB
)

func insertFile(c *gin.Context, fileMeta *auth.UploadFileMeta, dbFileMeta *db.FileMeta) error {
	var (
		err                 error
		insertUserFileResp  *db.InsertUserFileMetaResp
		insertGroupFileResp *db.InsertGroupFileResp
	)

	logger.Infof("insertFile userID:%v groupID:%v", fileMeta.UserID, fileMeta.GroupID)
	t := time.Now().Unix()
	if fileMeta.UserID != 0 {
		if insertUserFileResp, err = dbService.InsertUserFile(c, &db.InsertUserFileMetaReq{
			UserFileMeta: &db.ListFileMeta{
				FileName:     fileMeta.FileName,
				IsDirectory:  false,
				UploadAt:     t,
				Directory:    fileMeta.Directory,
				LastUpdateAt: t,
				Status:       0,
			},
			FileMeta: dbFileMeta,
			UserID:   fileMeta.UserID,
		}); err != nil {
			common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
			return err
		}
	} else {
		if insertGroupFileResp, err = dbService.InsertGroupFile(c, &db.InsertGroupFileReq{
			GroupFileMeta: &db.ListFileMeta{
				FileName:     fileMeta.FileName,
				IsDirectory:  false,
				UploadAt:     t,
				Directory:    fileMeta.Directory,
				LastUpdateAt: t,
				Status:       0,
			},
			FileMeta: dbFileMeta,
			GroupID:  fileMeta.GroupID,
		}); err != nil {
			common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
			return err
		}
	}

	if insertUserFileResp != nil && insertUserFileResp.Err != nil {
		common.SetSimpleResponse(c, int(insertUserFileResp.Err.Code), insertUserFileResp.Err.Message)
		return err
	}

	if insertGroupFileResp != nil && insertGroupFileResp.Err != nil {
		common.SetSimpleResponse(c, int(insertGroupFileResp.Err.Code), insertGroupFileResp.Err.Message)
		return err
	}

	if insertUserFileResp != nil {
		c.JSON(http.StatusOK, UploadResponse{
			FileID:      insertUserFileResp.FileMeta.FileID,
			FileName:    fileMeta.FileName,
			UploadAt:    time.Unix(t, 0),
			IsDirectory: insertUserFileResp.FileMeta.IsDirectory,
			FileSize:    insertUserFileResp.FileMeta.Size,
		})
	} else {
		c.JSON(http.StatusOK, UploadResponse{
			FileID:      insertGroupFileResp.FileMeta.FileID,
			FileName:    fileMeta.FileName,
			UploadAt:    time.Unix(t, 0),
			IsDirectory: insertGroupFileResp.FileMeta.IsDirectory,
			FileSize:    insertGroupFileResp.FileMeta.Size,
		})
	}
	return nil
}
