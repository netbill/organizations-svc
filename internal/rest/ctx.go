package rest

import (
	"fmt"
	"net/http"

	"github.com/netbill/restkit/token"
)

type ctxKey int

const (
	AccountDataCtxKey ctxKey = iota
)

func AccountData(r *http.Request) (token.AccountData, error) {
	if r.Context() == nil {
		return token.AccountData{}, fmt.Errorf("missing context")
	}

	userData, ok := r.Context().Value(AccountDataCtxKey).(token.AccountData)
	if !ok {
		return token.AccountData{}, fmt.Errorf("missing context")
	}

	return userData, nil
}
