package controller

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetMyOrganizations = "get_my_organizations"

func (c *Controller) GetMyOrganizations(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyOrganizations)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	accountID := scope.AccountActor(r)
	log = log.WithField("account_id", accountID).WithField("limit", limit).WithField("offset", offset)

	res, err := c.modules.Organization.GetForUser(r.Context(), accountID, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organizations")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.OrgCollectionOption, 0)

	if slices.Contains(includes, "members") {
		organizationIDs := make([]uuid.UUID, res.Size)
		for i, organization := range res.Data {
			organizationIDs[i] = organization.ID
		}

		members, err := c.modules.Member.GetByAccountAndOrgs(r.Context(), accountID, organizationIDs)
		if err != nil {
			log.WithError(err).Error("failed to get members for organizations")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithOrganizationMembers(members))
	}

	render.Response(w, http.StatusOK, responses.Organizations(r, res, opts...))
}
