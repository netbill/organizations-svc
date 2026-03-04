package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetMyInvites = "get_my_invites"

func (c *Controller) GetMyInvites(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyInvites)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/limit": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	accountID := scope.AccountActor(r)
	log = log.WithField("account_id", accountID).WithField("limit", limit).WithField("offset", offset)

	invites, err := c.modules.Invite.GetListForAccount(r.Context(), accountID, limit, offset)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("account not found")
		render.ResponseError(w, problems.NotFound("account not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.InvitesCollectionOption, 0, 1)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 1)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	organizationIDs := make([]uuid.UUID, 0, invites.Size)
	for _, invite := range invites.Data {
		if !slices.Contains(organizationIDs, invite.OrganizationID) {
			organizationIDs = append(organizationIDs, invite.OrganizationID)
		}
	}

	if slices.Contains(includes, "organizations") {
		organization, err := c.modules.Organization.GetByIDs(r.Context(), organizationIDs)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionInvitesOrganizations(organization))
	}

	render.Response(w, http.StatusOK, responses.Invites(r, invites, opts...))
}
