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

func (c Controller) DeleteMember(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	memberId, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	err = c.core.DeleteMember(r.Context(), initiator.ID, memberId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to delete member"))
		case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
			ape.RenderErr(w, problems.Forbidden("cannot delete organization head member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
