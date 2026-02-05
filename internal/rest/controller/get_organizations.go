package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {
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
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	organizations, err := c.core.organization.GetList(r.Context(), params, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organizations")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Organizations(r, organizations))
}
