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

const operationDeleteUploadOrganizationIcon = "delete_upload_organization_icon"

func (c *Controller) DeleteOrganizationUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationIcon)

	req, err := requests.DeleteUploadOrgIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload organization icon requests")
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
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
