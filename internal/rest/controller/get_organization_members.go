package controller

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetOrganizationMembers = "get_organization_members"

func (c *Controller) GetOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationMembers)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Info("invalid organization id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	log = log.WithField("organization_id", organizationID).WithField("limit", limit).WithField("offset", offset)

	params := member.FilterParams{
		OrganizationID: &organizationID,
	}

	if v := r.URL.Query().Get("text"); v != "" {
		text := v
		params.BestMatch = &text
	}

	members, err := c.modules.Member.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organization members")
		render.ResponseError(w, problems.InternalError())
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.MembersCollectionOption, 0)

	if slices.Contains(includes, "profile") {
		profileIDs := make([]uuid.UUID, 0, members.Size)
		for _, invite := range members.Data {
			profileIDs = append(profileIDs, invite.AccountID)
		}

		profiles, err := c.modules.Profile.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionMembersProfiles(profiles))
	}

	if slices.Contains(includes, "organization") {
		organization, err := c.modules.Organization.GetByID(r.Context(), organizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionMembersOrganization(organization))
	}

	render.Response(w, http.StatusOK, responses.Members(r, members, opts...))
}
