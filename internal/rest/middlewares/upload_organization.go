package middlewares

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/headers"
	"github.com/netbill/restkit/problems"
)

func (p *Provider) UpdateOrganizationMediaContent() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor := scope.AccountActor(r)

			token, err := headers.GetUploadContent(r)
			if err != nil {
				scope.Log(r).Debug("upload token missing")
				p.responser.RenderErr(w, problems.Unauthorized("failed to get token"))

				return
			}

			uploadClaims, err := p.tokenManager.ParseUploadOrganizationContentToken(token)
			if err != nil {
				scope.Log(r).Info("upload token invalid")
				p.responser.RenderErr(w, problems.Unauthorized("invalid upload organization token"))

				return
			}

			if uploadClaims.GetAccountID() != actor {
				scope.Log(r).Info("account is not owner of the organization")
				p.responser.RenderErr(w, problems.Unauthorized("account is not owner of the organization"))

				return
			}

			next.ServeHTTP(w, r.WithContext(scope.CtxUploadContent(r.Context(), uploadClaims)))
		})
	}
}
