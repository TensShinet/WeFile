package handler

import (
	"context"
	"fmt"
	dbProto "github.com/TensShinet/WeFile/service/db/proto"
	idgProto "github.com/TensShinet/WeFile/service/id_generator/proto"
	"github.com/TensShinet/WeFile/utils"
	"path/filepath"
	"strconv"
)

func getID(ctx context.Context) (int64, error) {
	idResp, err := generateIDService.GenerateID(ctx, new(idgProto.IDReq))
	if err != nil {
		return -1, err
	}
	if idResp.Err != nil {
		return -1, fmt.Errorf(idResp.Err.Message)
	}
	return idResp.Id, nil
}

func getProtoError(err error, code int) *dbProto.Error {
	return &dbProto.Error{Code: int32(code), Message: err.Error()}
}

func getUserFileHash(id int64, directory, filename string) string {
	s := strconv.FormatInt(id, 16) + filepath.Join(directory, filename)
	return utils.Digest256([]byte(s))
}

func getGroupFileHash(id int64, directory, filename string) string {
	return getUserFileHash(id, directory, filename)
}
