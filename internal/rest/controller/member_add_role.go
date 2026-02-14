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

const operationMemberAddRole = "member_add_role"

func (c *Controller) MemberAddRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationMemberAddRole)

	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		log.WithError(err).Info("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id: %s", chi.URLParam(r, "role_id")))...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Info("invalid member id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")))...)
		return
	}

	log = log.WithField("role_id", roleID).WithField("member_id", memberID)

	err = c.modules.Role.AddForMember(r.Context(), scope.AccountActor(r), memberID, roleID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
		c.responser.RenderErr(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorRoleNotFound):
		log.Info("role not found")
		c.responser.RenderErr(w, problems.NotFound("role not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to add role to member")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to add role to member"))
	case errors.Is(err, errx.ErrorCannotAddHeadRoleToMember):
		log.Info("cannot add head role to member")
		c.responser.RenderErr(w, problems.Forbidden("cannot add head role to member"))
	case err != nil:
		log.WithError(err).Error("failed to add role to member")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
