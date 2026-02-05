package controller

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateOrganization(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create organization request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.organization.Create(r.Context(), initiator,
		organization.CreateParams{
			Name: req.Data.Attributes.Name,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create organization")
		switch {
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusCreated, responses.Organization(res))
}
