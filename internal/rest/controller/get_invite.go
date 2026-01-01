package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) GetInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid invite id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid invite id"))...)
		return
	}

	invite, err := s.core.GetInvite(r.Context(), inviteID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Invite(invite))
}
