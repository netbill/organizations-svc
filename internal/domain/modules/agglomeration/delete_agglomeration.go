package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
)

func (s Service) DeleteAgglomeration(ctx context.Context, accountID, ID uuid.UUID) error {
	agglomeration, err := s.GetAgglomeration(ctx, ID)
	if err != nil {
		return err
	}

	role, err := s.repo.GetAccountMaxRoleInAgglomeration(ctx, accountID, ID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member max role: %w", err),
		)
	}

	if role.Head != true {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only agglomeration head can delete agglomeration"),
		)
	}

	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.DeleteAgglomeration(ctx, ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationDeleted(ctx, agglomeration)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration delete event: %w", err),
			)
		}

		return nil
	})
}
