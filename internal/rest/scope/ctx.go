package scope

import (
	"context"
	"net/http"

	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/tokens"
)

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	AccountDataCtxKey
	UploadContentCtxKey
)

func CtxLog(ctx context.Context, log *logium.Entry) context.Context {
	return context.WithValue(ctx, LogCtxKey, log)
}

func Log(r *http.Request) *logium.Entry {
	log := r.Context().Value(LogCtxKey).(*logium.Entry)

	authClaims, ok := r.Context().Value(AccountDataCtxKey).(tokens.AccountAuthClaims)
	if ok {
		log = log.WithAccountAuthClaims(authClaims)
	}
	contentClaims, ok := r.Context().Value(UploadContentCtxKey).(tokens.UploadContentClaims)
	if ok {
		log = log.WithUploadContentClaims(contentClaims)
	}

	return log
}

func CtxAccountAuth(ctx context.Context, accountData tokens.AccountAuthClaims) context.Context {
	return context.WithValue(ctx, AccountDataCtxKey, accountData)
}

func AccountActor(r *http.Request) models.AccountActor {
	claims := r.Context().Value(AccountDataCtxKey).(tokens.AccountAuthClaims)
	return claims.GetAccountID()
}

func CtxUploadContent(ctx context.Context, content tokens.UploadContentClaims) context.Context {
	return context.WithValue(ctx, UploadContentCtxKey, content)
}

func UploadScope(r *http.Request) models.UploadScope {
	claims := r.Context().Value(UploadContentCtxKey).(tokens.UploadContentClaims)
	return claims.GetSessionID()
}
