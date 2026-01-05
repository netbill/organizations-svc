package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest"
)

func (c Controller) MemberRemoveRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid role id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id: %s", chi.URLParam(r, "role_id")))...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid member id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")))...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	err = c.core.RemoveMemberRole(r.Context(), initiator.ID, memberID, roleID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to remove role to member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorRoleNotFound):
			ape.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to remove role to member"))
		case errors.Is(err, errx.ErrorCannotRemoveHeadRoleFromMember):
			ape.RenderErr(w, problems.Forbidden("cannot remove head role from member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}
}
