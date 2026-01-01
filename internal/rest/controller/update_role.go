package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid role id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateRole(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid update role request")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid update role request"))...)
		return
	}

	res, err := s.core.UpdateRole(r.Context(), initiator.ID, roleID, role.UpdateParams{
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Color:       req.Data.Attributes.Color,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to update role")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Role(res, nil))
}
