package middlewares

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/media"
	"github.com/netbill/organizations-svc/internal/rest/scope"
)

func (p *Provider) ResolverUrl(resolver *media.Resolver) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(scope.CtxUrlResolver(r.Context(), resolver)))
		})
	}
}
