package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) GetOrganization(w http.ResponseWriter, r *http.Request) {
	organizationID, err := uuid.Parse(chi.URLParam(r, "organizationID"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization ID")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization ID"))...)
		return
	}

	agglo, err := c.core.GetOrganization(r.Context(), organizationID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organization")
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			ape.RenderErr(w, problems.NotFound("organization not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Organization(agglo))
}
