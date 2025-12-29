package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type UpdateParams struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

func (s Service) UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params UpdateParams) (agglo models.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.UpdateAgglomeration(ctx, ID, params)
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

func (s Service) UpdateAgglomerationByUser(
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
	err = s.checkPermissionForManageAgglomeration(
		ctx,
		accountID,
		agglomerationID,
	)
	if err != nil {
		return models.Agglomeration{}, err
	}

	return s.UpdateAgglomeration(ctx, agglomerationID, params)
}
