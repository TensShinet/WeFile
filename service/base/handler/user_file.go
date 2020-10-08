package handler

import (
	"github.com/gin-gonic/gin"
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
}

// swagger:parameters GetUserFileListParam
type GetUserFileListParam struct {
	// 存储 session id
	// Required: true
	Cookie string `json:"cookie"`
	// in: body
	Body struct {
		// 目录
		// Required: true
		// 比如 / /dir /dir1/dir2
		Directory string `json:"directory"`
		// csrf token
		// Required: true
		Token string `json:"token"`
	}
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

// swagger:parameters GetUploadAddress
type GetUploadAddressParam struct {
	Cookie string `json:"cookie"`
	// in: body
	Body struct {
		// 文件名
		// Required: true
		FileName string `json:"file_name"`
		// 目录
		// 比如 / /dir /dir1/dir2
		// Required: true
		Directory string `json:"directory"`
		// csrf Token
		// Required: true
		Token string `json:"token"`
	}
}

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

// swagger:parameters GetDownloadAddress
type GetDownloadAddressParam struct {
	Cookie string `json:"cookie"`
	// in: body
	Body struct {
		// 文件唯一 ID
		// Required: true
		FileID int64 `json:"file_id"`
		// 文件名
		// Required: true
		FileName string `json:"file_name"`
		// csrf Token
		// Required: true
		Token string `json:"token"`
	}
}

// swagger:route POST /user/{id}/download_address User GetDownloadAddress
//
// GetDownloadAddress
//
// 得到用户文件下载地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
//		 401: UnauthorizedResponse
//       500: ErrorResponse
func GetDownloadAddress(c *gin.Context) {}
