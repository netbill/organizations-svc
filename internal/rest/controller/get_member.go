package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) GetMember(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		s.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := s.core.GetMemberByID(r.Context(), memberId)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get member by id")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, res)
}
