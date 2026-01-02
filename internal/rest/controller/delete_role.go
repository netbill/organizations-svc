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

func (c Controller) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid role id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	if err = c.core.DeleteRole(r.Context(), initiator.ID, roleID); err != nil {
		c.log.WithError(err).Errorf("failed to delete role")
		switch {
		case errors.Is(err, errx.ErrorRoleNotFound):
			ape.RenderErr(w, problems.NotFound("role not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to delete role"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
