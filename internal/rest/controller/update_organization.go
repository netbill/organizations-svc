package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateOrganization(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update organization request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	uploadFilesData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	res, err := c.core.organization.UpdateWithSession(
		r.Context(),
		initiator,
		req.Data.Id,
		organization.UpdateParams{
			Name: req.Data.Attributes.Name,
			Media: organization.UpdateMediaParams{
				UploadSessionID: uploadFilesData.GetUploadSessionID(),
				DeletedBanner:   req.Data.Attributes.DeleteBanner,
				DeletedIcon:     req.Data.Attributes.DeleteIcon,
			},
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update organization")
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			c.responser.RenderErr(w, problems.NotFound("organization not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
		case errors.Is(err, errx.ErrorOrganizationIconTooLarge),
			errors.Is(err, errx.ErrorOrganizationIconContentFormatNotAllowed),
			errors.Is(err, errx.ErrorOrganizationIconContentTypeNotAllowed),
			errors.Is(err, errx.ErrorOrganizationIconResolutionNotAllowed):
			c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
				"icon": fmt.Errorf(err.Error()),
			})...)
		case errors.Is(err, errx.ErrorOrganizationBannerTooLarge),
			errors.Is(err, errx.ErrorOrganizationBannerContentFormatNotAllowed),
			errors.Is(err, errx.ErrorOrganizationBannerContentTypeNotAllowed),
			errors.Is(err, errx.ErrorOrganizationBannerResolutionNotAllowed):
			c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
				"banner": fmt.Errorf(err.Error()),
			})...)
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Organization(res))
}
