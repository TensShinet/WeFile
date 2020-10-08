package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

// swagger:model File
type File struct {
	// 用户 id
	UserID string `json:"user_id"`
	// 文件 id 用于下载
	FileID string `json:"file_id"`
	// 文件名
	FileName string `json:"file_name"`
	// 上传时间
	UploadAt time.Time `json:"upload_at"`
	// 是否是目录
	IsDirectory bool `json:"is_directory"`
}

// swagger:parameters GetUserFileListParam
type GetUserFileListParam struct {
	// 存储 session id
	// Required: true
	// in: header
	Cookie string `json:"cookie"`
	// 目录
	// Required: true
	// 比如 / /dir /dir1/dir2
	// in: body
	Directory string `json:"directory"`
}

// File List
// swagger:response FileListResponse
type FileListResponse struct {
	// in: body
	Body struct {
		DirectoryName string  `json:"directory_name"`
		Files         []*File `json:"files"`
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

// swagger:route POST /user/{id}/file_list User GetUserFileList
//
// GetUserFileList
//
// 用户目录下的文件
//     Responses:
//       200: FileListResponse
//  	 401: UnauthorizedResponse
//       500: ErrorResponse
func GetUserFileList(c *gin.Context) {}

// swagger:route GET /user/{id}/upload_address User GetUploadAddress
//
// GetUploadAddress
//
// 得到用户上传地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
//       401: UnauthorizedResponse
//       500: ErrorResponse
func GetUploadAddress(c *gin.Context) {}

// swagger:route GET /user/{id}/download_address User GetUploadAddress
//
// GetDownloadAddress
//
// 得到用户文件下载地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
//		 401: UnauthorizedResponse
//       500: ErrorResponse
func GetDownloadAddress(c *gin.Context) {}
