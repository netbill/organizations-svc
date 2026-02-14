package controller

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCreateOrganization = "create_organization"

func (c *Controller) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganization)

	req, err := request.CreateOrganization(r)
	if err != nil {
		log.WithError(err).Info("invalid create organization request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.organization.Create(
		r.Context(),
		scope.AccountActor(r),
		organization.CreateParams{
			Name: req.Data.Attributes.Name,
		},
	)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to create organization")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusCreated, responses.Organization(res))
	}
}
