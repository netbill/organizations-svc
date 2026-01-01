package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) SuspendAgglomeration(w http.ResponseWriter, r *http.Request) {
	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomerationID"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid agglomeration ID")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration ID"))...)
		return
	}

	agglomeration, err := s.core.SuspendAgglomeration(r.Context(), agglomerationID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to suspend agglomeration")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Agglomeration(agglomeration))
}
