package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) DeactivateAgglomeration(
	ctx context.Context,
	accountID,
	agglomerationID uuid.UUID,
) (models.Agglomeration, error) {
	agglo, err := s.GetAgglomeration(ctx, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, err
	}

	if agglo.Status == models.AgglomerationStatusSuspended {
		return models.Agglomeration{}, errx.ErrorAgglomerationIsSuspended.Raise(
			fmt.Errorf("agglomeration is not suspended"),
		)
	}

	err = s.checkPermissionForManageAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomerationStatus(ctx, agglomerationID, models.AgglomerationStatusInactive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to deactivate agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationDeactivated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration deactivate event: %w", err))
		}

		return nil
	}); err != nil {
		return models.Agglomeration{}, err
	}

	return agglo, nil
}
