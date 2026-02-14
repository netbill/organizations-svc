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
)

const operationAcceptInvite = "accept_invite"

func (c *Controller) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationAcceptInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Info("invalid invite id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.core.invite.Accept(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotFound):
		log.Info("invite not found")
		c.responser.RenderErr(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.Info("account has no rights to accept this invite")
		c.responser.RenderErr(w, problems.Forbidden("account has no rights to accept this invite"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.Info("invite already accepted")
		c.responser.RenderErr(w, problems.Conflict("invite already accepted"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.Info("invite has expired")
		c.responser.RenderErr(w, problems.Forbidden("invite has expired"))
	case err != nil:
		log.WithError(err).Error("failed to accept invite")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Invite(res))
	}
}
