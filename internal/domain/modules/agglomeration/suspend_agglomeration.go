package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (agglo models.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomerationStatus(ctx, ID, models.AgglomerationStatusSuspended)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to suspend agglomeration: %w", err))
		}

		err = s.messenger.WriteAgglomerationSuspended(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration suspend event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Agglomeration{}, err
	}

	return agglo, nil
}
