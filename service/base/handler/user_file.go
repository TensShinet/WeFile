package handler

import (
	"github.com/TensShinet/WeFile/conf"
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// swagger:model File
type File struct {
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

// swagger:parameters GetUserFileList
type GetUserFileListParam struct {
	// in: header
	// 存储 session id
	// Required: true
	Cookie string `json:"cookie"`

	// in: query
	// 目录
	// Required: true
	// 比如 / /dir /dir1/dir2
	Directory string `json:"directory"`
}

// File List
// swagger:response FileListResponse
type FileListResponse struct {
	// in: body
	Body struct {
		Directory string  `json:"directory"`
		Files     []*File `json:"files"`
	} `json:"body"`
}

// Address
// swagger:response AddressResponse
type AddressResponse struct {
	// 采用 jwt 的方法认证
	Authorization string
	// in: body
	Body struct {
		Address string `json:"address"`
	}
}

// swagger:route GET /file_list/{user_id} User GetUserFileList
//
// GetUserFileList
//
// 用户目录下的文件
//     Responses:
//       200: FileListResponse
//  	 401: UnauthorizedResponse
//		 404: NotFoundError
//       500: ErrorResponse
func GetUserFileList(c *gin.Context) {
	var (
		res *db.ListUserFileMetaResp
		err error
	)
	// TODO:检查
	id, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	directory := c.Query("directory")
	logger.Infof("GetUserFileList user_id:%v directory:%v", id, directory)
	if directory == "" {
		c.JSON(http.StatusNotFound, common.Response404{
			Message: "目录不存在",
		})
		return
	}

	if res, err = dbService.ListUserFile(c, &db.ListUserFileMetaReq{
		UserID:    id,
		Directory: directory,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}

	if res.Err != nil && res.Err.Code == common.DBNotFoundCode {
		c.JSON(http.StatusNotFound, common.Response404{
			Message: "目录不存在",
		})
		return
	}

	list := res.UserFileMetaList
	files := make([]*File, len(list))
	for i := 0; i < len(res.UserFileMetaList); i++ {
		files[i] = &File{
			FileID:      list[i].FileID,
			FileName:    list[i].FileName,
			UploadAt:    time.Unix(list[i].UploadAt, 0),
			IsDirectory: list[i].IsDirectory,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"directory": directory,
		"files":     files,
	})
}

// swagger:parameters GetUploadAddress
type GetUploadAddressParam struct {
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// 文件名
	// Required: true
	// in: query
	FileName string `json:"file_name"`
	// 目录
	// 比如 / /dir /dir1/dir2
	// Required: true
	Directory string `json:"directory"`
}

// swagger:route GET /upload_address/{user_id} User GetUploadAddress
//
// GetUploadAddress
//
// 得到用户上传地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//       500: ErrorResponse
func GetUploadAddress(c *gin.Context) {

	var (
		err error
		res *auth.EncodeResp
	)

	// TODO: 检查
	id, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	directory := c.Query("directory")
	fileName := c.Query("file_name")
	logger.Infof("GetUploadAddress id:%v directory:%v fileName:%v", id, directory, fileName)
	if directory == "" || fileName == "" {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{
			Message: "目录或者文件名为空",
		})
		return
	}

	// jwt encode
	if res, err = authService.UploadJWTEncode(c, &auth.UploadFileMeta{
		UserID:    id,
		Directory: directory,
		FileName:  fileName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Authorization", "Bearer "+res.Token)
	config := conf.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"address": config.BaseAPI.FileAPIAddress,
	})
}

// swagger:parameters GetDownloadAddress
type GetDownloadAddressParam struct {
	// in: header
	// Required: true
	Cookie string `json:"cookie"`
	// in: query
	// 文件 ID
	// Required: true
	FileID int64 `json:"file_id"`
	// in: query
	// 文件名
	// Required: true
	FileName string `json:"file_name"`
	// in: query
	// 目录
	// 比如 / /dir /dir1/dir2
	// Required: true
	Directory string `json:"directory"`
}

// swagger:route GET /download_address/{user_id} User GetDownloadAddress
//
// GetDownloadAddress
//
// 得到用户文件下载地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
//		 400: BadRequestResponse
//		 401: UnauthorizedResponse
//       500: ErrorResponse
func GetDownloadAddress(c *gin.Context) {
	var (
		err    error
		fileID int64
		userID int64
		res    *auth.EncodeResp
		res1   *db.QueryUserFileResp
	)

	// TODO: 检查
	if fileID, err = strconv.ParseInt(c.Query("file_id"), 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{Message: "file id 无效"})
		return
	}
	fileName := c.Request.FormValue("file_name")
	userID, _ = strconv.ParseInt(c.Param("user_id"), 10, 64)
	directory := c.Query("directory")

	if directory == "" || fileName == "" {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{
			Message: "目录或者文件名为空",
		})
		return
	}

	logger.Infof("GetDownloadAddress fileID:%v directory:%v fileName:%v", fileID, directory, fileName)

	// 查询文件是否存在
	if res1, err = dbService.QueryUserFile(c, &db.QueryUserFileReq{
		UserID:    userID,
		Directory: directory,
		FileName:  fileName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if res1.Err != nil && res1.Err.Code == common.DBNotFoundCode {
		c.JSON(http.StatusNotFound, common.Response404{Message: "file not found"})
		return
	}

	// jwt encode
	if res, err = authService.DownloadJWTEncode(c, &auth.DownloadFileMeta{
		FileID:   fileID,
		FileName: fileName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Authorization", "Bearer "+res.Token)

	config := conf.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"address": config.BaseAPI.FileAPIAddress,
	})
}
