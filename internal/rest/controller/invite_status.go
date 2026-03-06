package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationAcceptInvite = "accept_invite"

func (c *InviteController) Accept(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationAcceptInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/invite_id": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.invites.Accept(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not exists")
		render.ResponseError(w, problems.NotFound("invite not exists"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.WithError(err).Warn("account has no rights to accept this invite")
		render.ResponseError(w, problems.Forbidden("account has no rights to accept this invite"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already accepted")
		render.ResponseError(w, problems.Conflict("invite already accepted"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.WithError(err).Warn("invite has expired")
		render.ResponseError(w, problems.Forbidden("invite has expired"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to accept invite")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Invite(r, res))
	}
}

const operationDeclineInvite = "decline_invite"

func (c *InviteController) Decline(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeclineInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/invite_id": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.invites.Decline(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.WithError(err).Warn("invite not for this account")
		render.ResponseError(w, problems.Forbidden("invite not for this account"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already answered")
		render.ResponseError(w, problems.Conflict("invite already answered"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.WithError(err).Warn("invite has expired")
		render.ResponseError(w, problems.Forbidden("invite has expired"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case err != nil:
		log.WithError(err).Error("failed to decline invite")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Invite(r, res))
	}
}

const operationCancelledInvite = "delete_invite"

func (c *InviteController) Cancelled(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCancelledInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/invite_id": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	invite, err := c.invites.Cancelled(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("only organization members can cancel invite")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can cancel invite")
		render.ResponseError(w, problems.Forbidden("only organization head can cancel invite"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already answered")
		render.ResponseError(w, problems.Forbidden("invite already answered"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to delete invite")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("invite deleted")
		render.Response(w, http.StatusOK, responses.Invite(r, invite))
	}
}
