package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error) {
	agglo, err := s.repo.UpdateAgglomerationStatus(ctx, ID, entity.AgglomerationStatusInactive)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to deactivate agglomeration: %w", err),
		)
	}

	err = s.messager.WriteAgglomerationDeactivated(ctx, agglo)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish agglomeration deactivate event: %w", err))
	}

	return agglo, nil
}

func (s Service) DeactivateAgglomerationByUser(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
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
			fmt.Errorf("initiator has no access to deactivate agglomeration"),
		)
	}

	return s.DeactivateAgglomeration(ctx, agglomerationID)
}
