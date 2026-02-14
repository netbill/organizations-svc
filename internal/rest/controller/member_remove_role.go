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

const operationMemberRemoveRole = "member_remove_role"

func (c *Controller) MemberRemoveRole(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationMemberRemoveRole)

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

	err = c.core.role.RemoveFromMember(
		r.Context(),
		scope.AccountActor(r),
		memberID,
		roleID,
	)
	switch {
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
		c.responser.RenderErr(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorRoleNotFound):
		log.Info("role not found")
		c.responser.RenderErr(w, problems.NotFound("role not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to remove role from member")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to remove role from member"))
	case errors.Is(err, errx.ErrorCannotRemoveHeadRoleFromMember):
		log.Info("cannot remove head role from member")
		c.responser.RenderErr(w, problems.Forbidden("cannot remove head role from member"))
	case err != nil:
		log.WithError(err).Error("failed to remove role from member")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
