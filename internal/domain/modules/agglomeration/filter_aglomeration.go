package agglomeration

import (
	"context"
	"fmt"

	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
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
