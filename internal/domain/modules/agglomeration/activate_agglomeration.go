package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (agglo models.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomerationStatus(ctx, ID, models.AgglomerationStatusActive)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to activate agglomeration: %w", err))
		}

		err = s.messenger.WriteAgglomerationActivated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration activated event: %w", err))
		}

		return nil
	}); err != nil {
		return models.Agglomeration{}, err
	}

	return agglo, nil
}

func (s Service) ActivateAgglomerationByUser(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (models.Agglomeration, error) {
	agglo, err := s.GetAgglomeration(ctx, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, err
	}

	if agglo.Status == models.AgglomerationStatusSuspended {
		return models.Agglomeration{}, errx.ErrorAgglomerationIsSuspended.Raise(
			fmt.Errorf("agglomeration is suspended"),
		)
	}

	err = s.checkPermissionForManageAgglomeration(ctx, accountID, agglomerationID)
	if err != nil {
		return models.Agglomeration{}, err
	}

	return s.ActivateAgglomeration(ctx, agglomerationID)
}
