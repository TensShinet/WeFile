package handler

import (
	"context"
	"fmt"
	"github.com/TensShinet/WeFile/conf"
	"github.com/TensShinet/WeFile/logging"
	"github.com/TensShinet/WeFile/service/auth/proto"
	"github.com/TensShinet/WeFile/service/common"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var logger = logging.GetLogger("auth_service_handler")

type Service struct{}

type downloadJWTClaims struct {
	jwt.StandardClaims

	// 追加 file 信息
	FileName string `json:"file_name"`
	FileID   int64  `json:"file_id"`
}

func (s *Service) DownloadJWTEncode(_ context.Context, fileMeta *proto.DownloadFileMeta, res *proto.EncodeResp) error {
	config := conf.GetConfig()
	claims := &downloadJWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.JWT.ValidTime)).Unix(),
		},
		fileMeta.FileName,
		fileMeta.FileID,
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := raw.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
		return err
	}
	res.Token = token
	return nil
}

func (s *Service) DownloadJWTDecode(_ context.Context, req *proto.DecodeReq, res *proto.DownloadJWTDecodeResp) error {
	config := conf.GetConfig()
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT.Secret), nil
	})

	if err != nil {
		if token != nil && !token.Valid {
			res.Err = &proto.Error{
				Code:    common.UnauthorizedCode,
				Message: err.Error(),
			}
		} else {
			res.Err = &proto.Error{
				Code:    -1,
				Message: err.Error(),
			}
			return err
		}
	}

	if token != nil {
		claims, _ := token.Claims.(jwt.MapClaims)

		fileID, _ := claims["file_id"].(float64)
		fileName, _ := claims["file_name"].(string)

		res.FileMeta = &proto.DownloadFileMeta{
			FileID:   int64(fileID),
			FileName: fileName,
		}
	} else {
		logger.Panic("impossible!")
	}

	return nil
}

type uploadJWTClaims struct {
	jwt.StandardClaims

	// 追加 file 信息
	FileName  string `json:"file_name"`
	Directory string `json:"directory"`
	UserID    int64  `json:"user_id"`
}

func (s *Service) UploadJWTEncode(_ context.Context, fileMeta *proto.UploadFileMeta, res *proto.EncodeResp) error {
	config := conf.GetConfig()
	claims := &uploadJWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.JWT.ValidTime)).Unix(),
		},
		fileMeta.FileName,
		fileMeta.Directory,
		fileMeta.UserID,
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := raw.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		res.Err = &proto.Error{
			Code:    -1,
			Message: err.Error(),
		}
		return err
	}
	res.Token = token
	return nil
}

func (s *Service) UploadJWTDecode(_ context.Context, req *proto.DecodeReq, res *proto.UploadJWTDecodeResp) error {
	config := conf.GetConfig()
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT.Secret), nil
	})

	if err != nil {
		if token != nil && !token.Valid {
			res.Err = &proto.Error{
				Code:    common.UnauthorizedCode,
				Message: err.Error(),
			}
		} else {
			res.Err = &proto.Error{
				Code:    -1,
				Message: err.Error(),
			}
			return err
		}
	}

	if token != nil {

		claims, _ := token.Claims.(jwt.MapClaims)

		userID, _ := claims["user_id"].(float64)
		fileName, _ := claims["file_name"].(string)
		directory, _ := claims["directory"].(string)

		res.FileMeta = &proto.UploadFileMeta{
			FileName:  fileName,
			UserID:    int64(userID),
			Directory: directory,
		}
	} else {
		logger.Panic("impossible!")
	}

	return nil
}
