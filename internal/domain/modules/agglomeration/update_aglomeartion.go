package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

type UpdateParams struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

func (s Service) UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params UpdateParams) (entity.Agglomeration, error) {
	res, err := s.repo.UpdateAgglomeration(ctx, ID, params)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update agglomeration: %w", err),
		)
	}

	err = s.messager.WriteAgglomerationUpdated(ctx, res)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish agglomeration updated event: %w", err),
		)
	}

	return res, nil
}

func (s Service) UpdateAgglomerationByUser(ctx context.Context, accountID, agglomerationID uuid.UUID, params UpdateParams) (entity.Agglomeration, error) {
	agglo, err := s.GetAgglomeration(ctx, agglomerationID)
	if err != nil {
		return entity.Agglomeration{}, err
	}

	if agglo.Status == entity.AgglomerationStatusSuspended {
		return entity.Agglomeration{}, errx.ErrorAgglomerationIsSuspended.Raise(
			fmt.Errorf("agglomeration is suspended"),
		)
	}

	access, err := s.repo.CheckAccountHavePermissionByCode(
		ctx,
		accountID,
		entity.RolePermissionManageAgglomeration.String(),
	)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return entity.Agglomeration{}, errx.ErrorNotEnoughRightsForAgglomeration.Raise(
			fmt.Errorf("initiator has no access to activate agglomeration"),
		)
	}

	return s.UpdateAgglomeration(ctx, agglomerationID, params)
}
