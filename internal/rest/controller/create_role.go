package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CreateRole(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateRole(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create role request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := c.core.role.Create(r.Context(), initiator,
		role.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			Name:           req.Data.Attributes.Name,
			Rank:           req.Data.Attributes.Rank,
			Description:    req.Data.Attributes.Description,
			Color:          req.Data.Attributes.Color,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create role")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to create role"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusCreated, responses.Role(res, nil))
}
