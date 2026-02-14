package middlewares

import (
	"net/http"

	"github.com/netbill/restkit/tokens"
)

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type tokenManager interface {
	ParseAccountAuthAccessClaims(token string) (tokens.AccountAuthClaims, error)
	ParseUploadOrganizationContentToken(token string) (tokens.UploadContentClaims, error)
}

type Provider struct {
	tokenManager tokenManager
	responser    responser
}

func New(
	responser responser,
	tokenManager tokenManager,
) *Provider {
	return &Provider{
		tokenManager: tokenManager,
		responser:    responser,
	}
}
