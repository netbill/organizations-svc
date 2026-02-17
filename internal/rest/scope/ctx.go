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

	return log
}

func CtxAccountAuth(ctx context.Context, accountData tokens.AccountAuthClaims) context.Context {
	return context.WithValue(ctx, AccountDataCtxKey, accountData)
}

func AccountActor(r *http.Request) models.AccountActor {
	claims := r.Context().Value(AccountDataCtxKey).(tokens.AccountAuthClaims)
	return claims.GetAccountID()
}
