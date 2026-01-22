package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/pagi"
)

func (c Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := organization.FilterParams{}

	if _, ok := q["name"]; ok {
		name := strings.TrimSpace(q.Get("name"))
		params.Name = &name
	}

	if _, ok := q["status"]; ok {
		status := strings.TrimSpace(q.Get("status"))
		params.Status = &status
	}

	limit, offset := pagi.GetPagination(r)
	if limit < 1 || limit > 100 {
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	organizations, err := c.core.GetOrganizations(r.Context(), params, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organizations")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Organizations(r, organizations))
}
