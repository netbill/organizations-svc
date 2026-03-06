package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationSuspendOrganization = "suspend_organization"

func (c *OrganizationController) Suspend(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationSuspendOrganization)

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

	log = log.WithField("organization_id", organizationID)

	org, err := c.organizations.Suspend(r.Context(), organizationID, true)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to suspend organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization suspended")
		render.Response(w, http.StatusOK, responses.Organization(r, org))
	}
}

const operationUnsuspendOrganization = "unsuspend_organization"

func (c *OrganizationController) Unsuspend(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUnsuspendOrganization)

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

	log = log.WithField("organization_id", organizationID)

	org, err := c.organizations.Suspend(r.Context(), organizationID, false)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to unsuspend organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization unsuspended")
		render.Response(w, http.StatusOK, responses.Organization(r, org))
	}
}

const operationActivateOrganization = "activate_organization"

func (c *OrganizationController) Activate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationActivateOrganization)

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

	log = log.WithField("organization_id", organizationID)

	res, err := c.organizations.Activate(r.Context(), scope.AccountActor(r), organizationID, true)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can activate organization")
		render.ResponseError(w, problems.Forbidden("only organization head can activate organization"))
	case err != nil:
		log.WithError(err).Error("failed to activate organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization activated")
		render.Response(w, http.StatusOK, responses.Organization(r, res))
	}
}

const operationDeactivateOrganization = "deactivate_organization"

func (c *OrganizationController) Deactivate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeactivateOrganization)

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

	log = log.WithField("organization_id", organizationID)

	res, err := c.organizations.Activate(r.Context(), scope.AccountActor(r), organizationID, false)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to deactivate organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization deactivated")
		render.Response(w, http.StatusOK, responses.Organization(r, res))
	}
}
