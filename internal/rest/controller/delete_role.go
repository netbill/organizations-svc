package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeleteRole(w http.ResponseWriter, r *http.Request) {
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

	if err = c.core.DeleteRole(r.Context(), initiator.GetAccountID(), roleID); err != nil {
		c.log.WithError(err).Errorf("failed to delete role")
		switch {
		case errors.Is(err, errx.ErrorRoleNotFound):
			c.responser.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorCannotDeleteHeadRole):
			c.responser.RenderErr(w, problems.Forbidden("cannot delete head role"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete role"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
