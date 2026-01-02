package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

type UpdateParams struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

func (s Service) UpdateAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
	params UpdateParams,
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

	initiator, err := s.getInitiator(ctx, accountID, agglo.ID)
	if err != nil {
		return models.Agglomeration{}, err
	}

	if err = s.chekPermissionForManageAgglomeration(
		ctx,
		initiator.ID,
	); err != nil {
		return models.Agglomeration{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomeration(ctx, agglomerationID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationUpdated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration updated event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Agglomeration{}, err
	}

	return agglo, nil
}
