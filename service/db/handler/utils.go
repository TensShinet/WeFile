package handler

import (
	"context"
	"fmt"
	dbProto "github.com/TensShinet/WeFile/service/db/proto"
	idgProto "github.com/TensShinet/WeFile/service/id_generator/proto"
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
