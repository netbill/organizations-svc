package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCreateInvite = "create_invite"

func (c *Controller) CreateInvite(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateInvite)

	req, err := request.SentInvite(r)
	if err != nil {
		log.WithError(err).Info("invalid create invite request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId).WithField("account_id", req.Data.Attributes.AccountId)

	inv, err := c.core.invite.Create(
		r.Context(),
		scope.AccountActor(r),
		invite.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			AccountID:      req.Data.Attributes.AccountId,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		},
	)
	switch {
	case errors.Is(err, errx.ErrorAccountAlreadyMember):
		log.Info("account is already a member of the organization")
		c.responser.RenderErr(w, problems.Conflict("account is already a member of the organization"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create invite")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to create invite"))
	case err != nil:
		log.WithError(err).Error("failed to create invite")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusCreated, responses.Invite(inv))
	}
}
