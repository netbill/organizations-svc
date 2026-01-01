package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/rest"
	"github.com/umisto/cities-svc/internal/rest/request"
	"github.com/umisto/cities-svc/internal/rest/responses"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type Domain interface {
	CreateAgglomeration(
		ctx context.Context,
		accountID uuid.UUID,
		params agglomeration.CreateParams,
	) (models.Agglomeration, error)

	GetAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		params agglomeration.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)

	UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params agglomeration.UpdateParams) (models.Agglomeration, error)
	UpdateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		params agglomeration.UpdateParams,
	) (models.Agglomeration, error)

	ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	ActivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)

	DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	DeactivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)

	SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error
}

type AggloController struct {
	domain Domain
	log    logium.Logger
}

func (c AggloController) CreateAgglomeration(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateAgglomeration(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create agglomeration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.domain.CreateAgglomeration(
		r.Context(),
		req.Data.Attributes.Head,
		agglomeration.CreateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create agglomeration")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Agglomeration(res))
}

func (c AggloController) GetAgglomeration(w http.ResponseWriter, r *http.Request) {
	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomerationID"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid agglomeration ID")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration ID"))...)
		return
	}

	agglo, err := c.domain.GetAgglomeration(r.Context(), agglomerationID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get agglomeration")
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

func (c AggloController) ActivateAgglomeration(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomeration_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid agglomeration id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration id"))...)
		return
	}

	res, err := c.domain.ActivateAgglomerationByUser(
		r.Context(),
		initiator.ID,
		agglomerationID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to activate agglomeration")
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

func (c AggloController) DeactivateAgglomeration(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	agglomerationID, err := uuid.Parse(chi.URLParam(r, "agglomeration_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid agglomeration id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid agglomeration id"))...)
		return
	}

	res, err := c.domain.DeactivateAgglomerationByUser(
		r.Context(),
		initiator.ID,
		agglomerationID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to deactivate agglomeration")
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

func (c AggloController) UpdateAgglomeration(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := request.UpdateAgglomeration(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update agglomeration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.domain.UpdateAgglomerationByUser(
		r.Context(),
		initiator.ID,
		req.Data.Id,
		agglomeration.UpdateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update agglomeration")
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
