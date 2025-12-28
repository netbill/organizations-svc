package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (agglo entity.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomerationStatus(ctx, ID, entity.AgglomerationStatusInactive)
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
		return entity.Agglomeration{}, err
	}

	return agglo, nil
}

func (s Service) DeactivateAgglomerationByUser(
	ctx context.Context,
	memberID, agglomerationID uuid.UUID,
) (entity.Agglomeration, error) {
	agglo, err := s.GetAgglomeration(ctx, agglomerationID)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	if agglo.Status == entity.AgglomerationStatusSuspended {
		return entity.Agglomeration{}, errx.ErrorAgglomerationIsSuspended.Raise(
			fmt.Errorf("agglomeration is not suspended"),
		)
	}

	err = s.checkPermissionForManageAgglomeration(
		ctx,
		memberID,
		agglomerationID,
	)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return s.DeactivateAgglomeration(ctx, agglomerationID)
}
