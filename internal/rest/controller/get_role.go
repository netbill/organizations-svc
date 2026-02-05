package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"

	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetRole(w http.ResponseWriter, r *http.Request) {
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

	role, perm, err := c.core.GetRoleWithPermissions(
		r.Context(),
		models.InitiatorData{
			AccountID: initiator.GetAccountID(),
		},
		roleID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get role")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Role(role, &perm))
}
