package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
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
