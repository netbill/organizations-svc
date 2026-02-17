package controller

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetAllPermissions = "get_all_permissions"

func (c *Controller) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetAllPermissions)

	perms, err := c.modules.Role.GetAllPermissions(r.Context())
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get all permissions")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.RolePermissions(perms))
	}
}
