package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
)

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	agglomeration, err := s.GetAgglomeration(ctx, ID)
	if err != nil {
		return err
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
