package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

const operationGetOrganizations = "get_organizations"

func (c *Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizations)

	q := r.URL.Query()
	params := organization.FilterParams{}

	if v := strings.TrimSpace(q.Get("name")); v != "" {
		params.Name = &v
	}

	if v := strings.TrimSpace(q.Get("status")); v != "" {
		params.Status = &v
	}

	limit, offset := pagi.GetPagination(r)
	if limit < 1 || limit > 100 {
		log.Info("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	organizations, err := c.core.organization.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organizations")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Organizations(r, organizations))
	}
}
