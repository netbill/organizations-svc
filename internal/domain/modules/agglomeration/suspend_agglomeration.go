package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	agglo, err := s.repo.UpdateAgglomerationStatus(ctx, ID, entity.AgglomerationStatusSuspended)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to suspend agglomeration: %w", err))
	}

	err = s.messager.WriteAgglomerationSuspended(ctx, agglo)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish agglomeration suspend event: %w", err),
		)
	}

	return agglo, nil
}
