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

func (c *Controller) MemberRemoveRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid role id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id: %s", chi.URLParam(r, "role_id")))...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid member id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")))...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	if err = c.core.role.RemoveFromMember(
		r.Context(),
		initiator,
		memberID,
		roleID,
	); err != nil {
		c.log.WithError(err).Errorf("failed to remove role to member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			c.responser.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorRoleNotFound):
			c.responser.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to remove role to member"))
		case errors.Is(err, errx.ErrorCannotRemoveHeadRoleFromMember):
			c.responser.RenderErr(w, problems.Forbidden("cannot remove head role from member"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}
}
