package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateRolePermissions = "update_role_permissions"

func (c *Controller) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateRolePermissions)

	req, err := request.UpdateRolePermissions(r)
	if err != nil {
		log.WithError(err).Info("invalid update role permissions request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("role_id", req.Data.Id)

	params := role.SetPermissions{}
	for _, p := range req.Data.Attributes.Permissions {
		params[p.Id] = p.Status
	}

	role, perm, err := c.modules.Role.UpdatePermissions(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		params,
	)
	switch {
	case errors.Is(err, errx.ErrorRoleNotFound):
		log.Info("role not found")
		c.responser.RenderErr(w, problems.NotFound("role not found"))
	case errors.Is(err, errx.ErrorCannotUpdatePermissionsHeadRole):
		log.Info("cannot update permissions of head role")
		c.responser.RenderErr(w, problems.Forbidden("cannot update permissions of head role"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update role permissions")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update role permissions"))
	case err != nil:
		log.WithError(err).Error("failed to update role permissions")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Role(role, &perm))
	}
}
