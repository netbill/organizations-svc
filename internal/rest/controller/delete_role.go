package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteRole = "delete_role"

func (c *Controller) DeleteRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteRole)

	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		log.WithError(err).Info("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	log = log.WithField("role_id", roleID)

	err = c.core.role.Delete(r.Context(), scope.AccountActor(r), roleID)
	switch {
	case errors.Is(err, errx.ErrorRoleNotFound):
		log.Info("role not found")
		c.responser.RenderErr(w, problems.NotFound("role not found"))
	case errors.Is(err, errx.ErrorCannotDeleteHeadRole):
		log.Info("cannot delete head role")
		c.responser.RenderErr(w, problems.Forbidden("cannot delete head role"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete role")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete role"))
	case err != nil:
		log.WithError(err).Error("failed to delete role")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
