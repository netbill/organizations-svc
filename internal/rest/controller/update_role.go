package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateRole(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update role request")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid update role request"))...)
		return
	}

	res, err := c.core.UpdateRole(r.Context(), initiator.GetAccountID(), roleID, role.UpdateParams{
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Color:       req.Data.Attributes.Color,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update role")
		switch {
		case errors.Is(err, errx.ErrorRoleNotFound):
			c.responser.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update role"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Role(res, nil))
}
