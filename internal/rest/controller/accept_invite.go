package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid invite id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.AcceptInvite(r.Context(), initiator.AccountID, inviteID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to accept invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorInviteNotForInitiator):
			ape.RenderErr(w, problems.Forbidden("account has no rights to accept this invite"))
		case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
			ape.RenderErr(w, problems.Conflict("invite already accepted"))
		case errors.Is(err, errx.ErrorInviteExpired):
			ape.RenderErr(w, problems.Forbidden("invite has expired"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Invite(res))
}
