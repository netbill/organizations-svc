package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateOrganization = "update_organization"

func (c *Controller) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateOrganization)

	req, err := request.UpdateOrganization(r)
	if err != nil {
		log.WithError(err).Info("invalid update organization request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
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
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		c.responser.RenderErr(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update organization")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIconKeyIsInvalid):
		log.WithError(err).Info("icon key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconFormatIsNotAllowed):
		log.WithError(err).Info("icon format is not allowed")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconContentIsExceedsMax):
		log.WithError(err).Info("icon content is exceeds max")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationIconResolutionIsInvalid):
		log.WithError(err).Info("icon resolution is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerKeyIsInvalid):
		log.WithError(err).Info("banner key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerFormatIsNotAllowed):
		log.WithError(err).Info("banner format is not allowed")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerContentIsExceedsMax):
		log.WithError(err).Info("banner content is exceeds max")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerResolutionIsInvalid):
		log.WithError(err).Info("banner resolution is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Organization(res))
	}
}
