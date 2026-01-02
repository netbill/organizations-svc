package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

func (s Service) GetAgglomeration(ctx context.Context, agglomerationID uuid.UUID) (models.Agglomeration, error) {
	res, err := s.repo.GetAgglomerationByID(ctx, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get agglomeration by id: %w", err),
		)
	}
	if res.IsNil() {
		return models.Agglomeration{}, errx.ErrorAgglomerationNotFound.Raise(
			fmt.Errorf("agglomeration with id %s not found", agglomerationID),
		)
	}

	return res, nil
}

type FilterParams struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (s Service) GetAgglomerations(
	ctx context.Context,
	params FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.Agglomeration], error) {
	res, err := s.repo.GetAgglomerations(ctx, params, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Agglomeration]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("filter agglomerations: %w", err),
		)
	}

	return res, nil
}

func (s Service) GetAgglomerationForUser(
	ctx context.Context,
	accountID uuid.UUID,
	offset, limit uint,
) (pagi.Page[[]models.Agglomeration], error) {
	res, err := s.repo.GetAgglomerationsForUser(ctx, accountID, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Agglomeration]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get agglomeration for user: %w", err),
		)
	}

	return res, nil
}
