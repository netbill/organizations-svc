package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetOrganization(w http.ResponseWriter, r *http.Request) {
	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization ID")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization ID"))...)
		return
	}

	org, err := c.core.GetOrganization(r.Context(), organizationID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organization")
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			c.responser.RenderErr(w, problems.NotFound("organization not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Organization(org))
}
