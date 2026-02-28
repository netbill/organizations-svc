package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeleteMember = "delete_member"

func (c *Controller) DeleteMember(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Info("invalid member id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	log = log.WithField("member_id", memberID)

	err = c.modules.Member.Delete(r.Context(), scope.AccountActor(r), memberID)
	switch {
	case errors.Is(err, errx.ErrorCannotDeleteSelf):
		log.Info("cannot delete self")
		render.ResponseError(w, problems.Forbidden("cannot delete self"))
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorMemberDeleted):
		log.Info("member already deleted")
		render.Response(w, http.StatusNoContent, nil)
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.Info("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete member")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete member"))
	case errors.Is(err, errx.ErrorOrganizationNotFound):
		log.Info("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
		log.Info("cannot delete organization head member")
		render.ResponseError(w, problems.Forbidden("cannot delete organization head member"))
	case err != nil:
		log.WithError(err).Error("failed to delete member")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("member deleted")
		render.Response(w, http.StatusNoContent, nil)
	}
}
