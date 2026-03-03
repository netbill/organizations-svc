package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdateMember = "update_member"

func (c *Controller) UpdateMember(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMember)

	req, err := requests.UpdateMember(r)
	if err != nil {
		log.WithError(err).Warn("invalid update member requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.modules.Member.Update(
		r.Context(),
		scope.AccountActor(r),
		memberID,
		member.UpdateParams{
			Position: req.Data.Attributes.Position,
			Label:    req.Data.Attributes.Label,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update member")
		render.ResponseError(w, problems.Forbidden("not enough rights to update member"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to update member")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Member(res))
	}
}
