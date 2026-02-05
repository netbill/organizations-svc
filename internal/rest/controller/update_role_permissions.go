package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateRolePermissions(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update role permissions request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	permissions := models.OrgRolePermissionAccess{}
	for _, p := range req.Data.Attributes.Permissions {
		switch p.Code {
		case models.RolePermissionOrganizationUpdate:
			permissions.OrganizationUpdate = true
		case models.RolePermissionInvitesManage:
			permissions.InvitesManage = true
		case models.RolePermissionMembersDelete:
			permissions.MembersDelete = true
		case models.RolePermissionMembersUpdate:
			permissions.MembersUpdate = true
		case models.RolePermissionRolesManage:
			permissions.RolesManage = true
		default:
			c.log.Errorf("invalid permission code: %s", p)
			c.responser.RenderErr(w, problems.BadRequest(
				fmt.Errorf("invalid permission code: %s", p),
			)...)
			return
		}
	}

	role, perm, err := c.core.permissions.SetForRole(
		r.Context(),
		initiator,
		req.Data.Id,
		permissions,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update role permissions")
		switch {
		case errors.Is(err, errx.ErrorRoleNotFound):
			c.responser.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorCannotUpdatePermissionsHeadRole):
			c.responser.RenderErr(w, problems.Forbidden("cannot update permissions of head role"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update role permissions"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Role(role, &perm))
}
