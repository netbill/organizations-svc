package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetOrganization = "get_organization"

func (c *Controller) GetOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Info("invalid organization id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	org, err := c.modules.Organization.GetByID(r.Context(), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to get organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organization(org))
	}
}
