package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateOrganization(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update organization request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.UpdateOrganization(
		r.Context(),
		initiator.ID,
		req.Data.Id,
		organization.UpdateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update organization")
		switch {
		case errors.Is(err, errx.ErrorOrganizationIsSuspended):
			ape.RenderErr(w, problems.Forbidden("organization is suspended"))
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			ape.RenderErr(w, problems.NotFound("organization not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Organization(res))
}
