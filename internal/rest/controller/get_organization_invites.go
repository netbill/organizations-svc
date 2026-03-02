package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetOrganizationInvites = "get_organization_invites"

func (c *Controller) GetOrganizationInvites(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationInvites)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID).
		WithField("limit", limit).
		WithField("offset", offset)

	invites, err := c.modules.Invite.GetListForOrganization(
		r.Context(),
		scope.AccountActor(r),
		organizationID,
		limit, offset,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
		return
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("not enough rights to access organization invites")
		render.ResponseError(w, problems.Forbidden("not enough rights to access organization invites"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get organization invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.InvitesCollectionOption, 0, 2)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 2)

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

	if slices.Contains(includes, "profile") {
		profileIDs := make([]uuid.UUID, 0, invites.Size)
		for _, invite := range invites.Data {
			profileIDs = append(profileIDs, invite.AccountID)
		}

		profiles, err := c.modules.Profile.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionInvitesProfiles(profiles))
	}

	if slices.Contains(includes, "organization") {
		organization, err := c.modules.Organization.GetByID(r.Context(), organizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionInvitesOrganization(organization))
	}

	render.Response(w, http.StatusOK, responses.Invites(r, invites, opts...))
}
