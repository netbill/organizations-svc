package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreateOrganizationUploadMediaLink = "create_organization_upload_media_link"

func (c *Controller) CreateOrganizationUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganizationUploadMediaLink)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	org, media, err := c.modules.Organization.CreateOrgUploadMediaLinks(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can create organization upload media link")
		render.ResponseError(w, problems.Forbidden("not enough rights to create organization upload media link"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to create organization upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadOrganizationMediaLinks(org, media))
	}
}

const operationDeleteUploadOrganizationIcon = "delete_upload_organization_icon"

func (c *Controller) DeleteOrganizationUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationIcon)

	req, err := requests.DeleteUploadOrgIcon(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload organization icon requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	err = c.modules.Organization.DeleteOrgUploadIcon(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.Forbidden("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, nil)
	}
}

const operationDeleteUploadOrganizationBanner = "delete_upload_organization_banner"

func (c *Controller) DeleteOrganizationUploadBanner(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationBanner)

	req, err := requests.DeleteUploadOrgBanner(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload organization banner requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	err = c.modules.Organization.DeleteOrgUploadBanner(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.BannerKey,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.Forbidden("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization banner in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, nil)
	}
}
