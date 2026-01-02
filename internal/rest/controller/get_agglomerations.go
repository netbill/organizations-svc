package controller

import (
	"fmt"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/domain/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/pagi"
)

func (c Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	if name := r.URL.Query().Get("name"); name != "" {
		c.log.WithError(fmt.Errorf("filter by name is not supported yet")).Errorf("filter by name is not supported yet")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("filter by name is not supported yet"))...)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		c.log.WithError(fmt.Errorf("filter by status is not supported yet")).Errorf("filter by status is not supported yet")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("filter by status is not supported yet"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit == 0 || limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	organizations, err := c.core.GetOrganizations(r.Context(), organization.FilterParams{}, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organizations")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Organizations(r, organizations))
}
