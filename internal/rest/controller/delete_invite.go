package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeleteInvite = "delete_invite"

func (c *Controller) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Info("invalid invite id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	err = c.modules.Invite.Delete(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotFound):
		log.Info("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete invite")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete invite"))
	case err != nil:
		log.WithError(err).Error("failed to delete invite")
		render.ResponseError(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
