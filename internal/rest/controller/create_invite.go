package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CreateInvite(w http.ResponseWriter, r *http.Request) {
	req, err := request.SentInvite(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create invite request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	inv, err := c.core.CreateInvite(r.Context(),
		models.InitiatorData{
			AccountID: initiator.GetAccountID(),
		},
		invite.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			AccountID:      req.Data.Attributes.AccountId,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create invite")
		switch {
		case errors.Is(err, errx.ErrorAccountAlreadyMember):
			c.responser.RenderErr(w, problems.Conflict("account is already a member of the organization"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to create invite"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusCreated, responses.Invite(inv))
}
