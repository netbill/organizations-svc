package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/responses"
)

func (c Controller) UpdateMember(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateMember(r)
	if err != nil {
		c.log.Errorf("invalid update member request, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	memberId, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.core.UpdateMember(r.Context(), initiator.ID, memberId, member.UpdateParams{
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Member(res))
}
