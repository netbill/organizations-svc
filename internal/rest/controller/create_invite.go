package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreateInvite = "create_invite"

func (c *Controller) CreateInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateInvite)

	req, err := requests.SentInvite(r)
	if err != nil {
		log.WithError(err).Warn("invalid create invite requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId).WithField("account_id", req.Data.Attributes.AccountId)

	inv, err := c.modules.Invite.Create(
		r.Context(),
		scope.AccountActor(r),
		invite.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			AccountID:      req.Data.Attributes.AccountId,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		},
	)
	switch {
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to create invite")
		render.ResponseError(w, problems.Forbidden("not enough rights to create invite"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for account not found")
		render.ResponseError(w, problems.NotFound("profile for account not found"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorAccountAlreadyMember):
		log.WithError(err).Warn("account is already a member of the organization")
		render.ResponseError(w, problems.Conflict("account is already a member of the organization"))
	case errors.Is(err, errx.ErrorActiveInviteAlreadyExists):
		log.WithError(err).Warn("active invite for this account and organization already exists")
		render.ResponseError(w, problems.Conflict("active invite for this account and organization already exists"))
	case err != nil:
		log.WithError(err).Error("failed to create invite")
		render.ResponseError(w, problems.InternalError())
	default:
		log.WithField("invite_id", inv.ID).Info("invite created successfully")
		render.Response(w, http.StatusCreated, responses.Invite(inv))
	}
}
