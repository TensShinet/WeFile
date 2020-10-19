package handler

import (
	"bytes"
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/service/file/conf"
	"github.com/TensShinet/WeFile/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// swagger:parameters Upload
type UploadParam struct {
	// 采用 jwt 验证方式
	// in: header
	Authorization string
}

// Upload Response
// swagger:response UploadResponse
type UploadResponse struct {
	// 文件 id 用于下载
	FileID int64 `json:"file_id"`
	// 文件名
	FileName string `json:"file_name"`
	// 上传时间
	UploadAt time.Time `json:"upload_at"`
	// 是否是目录
	IsDirectory bool `json:"is_directory"`
	// 文件大小 单位 Byte
	FileSize int64 `json:"file_size"`
}

// swagger:route POST /upload File Upload
//
// Upload
//
// 用户上传
//
//     Responses:
//       200: UploadResponse
//		 400: BadRequestResponse
//		 401: UnauthorizedResponse
//       500: ErrorResponse
func Upload(c *gin.Context) {
	val, ok := c.Get(defaultAuthKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "auth load failed"})
		return
	}
	fileMeta, _ := val.(*auth.UploadFileMeta)
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		logger.Infof("Upload StatusBadRequest err:%v", err.Error())
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: err.Error()})
		return
	}
	logger.Infof("Upload filename:%v size:%v", head.Filename, head.Size)
	defer file.Close()

	// 小文件上传
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		logger.Errorf("Upload StatusInternalServerError err:%v", err.Error())
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}

	if head.Filename != fileMeta.FileName {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "invalid filename"})
		return
	}

	config := conf.GetConfig()
	// 本地保存小文件
	hash := utils.Digest256(buf.Bytes())
	location := filepath.Join(config.FileAPI.LocalTempStore, hash, hash)
	err = os.MkdirAll(filepath.Join(config.FileAPI.LocalTempStore, hash), 0700)
	logger.Infof("Upload WriteFile cao:%v", err)
	if err := ioutil.WriteFile(location, buf.Bytes(), 0600); err != nil {
		logger.Infof("Upload WriteFile err:%v", err.Error())
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}
	logger.Debugf("Upload filename:%v hash:%v location:%v", head.Filename, hash, location)
	// 插入数据库
	t := time.Now().Unix()
	res1, err := dbService.InsertUserFile(c, &db.InsertUserFileMetaReq{
		UserFileMeta: &db.UserFileMeta{
			FileName:     head.Filename,
			IsDirectory:  false,
			UploadAt:     t,
			Directory:    fileMeta.Directory,
			LastUpdateAt: t,
			Status:       0,
		},
		FileMeta: &db.FileMeta{
			Hash:          hash,
			HashAlgorithm: "SHA256",
			Size:          int64(len(buf.Bytes())),
			Location:      location,
			CreateAt:      t,
			Status:        0,
		},
		UserID: fileMeta.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		_ = os.RemoveAll(location)
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		FileID:      res1.FileMeta.FileID,
		FileName:    head.Filename,
		UploadAt:    time.Unix(t, 0),
		IsDirectory: res1.FileMeta.IsDirectory,
		FileSize:    int64(len(buf.Bytes())),
	})
}
