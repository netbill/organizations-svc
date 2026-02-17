package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteUploadOrganizationIcon = "delete_upload_organization_icon"

func (c *Controller) DeleteOrganizationUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationIcon)

	req, err := request.DeleteUploadOrgIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload organization icon request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

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
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization does not exist")
		c.responser.RenderErr(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update organization")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization icon in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
