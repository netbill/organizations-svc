package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (agglo entity.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomerationStatus(ctx, ID, entity.AgglomerationStatusActive)
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
		return entity.Agglomeration{}, err
	}

	return agglo, nil
}

func (s Service) ActivateAgglomerationByUser(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (entity.Agglomeration, error) {
	agglo, err := s.GetAgglomeration(ctx, agglomerationID)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	if agglo.Status == entity.AgglomerationStatusSuspended {
		return entity.Agglomeration{}, errx.ErrorAgglomerationIsSuspended.Raise(
			fmt.Errorf("agglomeration is suspended"),
		)
	}

	err = s.checkPermissionByCode(
		ctx,
		accountID,
		agglomerationID,
		entity.RolePermissionManageAgglomeration,
	)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	return s.ActivateAgglomeration(ctx, agglomerationID)
}
