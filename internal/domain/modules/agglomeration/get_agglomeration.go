package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) GetAgglomeration(ctx context.Context, agglomerationID uuid.UUID) (entity.Agglomeration, error) {
	res, err := s.repo.GetAgglomerationByID(ctx, agglomerationID)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get agglomeration by id: %w", err),
		)
	}
	if res.IsNil() {
		return entity.Agglomeration{}, errx.ErrorAgglomerationNotFound.Raise(
			fmt.Errorf("agglomeration with id %s not found", agglomerationID),
		)
	}

	return res, nil
}
