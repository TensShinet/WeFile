package handler

import (
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/TensShinet/WeFile/service/file/conf"
	"github.com/TensShinet/WeFile/store"
	"github.com/TensShinet/WeFile/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io"
	"math"
	"net/http"
	"strings"
	"syscall"
	"time"
)

// ChunkSize ChunkCount int 足矣 100TB 没问题
// swagger:model MultipartUploadInfo
type MultipartUploadInfo struct {
	FileSize   int64  `json:"file_size"`
	UploadID   string `json:"upload_id"`
	ChunkSize  int    `json:"chunk_size"`
	ChunkCount int    `json:"chunk_count"`
}

// Multipart Upload Init Response
// swagger:response MultipartUploadInitResponse
type MultipartUploadInitResponse struct {
	// 采用 jwt 的方法认证 更新认证信息
	Authorization string
	// in: body
	multipartUploadInfo *MultipartUploadInfo
}

// 分块上传初始化
// swagger:parameters MultipartUploadInit
type MultipartUploadInitParam struct {
	// 采用 jwt 验证方式
	// in: header
	Authorization string
	// in: body
	Body struct {
		// 文件大小单位 Byte
		// Required: true
		FileSize string `json:"file_size"`
	}
}

const (
	defaultUploadIDLength = 30
)

// swagger:route POST /multipart_upload/init File MultipartUploadInit
//
// MultipartUploadInit
//
// 分块上传接口
//     Responses:
//       200: MultipartUploadInitResponse
// 		 401: UnauthorizedResponse
//       500: ErrorResponse
func MultipartUploadInit(c *gin.Context) {

	var (
		err      error
		fileSize int64
	)

	fileSize, err = utils.ParseInt64(c.Request.FormValue("file_size"))

	if err != nil {
		common.SetBadRequestResponse(c, err.Error())
		return
	}

	config := conf.GetConfig()

	// 获取一个连接
	rConn := redisPool.Get()
	defer rConn.Close()

	uploadInfo := MultipartUploadInfo{
		FileSize: fileSize,
		// TODO: 更好的生成 upload id de 办法因为要保证 upload id 唯一且不可猜出来
		UploadID:   utils.RandomString(defaultUploadIDLength),
		ChunkSize:  config.ChunkSize,
		ChunkCount: int(math.Ceil(float64(fileSize) / float64(config.ChunkSize))),
	}
	if _, err = rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunk_count", uploadInfo.ChunkCount); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err = rConn.Do("HSET", "MP_"+uploadInfo.UploadID, "file_size", uploadInfo.FileSize); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, uploadInfo)

}

// swagger:model ChunkInfo
type ChunkInfo struct {
	Size       int `json:"size"`
	ChunkIndex int `json:"chunk_index"`
}

// Multipart Upload Response
// swagger:response MultipartUploadResponse
type MultipartUploadResponse struct {
	// 采用 jwt 的方法认证 更新认证信息
	Authorization string
	// in: body
	ChunkInfo *ChunkInfo
}

// swagger:parameters MultipartUpload
type MultipartUploadParam struct {
	// 采用 jwt 验证方式
	// in: header
	// Required: true
	Authorization string
	// in: body
	Body struct {
		// 上传 id
		// Required: true
		UploadID string `json:"upload_id"`
		// 块索引 下标从 1 开始
		// Required: true
		ChunkIndex int `json:"chunk_index"`
	}
}

// swagger:route POST /multipart_upload File MultipartUpload
//
// MultipartUpload
//
// 分块上传接口
//     Responses:
//       200: MultipartUploadResponse
//		 400: BadRequestResponse
// 		 401: UnauthorizedResponse
//       500: ErrorResponse
func MultipartUpload(c *gin.Context) {
	var (
		err           error
		readLen       int64
		chunkIndex    int
		chunkIndexStr string
		chunk         store.Chunk
	)

	uploadID := c.Request.FormValue("upload_id")
	logger.Debugf("MultipartUpload :%v", uploadID)
	if !isUploadIDExist(uploadID) {
		common.SetSimpleResponse(c, http.StatusBadRequest, "invalid upload_id")
		return
	}
	// TODO: 检查 chunkIndex 的大小
	if chunkIndex, err = utils.ParseInt(c.Request.FormValue("chunk_index")); err != nil {
		common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	chunkIndexStr = c.Request.FormValue("chunk_index")

	logger.Infof("MultipartUpload uploadID:%v chunkIndex:%v", uploadID, chunkIndexStr)
	config := conf.GetConfig()

	rConn := redisPool.Get()
	defer rConn.Close()

	if chunk, err = fileStore.TempChunk(uploadID, chunkIndex, store.ChunkLimit{MaxSize: config.ChunkSize}); err != nil {
		logger.Errorf("MultipartUpload fileStore.TempChunk err:%v", err)
		if err == syscall.EEXIST {
			common.SetSimpleResponse(c, http.StatusBadRequest, "chunk exists")
		} else {
			common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
		logger.Infof("MultipartUpload fc.Request.FormFile err:%v", err)
		return
	}

	if readLen, err = io.Copy(chunk, file); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		logger.Errorf("MultipartUpload io.Copy err:%v", err)
		return
	}

	// 更新 redis 缓存状态
	if _, err = rConn.Do("HSET", "MP_"+uploadID, "chunk_"+chunkIndexStr, 1); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		logger.Errorf("MultipartUpload rConn.Do err:%v", err)
		return
	}

	c.JSON(http.StatusOK, ChunkInfo{
		Size:       int(readLen),
		ChunkIndex: chunkIndex,
	})
}

