package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationLeaveOrganization = "leave_organization"

func (c *Controller) LeaveOrganization(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationLeaveOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	err = c.modules.Member.DeleteSelf(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("account is not a member of organization")
		render.ResponseError(w, problems.Forbidden("account is not a member of organization"))
	case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
		log.WithError(err).Warn("cannot leave organization as head member")
		render.ResponseError(w, problems.Forbidden("cannot leave organization as head member"))
	case err != nil:
		log.WithError(err).Error("failed to leave organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("left organization")
		render.Response(w, http.StatusNoContent, nil)
	}
}
