package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeleteUploadOrganizationBanner(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid organization id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)

		return
	}

	uploadContentData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	err = c.core.DeleteUpdateOrganizationBannerInSession(
		r.Context(),
		initiator.GetAccountID(),
		organizationID,
		uploadContentData.GetUploadSessionID(),
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete organization banner in upload session")
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			c.responser.RenderErr(w, problems.NotFound("organization does not exist"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}
}
