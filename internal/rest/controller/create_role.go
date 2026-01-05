package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) CreateRole(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateRole(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create role request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.CreateRole(r.Context(), initiator.ID, role.CreateParams{
		OrganizationID: req.Data.Attributes.OrganizationId,
		Name:           req.Data.Attributes.Name,
		Rank:           req.Data.Attributes.Rank,
		Description:    req.Data.Attributes.Description,
		Color:          req.Data.Attributes.Color,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create role")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to create role"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Role(res, nil))
}
