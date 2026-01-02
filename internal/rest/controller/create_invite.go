package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) CreateInvite(w http.ResponseWriter, r *http.Request) {
	req, err := request.SentInvite(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	inv, err := c.core.CreateInvite(r.Context(),
		initiator.ID,
		invite.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			AccountID:      req.Data.Attributes.AccountId,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create invite")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to create invite"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
