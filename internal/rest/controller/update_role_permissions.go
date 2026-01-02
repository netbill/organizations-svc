package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateRolePermissions(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update role permissions request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	dict := make(map[string]bool)
	for _, item := range req.Data.Attributes.Roles {
		dict[item.Code] = item.Status
	}

	role, perm, err := c.core.SetRolePermissions(r.Context(), initiator.ID, req.Data.Id, dict)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update role permissions")
		switch {
		case errors.Is(err, errx.ErrorRoleNotFound):
			ape.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update role permissions"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Role(role, perm))
}
