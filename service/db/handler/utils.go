package handler

import (
	"context"
	"fmt"
)

func getID(ctx context.Context) (int64, error) {
	idResp, err := generateIDService.GenerateID(ctx, nil)
	if err != nil {
		return -1, err
	}
	if idResp.Err != nil {
		return -1, fmt.Errorf(idResp.Err.Message)
	}
	return idResp.Id, nil
}
