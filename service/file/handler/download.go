package handler

import (
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

// swagger:route GET /download File Download
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
func Download(c *gin.Context) {
	val, ok := c.Get(defaultAuthKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "auth load failed"})
		return
	}
	fileMeta, _ := val.(*auth.DownloadFileMeta)
	logger.Infof("Download fileID:%v, filename:%v", fileMeta.FileID, fileMeta.FileName)
	res1, err := dbService.QueryFileMeta(c, &db.QueryFileMetaReq{
		Id: fileMeta.FileID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "auth load failed"})
		return
	}

	if res1.Err != nil && res1.Err.Code == common.DBNotFoundCode {
		c.JSON(http.StatusNotFound, res1.Err.Message)
		return
	}
	logger.Debugf("Download location:%v filename:%v", res1.FileMeta.Location, fileMeta.FileName)
	c.FileAttachment(res1.FileMeta.Location, fileMeta.FileName)
}
