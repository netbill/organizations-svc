package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeclineInvite = "decline_invite"

func (c *Controller) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeclineInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Info("invalid invite id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.modules.Invite.Decline(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotFound):
		log.Info("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.Info("invite not for this account")
		render.ResponseError(w, problems.Forbidden("invite not for this account"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.Info("invite already answered")
		render.ResponseError(w, problems.Conflict("invite already answered"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.Info("invite has expired")
		render.ResponseError(w, problems.Forbidden("invite has expired"))
	case err != nil:
		log.WithError(err).Error("failed to decline invite")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Invite(res))
	}
}
