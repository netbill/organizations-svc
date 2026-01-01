package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/member"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) UpdateMember(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateMember(r)
	if err != nil {
		s.log.Errorf("invalid update member request, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiatorID, err := uuid.Parse(r.URL.Query().Get("initiator_id"))
	if err != nil {
		s.log.Errorf("failed to parse initiator id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid initiator id"))...)
		return
	}

	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		s.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := s.core.UpdateMember(r.Context(), initiatorID, memberId, member.UpdateParams{
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to update member")
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
