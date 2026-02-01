package controller

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := c.core.GetAllPermissions(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get all permissions")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	c.responser.Render(w, http.StatusOK, responses.RolePermissions(perms))
}
