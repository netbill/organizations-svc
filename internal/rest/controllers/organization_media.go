package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/organization"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreateOrganizationUploadMediaLink = "create_organization_upload_media_link"

func (c *OrganizationController) CreateUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganizationUploadMediaLink)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/organization_id": fmt.Errorf(
				"invalid organization id: %s", chi.URLParam(r, "organization_id"),
			),
		})...)
		return
	}

	org, media, err := c.organizations.CreateUploadMediaLinks(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
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
		render.Response(w, http.StatusOK, responses.UploadOrganizationMediaLinks(r, org, media))
	}
}

const operationDeleteUploadOrganizationMedia = "delete_upload_organization_media"

func (c *OrganizationController) DeleteUploadMedia(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationMedia)

	req, err := requests.DeleteUploadOrgMedia(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload organization requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	err = c.organizations.DeleteUploadMedia(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		organization.DeleteUploaderParams{
			Icon:   req.Data.Attributes.IconKey,
			Banner: req.Data.Attributes.BannerKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
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
