package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/middlewares"
)

func (c Controller) DeleteUploadOrganizationBanner(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)

		return
	}

	uploadFilesData, err := middlewares.UploadFilesData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		ape.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	err = c.core.DeleteUpdateOrganizationBannerInSession(
		r.Context(),
		initiator.AccountID,
		organizationID,
		uploadFilesData.UploadSessionID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete organization banner in upload session")
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			ape.RenderErr(w, problems.NotFound("organization does not exist"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}
}
