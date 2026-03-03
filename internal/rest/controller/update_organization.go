package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdateOrganization = "update_organization"

func (c *Controller) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateOrganization)

	req, err := requests.UpdateOrganization(r)
	if err != nil {
		log.WithError(err).Warn("invalid update organization requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	res, err := c.modules.Organization.Update(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		organization.UpdateParams{
			Name:      req.Data.Attributes.Name,
			IconKey:   req.Data.Attributes.IconKey,
			BannerKey: req.Data.Attributes.BannerKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorOrganizationIconKeyIsInvalid):
		log.WithError(err).Warn("icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconFormatIsNotAllowed):
		log.WithError(err).Warn("icon format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconContentIsExceedsMax):
		log.WithError(err).Warn("icon content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconResolutionIsInvalid):
		log.WithError(err).Warn("icon resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerKeyIsInvalid):
		log.WithError(err).Warn("banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerFormatIsNotAllowed):
		log.WithError(err).Warn("banner format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerContentIsExceedsMax):
		log.WithError(err).Warn("banner content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerResolutionIsInvalid):
		log.WithError(err).Warn("banner resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization updated successfully")
		render.Response(w, http.StatusOK, responses.Organization(res))
	}
}
