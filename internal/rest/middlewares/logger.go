package middlewares

import (
	"net/http"

	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/rest/scope"
)

func (p *Provider) Logger(log *logium.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(scope.CtxLog(r.Context(), log)))
		})
	}
}
