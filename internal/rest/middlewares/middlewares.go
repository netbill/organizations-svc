package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/restkit/grants"
	"github.com/netbill/restkit/problems"
)

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type Provider struct {
	log             *logium.Logger
	accountAccessSK string
	uploadFilesSK   string

	responser responser
}

type Config struct {
	AccountAccessSK string
	UploadFilesSK   string
}

func New(
	log *logium.Logger,
	cfg Config,
	responser responser,
) *Provider {
	return &Provider{
		accountAccessSK: cfg.AccountAccessSK,
		uploadFilesSK:   cfg.UploadFilesSK,
		log:             log,
		responser:       responser,
	}
}

func (p *Provider) AccountAuth(
	allowedRoles ...string,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res, err := grants.AccountAuthToken(
				r,
				p.accountAccessSK,
				"",
				allowedRoles...,
			)
			if err != nil {
				p.log.WithError(err).Errorf("account authentication failed")
				p.responser.RenderErr(w, problems.Unauthorized("account authentication failed"))

				return
			}

			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), contexter.AccountDataCtxKey, res)),
			)
		})
	}
}

func (p *Provider) UpdateOrganization() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res, err := grants.UploadContentGrant(r, p.uploadFilesSK, grants.UploadContentParams{
				Audience:   tokenmanager.OrganizationsActor,
				Resource:   tokenmanager.OrgResource,
				ResourceID: chi.URLParam(r, "organization_id"),
			})
			if err != nil {
				p.log.WithError(err).Errorf("upload content grant validation failed")
				p.responser.RenderErr(w, problems.Unauthorized("upload content grant validation failed"))

				return
			}

			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), contexter.UploadContentCtxKey, res)),
			)
		})
	}
}
