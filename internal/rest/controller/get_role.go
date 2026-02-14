package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetRole = "get_role"

func (c *Controller) GetRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetRole)

	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		log.WithError(err).Info("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	log = log.WithField("role_id", roleID)

	role, perm, err := c.core.role.GetWithPermissions(r.Context(), scope.AccountActor(r), roleID)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get role")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Role(role, &perm))
	}
}
