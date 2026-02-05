package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeleteInvite(w http.ResponseWriter, r *http.Request) {
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

	if err = c.core.invite.Delete(r.Context(), initiator, inviteID); err != nil {
		c.log.WithError(err).Errorf("failed to delete invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			c.responser.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete invite"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
