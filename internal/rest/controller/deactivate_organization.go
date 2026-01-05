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
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) DeactivateOrganization(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	res, err := c.core.DeactivateOrganization(
		r.Context(),
		initiator.ID,
		organizationID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to deactivate organization")
		switch {
		case errors.Is(err, errx.ErrorOrganizationIsSuspended):
			ape.RenderErr(w, problems.Forbidden("organization is suspended"))
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			ape.RenderErr(w, problems.NotFound("organization not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Organization(res))
}
