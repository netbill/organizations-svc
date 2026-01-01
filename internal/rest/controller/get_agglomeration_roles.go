package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/pagi"
)

func (s Service) GetAgglomerationRoles(w http.ResponseWriter, r *http.Request) {
	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomeration_id"))
	if err != nil {
		s.log.Errorf("failed to parse agglomeration id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration id"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit == 0 || limit > 100 {
		s.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	roles, err := s.core.GetRoles(r.Context(), role.FilterParams{
		AgglomerationID: &agglomerationID,
	}, limit, offset)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get agglomeration roles")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Roles(r, roles))
}
