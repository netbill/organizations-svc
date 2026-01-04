package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/pagi"
)

func (c Controller) GetOrganizationRoles(w http.ResponseWriter, r *http.Request) {
	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.Errorf("failed to parse organization id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	roles, err := c.core.GetRoles(r.Context(), role.FilterParams{
		OrganizationID: &organizationID,
	}, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organization roles")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Roles(r, roles))
}
