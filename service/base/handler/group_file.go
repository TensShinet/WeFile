package handler

import "C"
import (
	auth "github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/base/conf"
	"github.com/TensShinet/WeFile/service/common"
	db "github.com/TensShinet/WeFile/service/db/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// swagger:parameters DeleteGroupFile
type DeleteGroupFileParam struct {
	// in: header
	// Required: true
	Cookie string `json:"cookie"`

	// group id
	// in: query
	// Required: true
	GroupID int64 `json:"group_id"`
	// 在 Directory 下的文件名/目录名
	// in: query
	// Required: true
	Name string `json:"name"`
	// 目录
	// 比如 / /dir /dir1/dir2 注意一定要有 /
	// in: query
	// Required: true
	Directory string `json:"directory"`
	// csrf_token
	// 登录的时候才会更新
	// in: query
	// Required: true
	CSRFToken string `json:"csrf_token"`
}

// swagger:route DELETE /user/group/file_list Group DeleteGroupFile
//
// DeleteGroupFile
//
// 删除用户目录下的文件/目录
//     Responses:
//       200: DeleteFileResponse
//		 400: BadRequestResponse
//  	 401: UnauthorizedResponse
// 		 403: ForbiddenResponse
//		 404: NotFoundError
//       500: ErrorResponse
func DeleteGroupFile(c *gin.Context) {
	var (
		groupID int64
		userID  int64
		err     error
		res     *db.DeleteGroupFileResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID
	name := c.Request.FormValue("name")
	directory := c.Request.FormValue("directory")
	logger.Info("DeleteGroupFile groupID:%v userID:%v filename:%v directory:%v", groupID, userID, name, directory)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res, err = dbService.DeleteGroupFile(c, &db.DeleteGroupFileReq{
		GroupID:   groupID,
		Directory: directory,
		FileName:  name,
	}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil && res.Err.Code == common.DBNotFoundCode {
		common.SetSimpleResponse(c, http.StatusNotFound, "not found")
		return
	}

	c.JSON(http.StatusOK, File{
		FileID:      res.FileMeta.FileID,
		FileName:    res.FileMeta.FileName,
		UploadAt:    time.Unix(res.FileMeta.UploadAt, 0),
		IsDirectory: res.FileMeta.IsDirectory,
		FileSize:    res.FileMeta.Size,
	})
}

// swagger:parameters CreateGroupDirectory
type CreateGroupDirectoryParam struct {
	// in: header
	// Required: true
	Cookie string `json:"cookie"`

	// in: body
	Body struct {
		// group id
		// Required: true
		GroupID int64 `json:"group_id"`
		// 在 Directory 下的文件名/目录名
		// Required: true
		Name string `json:"name"`
		// 目录
		// 比如 / /dir /dir1/dir2 注意一定要有 /
		// Required: true
		Directory string `json:"directory"`
		// csrf_token
		// 登录的时候才会更新
		// Required: true
		CSRFToken string `json:"csrf_token"`
	}
}

// swagger:route POST /user/group/file_list Group CreateGroupDirectory
//
// CreateGroupDirectory
//
// 创建一个新目录
//     Responses:
//       200: CreateDirectoryResponse
//		 400: BadRequestResponse
//  	 401: UnauthorizedResponse
//		 403: ForbiddenResponse
//		 409: ConflictError
//       500: ErrorResponse
func CreateGroupDirectory(c *gin.Context) {
	var (
		groupID int64
		userID  int64
		err     error
		res     *db.InsertGroupFileResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroup failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID
	name := c.Request.FormValue("name")
	directory := c.Request.FormValue("directory")
	logger.Info("DeleteGroupFile groupID:%v userID:%v filename:%v directory:%v", groupID, userID, name, directory)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	t := time.Now()
	if res, err = dbService.InsertGroupFile(c, &db.InsertGroupFileReq{
		GroupFileMeta: &db.ListFileMeta{
			FileName:     name,
			IsDirectory:  true,
			UploadAt:     t.Unix(),
			Directory:    directory,
			LastUpdateAt: t.Unix(),
			Status:       0,
		},
		GroupID: groupID,
	}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil {
		common.SetSimpleResponse(c, int(res.Err.Code), res.Err.Message)
		return
	}

	c.JSON(http.StatusOK, File{
		FileName:    res.FileMeta.FileName,
		UploadAt:    time.Unix(res.FileMeta.UploadAt, 0),
		IsDirectory: true,
	})

}

// swagger:route GET /user/group/file_list Group GetGroupFileList
//
// GetGroupFileList
//
// 得到用户上传地址，并设置 jwt token
//     Responses:
//       200: FileListResponse
//  	 401: UnauthorizedResponse
//		 404: NotFoundError
//       500: ErrorResponse
func GetGroupFileList(c *gin.Context) {
	var (
		groupID   int64
		userID    int64
		directory string
		err       error
		res       *db.ListGroupFileResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroupFileList failed, for the reason:%v", err)
		}
	}()

	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	userID = getUser(c).UserID
	directory = c.Query("directory")
	logger.Info("GetGroupFileList groupID:%v userID:%v directory:%v", groupID, userID, directory)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	if res, err = dbService.ListGroupFile(c, &db.ListGroupFileReq{GroupID: groupID, Directory: directory}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil {
		common.SetSimpleResponse(c, int(res.Err.Code), res.Err.Message)
		return
	}

	list := res.GroupFileMetaList
	files := make([]*File, len(res.GroupFileMetaList))
	for i, v := range list {
		files[i] = &File{
			FileID:      v.FileID,
			FileName:    v.FileName,
			UploadAt:    time.Unix(v.UploadAt, 0),
			IsDirectory: v.IsDirectory,
			FileSize:    v.Size,
		}
	}

	c.JSON(http.StatusOK, files)

}

// swagger:route GET /user/group/download_address Group GetGroupDownloadAddress
//
// GetGroupDownloadAddress
//
// 得到用户文件下载地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
//		 400: BadRequestResponse
//		 401: UnauthorizedResponse
//       500: ErrorResponse
func GetGroupDownloadAddress(c *gin.Context) {
	var (
		groupID   int64
		userID    int64
		fileID    int64
		fileName  string
		directory string
		res       *db.QueryGroupFileResp
		authRes   *auth.EncodeResp
		err       error
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroupUploadAddress failed, for the reason:%v", err)
		}
	}()
	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	if fileID, err = getFileIDFromContext(c); err != nil {
		return
	}
	fileName = c.Query("file_name")
	directory = c.Query("directory")
	userID = getUser(c).UserID
	logger.Info("GetGroupDownloadAddress groupID:%v userID:%v fileID:%v filename:%v directory:%v", groupID, userID, fileID, fileName, directory)

	if err = checkUserInGroup(c, userID, groupID); err != nil {
		return
	}

	// 检查是不是存在文件
	if res, err = dbService.QueryGroupFile(c, &db.QueryGroupFileReq{
		GroupID:   groupID,
		Directory: directory,
		FileName:  fileName,
	}); err != nil {
		common.SetSimpleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res.Err != nil && res.Err.Code == common.DBNotFoundCode {
		common.SetSimpleResponse(c, http.StatusNotFound, res.Err.Message)
		return
	}

	// jwt encode
	if authRes, err = authService.DownloadJWTEncode(c, &auth.DownloadFileMeta{
		FileID:   fileID,
		FileName: fileName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Authorization", "Bearer "+authRes.Token)

	config := conf.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"address": config.BaseAPI.FileAPIAddress,
	})

}

// swagger:route GET /user/group/upload_address Group GetGroupUploadAddress
//
// GetGroupUploadAddress
//
// 得到用户上传地址，并设置 jwt token
//     Responses:
//       200: AddressResponse
// 		 400: BadRequestResponse
//       401: UnauthorizedResponse
//       500: ErrorResponse
func GetGroupUploadAddress(c *gin.Context) {

	var (
		err     error
		groupID int64
		res     *auth.EncodeResp
	)
	defer func() {
		if err != nil {
			logger.Errorf("GetGroupUploadAddress failed, for the reason:%v", err)
		}
	}()

	// TODO: 检查
	if groupID, err = getGroupIDFromContext(c); err != nil {
		return
	}
	directory := c.Query("directory")
	fileName := c.Query("file_name")
	logger.Infof("GetGroupUploadAddress groupID:%v directory:%v fileName:%v", groupID, directory, fileName)
	if directory == "" || fileName == "" {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse{
			Message: "目录或者文件名为空",
		})
		return
	}

	// jwt encode
	if res, err = authService.UploadJWTEncode(c, &auth.UploadFileMeta{
		GroupID:   groupID,
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
