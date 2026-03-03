package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

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
