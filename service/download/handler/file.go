package handler

import "github.com/gin-gonic/gin"

// swagger:parameters Download
type DownloadParam struct {
	// 采用 jwt 验证方式
	// in: header
	Authorization string
}

// Download Response
// swagger:response DownloadResponse
type DownloadResponse struct {
	ContentDisposition string `json:"Content-Disposition"`
}

// swagger:route POST /download File Download
//
// Download
//
// 用户下载
//     Produces:
//     - application/octect-stream
//
//     Responses:
//       200: DownloadResponse
//       500: ErrorResponse
func Download(c *gin.Context) {}
