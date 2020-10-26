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

func insertUserFile(c *gin.Context, fileMeta *auth.UploadFileMeta, dbFileMeta *db.FileMeta) error {

	var (
		err            error
		insertFileResp *db.InsertUserFileMetaResp
	)

	t := time.Now().Unix()
	if insertFileResp, err = dbService.InsertUserFile(c, &db.InsertUserFileMetaReq{
		UserFileMeta: &db.UserFileMeta{
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

	if insertFileResp.Err != nil && insertFileResp.Err.Code == common.DBConflictCode {
		common.SetSimpleResponse(c, http.StatusConflict, common.ErrConflict.Error())
		return err
	}

	c.JSON(http.StatusOK, UploadResponse{
		FileID:      insertFileResp.FileMeta.FileID,
		FileName:    fileMeta.FileName,
		UploadAt:    time.Unix(t, 0),
		IsDirectory: insertFileResp.FileMeta.IsDirectory,
		FileSize:    insertFileResp.FileMeta.Size,
	})
	return nil
}
