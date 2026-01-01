package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	params FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.Agglomeration], error) {
	res, err := s.repo.FilterAgglomerations(ctx, params, offset, limit)
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
	res, err := s.repo.GetAgglomerationForUser(ctx, accountID, offset, limit)
	if err != nil {
		return pagi.Page[[]models.Agglomeration]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("get agglomeration for user: %w", err),
		)
	}

	return res, nil
}
