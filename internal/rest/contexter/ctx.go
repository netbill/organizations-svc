package contexter

import (
	"context"
	"fmt"

	"github.com/netbill/restkit/tokens"
)

const (
	AccountDataCtxKey   = iota
	UploadContentCtxKey = iota
)

func AccountData(ctx context.Context) (tokens.AccountClaims, error) {
	if ctx == nil {
		return tokens.AccountClaims{}, fmt.Errorf("missing context")
	}

	userData, ok := ctx.Value(AccountDataCtxKey).(tokens.AccountClaims)
	if !ok {
		return tokens.AccountClaims{}, fmt.Errorf("missing context")
	}

	if err := userData.Validate(); err != nil {
		return tokens.AccountClaims{}, fmt.Errorf("invalid account data in context: %w", err)
	}

	return userData, nil
}

func UploadContentData(ctx context.Context) (tokens.UploadContentClaims, error) {
	if ctx == nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("missing context")
	}

	userData, ok := ctx.Value(UploadContentCtxKey).(tokens.UploadContentClaims)
	if !ok {
		return tokens.UploadContentClaims{}, fmt.Errorf("missing context")
	}

	if err := userData.Validate(); err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload content data in context: %w", err)
	}

	return userData, nil
}
