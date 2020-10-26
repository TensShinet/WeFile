package handler

import (
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/store"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

	var (
		err       error
		tempFile  store.File
		savedFile store.File
	)
	val, _ := c.Get(defaultAuthKey)
	fileMeta, _ := val.(*auth.UploadFileMeta)
	file, head, err := c.Request.FormFile("file")
	if err != nil || head.Size > maxSmallFileSize {
		if err == nil {
			err = store.ErrFileSizeExceed
		}
		logger.Errorf("Upload Request.FormFile err:%v", err.Error())
		common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Infof("Upload filename:%v size:%v", head.Filename, head.Size)
	defer file.Close()

	if tempFile, err = fileStore.TempFile(store.FileLimit{MaxSize: maxSmallFileSize}); err != nil {
		logger.Errorf("Upload fileStore.TempFile err:%v", err.Error())
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := io.Copy(tempFile, file); err != nil {
		logger.Errorf("Upload io.Copy err:%v", err.Error())
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if head.Filename != fileMeta.FileName {
		common.SetSimpleResponse(c, http.StatusBadRequest, "invalid filename")
		return
	}

	// 保存
	if savedFile, err = fileStore.SaveFile(tempFile); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 删除暂存文件
	_ = tempFile.Remove()

	logger.Infof("Upload filename:%v Hash:%v location:%v", head.Filename, tempFile.TotalHash(), savedFile.Location())
	// 插入数据库
	t := time.Now().Unix()
	if err = insertUserFile(c, fileMeta, &db.FileMeta{
		Hash:          savedFile.TotalHash(),
		SamplingHash:  savedFile.SamplingHash(),
		HashAlgorithm: "SHA256",
		Size:          savedFile.Size(),
		Location:      savedFile.Location(),
		CreateAt:      t,
		Status:        0,
	}); err != nil {
		// 删除保存的文件
		_ = savedFile.Remove()
	}
}

// swagger:parameters TryFastUpload
type TryFastUploadParam struct {
	// 采用 jwt 验证方式
	// in: header
	// Required: true
	Authorization string
	// in: body
	Body struct {
		// 采用抽样 Hash 的方式计算 Hash 快速上传
		FileSamplingHash string `json:"file_sampling_hash"`

		// 全量 Hash
		FileHash string `json:"file_hash"`
	}
}

// swagger:route POST /upload/try_fast File TryFastUpload
//
// TryFastUpload
//
// 尝试快速上传
//
//     Responses:
//       200: UploadResponse
//		 400: BadRequestResponse
//		 404: NotFoundError
//		 401: UnauthorizedResponse
//		 409: ConflictError
//       500: ErrorResponse

// 尝试快速上传
// 上传抽样 Hash 检查存不存在
//  如果不存在，可以上传
//	如果存在，需要全量 Hash
// 上传全量 Hash 如果存在直接写入
// 如果不存在需要上传文件
func TryFastUpload(c *gin.Context) {
	val, _ := c.Get(defaultAuthKey)
	fileMeta, _ := val.(*auth.UploadFileMeta)

	var (
		err                        error
		fileHash, fileSamplingHash string
		queryFileResp              *db.QueryFileMetaResp
	)
	fileHash = c.Request.FormValue("file_hash")
	fileSamplingHash = c.Request.FormValue("file_sampling_hash")

	if fileHash == "" && fileSamplingHash == "" {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "invalid file hash and file sampling sash"})
		return
	}
	logger.Infof("TryFastUpload fileHash:%v fileSamplingHash:%v", fileHash, fileSamplingHash)
	// Post 存在 sampling Hash
	if fileHash == "" {
		if queryFileResp, err = dbService.QueryFileMeta(c, &db.QueryFileMetaReq{
			SamplingHash: fileSamplingHash,
			Hash:         fileHash,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
			return
		}
		// 不存在 抽样 Hash 可以直接上传
		if queryFileResp.Err != nil && queryFileResp.Err.Code == common.DBNotFoundCode {
			c.JSON(http.StatusNotFound, common.Response404{Message: "file not found"})
			return
		}
		// 存在 抽样 Hash 需要全量 Hash
		c.JSON(http.StatusAccepted, common.AcceptedResponse{Message: "find sampling hash, need file total hash"})
		return
	}

	if queryFileResp, err = dbService.QueryFileMeta(c, &db.QueryFileMetaReq{
		Hash: fileHash,
	}); err != nil {
		// 全量 Hash 不存在需要上传
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}
	// 全量 Hash 存在 直接上传
	if queryFileResp.Err != nil && queryFileResp.Err.Code == common.DBNotFoundCode {
		c.JSON(http.StatusNotFound, common.Response404{Message: "File not found"})
		return
	}
	_ = insertUserFile(c, fileMeta, queryFileResp.FileMeta)
}
