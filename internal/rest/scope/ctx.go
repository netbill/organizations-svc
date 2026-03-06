package scope

import (
	"context"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/media"
	"github.com/netbill/organizations-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	AccountDataCtxKey
	BaseURLCtxKey
)

func CtxLog(ctx context.Context, log *log.Logger) context.Context {
	return context.WithValue(ctx, LogCtxKey, log)
}

func Log(r *http.Request) *log.Logger {
	logger := r.Context().Value(LogCtxKey).(*log.Logger)

	authClaims, ok := r.Context().Value(AccountDataCtxKey).(tokens.AccountAuthClaims)
	if ok {
		logger = logger.WithAccountAuthClaims(authClaims)
	}

	return logger.WithRequest(r)
}

func CtxAccountAuth(ctx context.Context, accountData tokens.AccountAuthClaims) context.Context {
	return context.WithValue(ctx, AccountDataCtxKey, accountData)
}

func AccountActor(r *http.Request) models.AccountActor {
	claims := r.Context().Value(AccountDataCtxKey).(tokens.AccountAuthClaims)
	return claims.GetAccountID()
}

func CtxUrlResolver(ctx context.Context, resolver *media.Resolver) context.Context {
	return context.WithValue(ctx, BaseURLCtxKey, resolver)
}

func ResolverURL(r *http.Request, key string) (url string) {
	return r.Context().Value(BaseURLCtxKey).(*media.Resolver).Resolve(key)
}
