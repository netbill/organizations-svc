package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateRole = "update_role"

func (c *Controller) UpdateRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateRole)

	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		log.WithError(err).Info("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	req, err := request.UpdateRole(r)
	if err != nil {
		log.WithError(err).Info("invalid update role request")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid update role request"))...)
		return
	}

	log = log.WithField("role_id", roleID)

	res, err := c.modules.Role.Update(
		r.Context(),
		scope.AccountActor(r),
		roleID,
		role.UpdateParams{
			Name:        req.Data.Attributes.Name,
			Description: req.Data.Attributes.Description,
			Color:       req.Data.Attributes.Color,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorRoleNotFound):
		log.Info("role not found")
		c.responser.RenderErr(w, problems.NotFound("role not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update role")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update role"))
	case err != nil:
		log.WithError(err).Error("failed to update role")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Role(res, nil))
	}
}
