package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateMember = "update_member"

func (c *Controller) UpdateMember(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMember)

	req, err := request.UpdateMember(r)
	if err != nil {
		log.WithError(err).Info("invalid update member request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Info("invalid member id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.core.member.Update(
		r.Context(),
		scope.AccountActor(r),
		memberID,
		member.UpdateParams{
			Position: req.Data.Attributes.Position,
			Label:    req.Data.Attributes.Label,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
		c.responser.RenderErr(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update member")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update member"))
	case err != nil:
		log.WithError(err).Error("failed to update member")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Member(res))
	}
}
