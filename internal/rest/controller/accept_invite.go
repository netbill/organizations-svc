package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid invite id")
		http.Error(w, "invalid invite id", http.StatusBadRequest)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := s.core.AcceptInvite(r.Context(), initiator.ID, inviteID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to accept invite")
		switch {
		case errors.Is(err, errx.ErrorInviteNotFound):
			ape.RenderErr(w, problems.NotFound("invite not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to accept invite"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, res)
}
