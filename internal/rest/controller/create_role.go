package controller

import (
	"errors"
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) CreateRole(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateRole(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create role request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	res, err := s.core.CreateRoleByUser(r.Context(), initiator.ID, role.CreateParams{
		AgglomerationID: req.Data.Attributes.AgglomerationId,
		Name:            req.Data.Attributes.Name,
		Rank:            req.Data.Attributes.Rank,
		Description:     req.Data.Attributes.Description,
		Color:           req.Data.Attributes.Color,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to create role")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to create role"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Role(res, nil))
}
