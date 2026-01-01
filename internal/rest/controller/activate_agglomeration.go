package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) ActivateAgglomeration(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomeration_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid agglomeration id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration id"))...)
		return
	}

	res, err := s.core.ActivateAgglomeration(
		r.Context(),
		initiator.ID,
		agglomerationID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to activate agglomeration")
		switch {
		case errors.Is(err, errx.ErrorAgglomerationIsSuspended):
			ape.RenderErr(w, problems.Forbidden("agglomeration is suspended"))
		case errors.Is(err, errx.ErrorAgglomerationNotFound):
			ape.RenderErr(w, problems.NotFound("agglomeration not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update agglomeration"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Agglomeration(res))
}
