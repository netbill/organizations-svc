package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdateMember(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateMember(r)
	if err != nil {
		c.log.Errorf("invalid update member request, cause %s", err)
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	memberId, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.core.member.Update(
		r.Context(),
		initiator,
		memberId,
		member.UpdateParams{
			Position: req.Data.Attributes.Position,
			Label:    req.Data.Attributes.Label,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			c.responser.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update member"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Member(res))
}
