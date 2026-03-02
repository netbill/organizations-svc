package controller

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreateOrganization = "create_organization"

func (c *Controller) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganization)

	req, err := requests.CreateOrganization(r)
	if err != nil {
		log.WithError(err).Warn("invalid create organization requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.modules.Organization.Create(
		r.Context(),
		scope.AccountActor(r),
		organization.CreateParams{
			Name: req.Data.Attributes.Name,
		},
	)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to create organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusCreated, responses.Organization(res))
	}
}
