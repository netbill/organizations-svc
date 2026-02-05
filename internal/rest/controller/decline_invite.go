package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid invite id")
		http.Error(w, "invalid invite id", http.StatusBadRequest)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.invite.Decline(r.Context(), initiator, inviteID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to decline invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			c.responser.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteNotForInitiator):
			c.responser.RenderErr(w, problems.Forbidden("invite not for this account"))
		case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
			c.responser.RenderErr(w, problems.Conflict("invite already answered"))
		case errors.Is(err, errx.ErrorInviteExpired):
			c.responser.RenderErr(w, problems.Forbidden("invite has expired"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Invite(res))
}
