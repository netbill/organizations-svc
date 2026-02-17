package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCreateOrganizationUploadMediaLink = "create_organization_upload_media_link"

func (c *Controller) CreateOrganizationUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganizationUploadMediaLink)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Info("invalid organization id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"organization_id": err,
		})...)
		return
	}

	org, media, err := c.modules.Organization.CreateOrgUploadMediaLinks(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization does not exist")
		c.responser.RenderErr(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create organization upload media link")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to create organization upload media link"))
	case err != nil:
		log.WithError(err).Error("failed to create organization upload media link")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.UploadOrganizationMediaLinks(org, media))
	}
}
