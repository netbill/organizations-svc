package controller

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (c Controller) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := c.core.GetAllPermissions(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get all permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.RolePermissions(perms))
}
