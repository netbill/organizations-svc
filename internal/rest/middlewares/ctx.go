package middlewares

import (
	"context"
	"fmt"

	"github.com/netbill/restkit/tokens"
)

const (
	accountDataCtxKey = iota
	uploadFilesCtxKey = iota
)

func AccountData(ctx context.Context) (tokens.AccountJwtData, error) {
	if ctx == nil {
		return tokens.AccountJwtData{}, fmt.Errorf("missing context")
	}

	userData, ok := ctx.Value(accountDataCtxKey).(tokens.AccountJwtData)
	if !ok {
		return tokens.AccountJwtData{}, fmt.Errorf("missing context")
	}

	return userData, nil
}

func UploadFilesData(ctx context.Context) (tokens.UploadFilesJwtData, error) {
	if ctx == nil {
		return tokens.UploadFilesJwtData{}, fmt.Errorf("missing context")
	}

	userData, ok := ctx.Value(uploadFilesCtxKey).(tokens.UploadFilesJwtData)
	if !ok {
		return tokens.UploadFilesJwtData{}, fmt.Errorf("missing context")
	}

	return userData, nil
}
