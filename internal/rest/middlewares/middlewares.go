package middlewares

import (
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/netbill/restkit/mdlv"
	"github.com/netbill/restkit/tokens"
)

type Provider struct {
	accountAccessSK string
	uploadFilesSK   string

	log *logium.Logger
}

type Config struct {
	AccountAccessSK string
	UploadFilesSK   string
}

func New(
	log *logium.Logger,
	cfg Config,
) Provider {
	return Provider{
		accountAccessSK: cfg.AccountAccessSK,
		uploadFilesSK:   cfg.UploadFilesSK,
		log:             log,
	}
}

func (p Provider) AccountAuth() func(next http.Handler) http.Handler {
	return mdlv.AccountAuth(p.log, accountDataCtxKey, p.accountAccessSK)
}

func (p Provider) AccountRoleGrant(
	allowedRoles map[string]bool,
) func(http.Handler) http.Handler {
	return mdlv.AccountRoleGrant(p.log, accountDataCtxKey, allowedRoles)
}

func (p Provider) UpdateOrganization() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userData, ok := ctx.Value(accountDataCtxKey).(tokens.AccountJwtData)
			if !ok {
				p.log.Errorf("missing account data in context")
				ape.RenderErr(w, problems.Unauthorized("missing account data in context"))
				return
			}

			confirm := mdlv.ConfirmUploadFiles(
				p.log,
				uploadFilesCtxKey,
				p.uploadFilesSK,
				mdlv.ConfirmUploadFilesParams{
					Audience:   tokenmanager.OrganizationsActor,
					Resource:   tokenmanager.UploadOrgResource,
					ResourceID: userData.AccountID.String(),
				},
			)

			confirm(next).ServeHTTP(w, r)
		})
	}
}
