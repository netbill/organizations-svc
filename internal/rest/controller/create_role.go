package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCreateRole = "create_role"

func (c *Controller) CreateRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateRole)

	req, err := request.CreateRole(r)
	if err != nil {
		log.WithError(err).Info("invalid create role request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId)

	res, err := c.modules.Role.Create(
		r.Context(),
		scope.AccountActor(r),
		role.CreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			Name:           req.Data.Attributes.Name,
			Rank:           req.Data.Attributes.Rank,
			Description:    req.Data.Attributes.Description,
			Color:          req.Data.Attributes.Color,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create role")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to create role"))
	case err != nil:
		log.WithError(err).Error("failed to create role")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusCreated, responses.Role(res, nil))
	}
}
