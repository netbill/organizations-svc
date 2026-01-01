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

func (s Service) GetAgglomeration(w http.ResponseWriter, r *http.Request) {
	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomerationID"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid agglomeration ID")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration ID"))...)
		return
	}

	agglo, err := s.core.GetAgglomeration(r.Context(), agglomerationID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get agglomeration")
		switch {
		case errors.Is(err, errx.ErrorAgglomerationNotFound):
			ape.RenderErr(w, problems.NotFound("agglomeration not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Agglomeration(agglo))
}
