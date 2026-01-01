package controller

import (
	"errors"
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) UpdateAgglomeration(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateAgglomeration(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid update agglomeration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := s.core.UpdateAgglomerationByUser(
		r.Context(),
		initiator.ID,
		req.Data.Id,
		agglomeration.UpdateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update agglomeration")
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
