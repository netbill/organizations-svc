package controller

import (
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateOrganization(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create organization request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.CreateOrganization(
		r.Context(),
		initiator.ID,
		organization.CreateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create organization")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Organization(res))
}
