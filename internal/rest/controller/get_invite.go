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

const operationGetInvite = "get_invite"

func (c *Controller) GetInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Info("invalid invite id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	inv, err := c.modules.Invite.GetForAccount(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to get invite")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to get invite"))
	case errors.Is(err, errx.ErrorInviteNotFound):
		log.Info("invite not found")
		c.responser.RenderErr(w, problems.NotFound("invite not found"))
	case err != nil:
		log.WithError(err).Error("failed to get invite")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Invite(inv))
	}
}
