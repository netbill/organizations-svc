package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetInvite = "get_invite"

// TODO
func (c *Controller) GetInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	invite, err := c.modules.Invite.GetForAccount(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.Info("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get invite")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.InviteOption, 0)

	if slices.Contains(includes, "profile") {
		profile, err := c.modules.Profile.GetByID(r.Context(), invite.AccountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotExists):
			log.WithError(err).
				WithField("account_id", invite.AccountID).
				Warn("profile not found")
			render.ResponseError(w, problems.NotFound("profile not found"))
			return
		case err != nil:
			log.WithError(err).
				WithField("account_id", invite.AccountID).
				Error("failed to get profile")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithInviteProfile(profile))
		}
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.modules.Organization.GetByID(r.Context(), invite.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists):
			log.WithError(err).Warn("organization not found")
			render.ResponseError(w, problems.NotFound("organization not found"))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithInviteOrganization(org))
		}
	}

	render.Response(w, http.StatusOK, responses.Invite(invite, opts...))
}
