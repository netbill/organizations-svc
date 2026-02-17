package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteUploadOrganizationBanner = "delete_upload_organization_banner"

func (c *Controller) DeleteOrganizationUploadBanner(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationBanner)

	req, err := request.DeleteUploadOrgBanner(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload organization banner request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

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
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization does not exist")
		c.responser.RenderErr(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update organization")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization banner in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
