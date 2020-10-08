package handler

import (
	"github.com/gin-gonic/gin"
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
}

// swagger:route POST /Upload File Upload
//
// Upload
//
// 用户上传
//
//     Responses:
//       200: UploadResponse
//       500: ErrorResponse
func Upload(c *gin.Context) {}
