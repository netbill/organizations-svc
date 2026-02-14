package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
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

	permissions := models.OrgRolePermissionEnable{}
	for _, p := range req.Data.Attributes.Permissions {
		switch p.Code {
		case models.RolePermissionOrgUpdate:
			permissions.OrgUpdate = true
		case models.RolePermissionInvitesManage:
			permissions.InvitesManage = true
		case models.RolePermissionMembersDelete:
			permissions.MembersDelete = true
		case models.RolePermissionMembersUpdate:
			permissions.MembersUpdate = true
		case models.RolePermissionRolesManage:
			permissions.RolesManage = true
		case models.RolePermissionPlaceCreate:
			permissions.PlaceCreate = true
		case models.RolePermissionPlaceDelete:
			permissions.PlaceDelete = true
		case models.RolePermissionPlaceUpdate:
			permissions.PlaceUpdate = true
		default:
			log.Info("invalid permission code")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid permission code: %s", p.Code))...)
			return
		}
	}

	role, perm, err := c.modules.Permissions.SetForRole(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		permissions,
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