// 完成分块上传
// swagger:parameters CompleteUpload
type CompleteUploadParam struct {
	// 采用 jwt 验证方式
	// in: header
	// Required: true
	Authorization string
	// in: body
	Body struct {
		// 上传 id
		//
		// Required: true
		UploadID string `json:"upload_id"`
	}
}

// swagger:route POST /multipart_upload/complete File CompleteUpload
//
// CompleteUpload
//
// 合并上传分块
//     Responses:
//       200: UploadResponse
// 		 400: BadRequestResponse
// 		 401: UnauthorizedResponse
//		 409: ConflictError
//       500: ErrorResponse
func CompleteUpload(c *gin.Context) {
	var (
		err       error
		savedFile store.File
		chunkNum  int
	)
	val, _ := c.Get(defaultAuthKey)
	fileMeta, _ := val.(*auth.UploadFileMeta)

	uploadID := c.Request.FormValue("upload_id")
	if !isUploadIDExist(uploadID) {
		common.SetSimpleResponse(c, http.StatusBadRequest, "invalid upload_id")
		return
	}

	logger.Infof("MultipartUpload uploadID:%v", uploadID)
	if chunkNum, err = checkChunkNum(uploadID); err != nil {
		logger.Errorf("CompleteUpload failed in checkChunkNum, for the reason:%v", err)
		if err == errChunkNumLess || err == errChunkNumMore {
			common.SetSimpleResponse(c, http.StatusBadRequest, err.Error())
		} else {
			common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 将所有 chunk 合并起来
	if savedFile, err = fileStore.MergeChunksToSave(fileStore.GetChunksPath(uploadID), chunkNum, store.FileLimit{MaxSize: maxFileSize}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		logger.Errorf("CompleteUpload failed in fileStore.MergeChunksToSave, for the reason: %v", err.Error())
		return
	}

	t := time.Now().Unix()
	logger.Debugf("savedFile.SamplingHash:%v savedFile.Size():%v savedFile.Location():%v", savedFile.SamplingHash(), savedFile.Size(), savedFile.Location())
	if err = insertFile(c, fileMeta, &db.FileMeta{
		Hash:          savedFile.TotalHash(),
		SamplingHash:  savedFile.SamplingHash(),
		HashAlgorithm: "SHA256",
		Size:          savedFile.Size(),
		Location:      savedFile.Location(),
		CreateAt:      t,
		Status:        0,
	}); err != nil {
		_ = savedFile.Remove()
	}

	// 删除上传的所有信息
	rConn := redisPool.Get()
	defer rConn.Close()

	if _, err := rConn.Do("DEL", "MP_"+uploadID); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		logger.Errorf("CompleteUpload failed in fileStore.MergeChunksToSave, for the reason: %v", err.Error())
	}

}

// 完成快速上传
// swagger:parameters UploadProgress
type UploadProgressParam struct {
	// 采用 jwt 验证方式
	// in: header
	// Required: true
	Authorization string
	// 上传 id
	//
	// in: query
	// Required: true
	UploadID string `json:"upload_id"`
}

// 上传进度
// swagger:response UploadProgressResponse
type UploadProgressResponse struct {
	// 采用 jwt 的方法认证 更新认证信息
	Authorization string
	// in: body
	Body struct {
		// chunk 数组
		//
		// 编号哪些 chunk 上传了
		Chunks []int `json:"chunks"`
	}
}

// swagger:route GET /multipart_upload/progress File UploadProgress
//
// UploadProgress
//
// 合并上传分块
//     Responses:
//       200: UploadProgressResponse
// 		 400: BadRequestResponse
// 		 401: UnauthorizedResponse
//       500: ErrorResponse
func UploadProgress(c *gin.Context) {
	uploadID := c.Query("upload_id")
	if !isUploadIDExist(uploadID) {
		common.SetSimpleResponse(c, http.StatusBadRequest, "invalid upload_id")
		return
	}

	rConn := redisPool.Get()
	defer rConn.Close()

	// 需要的 chunk num 和 已经上传的 chunk num 比较
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil || len(data) < 4 {
		if err == nil {
			common.SetSimpleResponse(c, http.StatusBadRequest, "invalid upload_id")
			return
		}
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		logger.Error("UploadProgress failed in redis.Values, for the reason:%v", err)
		return
	}

	chunks, idx := make([]int, (len(data)-4)/2), 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if strings.HasPrefix(k, "chunk_") && v == "1" {
			chunks[idx], _ = utils.ParseInt(strings.TrimPrefix(k, "chunk_"))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"chunks": chunks,
	})

}

func checkChunkNum(uploadID string) (int, error) {
	var (
		err          error
		chunkCount   = 0
		needChunkNum = 0
	)

	rConn := redisPool.Get()
	defer rConn.Close()

	// 需要的 chunk num 和 已经上传的 chunk num 比较
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunk_count" {
			needChunkNum, _ = utils.ParseInt(v)
		} else if strings.HasPrefix(k, "chunk_") && v == "1" {
			chunkCount++
		}
	}

	if needChunkNum > chunkCount {
		return 0, errChunkNumLess
	} else if needChunkNum < chunkCount {
		return 0, errChunkNumMore
	}

	return needChunkNum, nil
}

func isUploadIDExist(uploadID string) bool {
	var (
		err  error
		data interface{}
	)
	rConn := redisPool.Get()
	defer rConn.Close()
	logger.Debugf("isUploadIDExist :%v", uploadID)
	if data, err = rConn.Do("HEXISTS", "MP_"+uploadID, "file_size"); err != nil {
		return false
	}

	if t, _ := data.(int64); t != 1 {
		return false
	}

	return true
}
