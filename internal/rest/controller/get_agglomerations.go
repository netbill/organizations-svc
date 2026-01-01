package controller

import (
	"fmt"
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/pagi"
)

func (s Service) GetAgglomerations(w http.ResponseWriter, r *http.Request) {
	if name := r.URL.Query().Get("name"); name != "" {
		s.log.WithError(fmt.Errorf("filter by name is not supported yet")).Errorf("filter by name is not supported yet")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("filter by name is not supported yet"))...)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		s.log.WithError(fmt.Errorf("filter by status is not supported yet")).Errorf("filter by status is not supported yet")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("filter by status is not supported yet"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit == 0 || limit > 100 {
		s.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	agglomerations, err := s.core.GetAgglomerations(r.Context(), agglomeration.FilterParams{}, limit, offset)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get agglomerations")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Agglomerations(r, agglomerations))
}
