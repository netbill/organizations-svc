package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteOrganization = "delete_organization"

func (c *Controller) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteOrganization)

	orgID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Info("invalid organization id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	log = log.WithField("organization_id", orgID)

	err = c.modules.Organization.Delete(r.Context(), scope.AccountActor(r), orgID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		c.responser.RenderErr(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete organization")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
