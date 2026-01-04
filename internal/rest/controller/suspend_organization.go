package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) SuspendOrganization(w http.ResponseWriter, r *http.Request) {
	organizationID, err := uuid.Parse(chi.URLParam(r, "organizationID"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization ID")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization ID"))...)
		return
	}

	organization, err := c.core.SuspendOrganization(r.Context(), organizationID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to suspend organization")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Organization(organization))
}
