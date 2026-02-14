package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationConfirmUpdateOrganization = "confirm_update_organization"

func (c *Controller) ConfirmUpdateOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationConfirmUpdateOrganization)

	req, err := request.UpdateOrganization(r)
	if err != nil {
		log.WithError(err).Info("invalid update organization request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	res, err := c.modules.Organization.UpdateWithSession(
		r.Context(),
		scope.AccountActor(r),
		scope.UploadScope(r),
		req.Data.Id,
		organization.UpdateParams{
			Name: req.Data.Attributes.Name,
			Media: organization.UpdateMediaParams{
				DeletedBanner: req.Data.Attributes.DeleteBanner,
				DeletedIcon:   req.Data.Attributes.DeleteIcon,
			},
		},
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		c.responser.RenderErr(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update organization")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIconInvalid):
		log.Info("organization icon is not valid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": fmt.Errorf(err.Error()),
		})...)
	case errors.Is(err, errx.ErrorOrganizationBannerInvalid):
		log.Info("organization banner is not valid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": fmt.Errorf(err.Error()),
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Organization(res))
	}
}
