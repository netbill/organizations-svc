package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
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

	log = log.WithField("organization_id", organizationID).
		WithField("limit", limit).WithField("offset", offset)

	invites, err := c.modules.Invite.GetListForOrganization(
		r.Context(),
		scope.AccountActor(r),
		organizationID,
		limit, offset,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
		return
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to access organization invites")
		render.ResponseError(w, problems.Forbidden("not enough rights to access organization invites"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get organization invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.InvitesCollectionOption, 0)

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
