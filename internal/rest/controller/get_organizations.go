package controller

import (
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetOrganizations = "get_organizations"

func (c *Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizations)

	q := r.URL.Query()
	params := organization.FilterParams{}

	if v := strings.TrimSpace(q.Get("text")); v != "" {
		params.Text = &v
	}

	if v := strings.TrimSpace(q.Get("status")); v != "" {
		params.Status = &v
	}

	limit, offset := pagi.GetPagination(r)
	if limit < 1 || limit > 100 {
		log.Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit must be between 1 and 100"),
		})...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	organizations, err := c.modules.Organization.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organizations")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organizations(r, organizations))
	}
}
